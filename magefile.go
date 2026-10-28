// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/deanishe/awgo/util/build"
)

const (
	buildDir   = "./build"
	distDir    = "./dist"
	binaryName = "alfred-firefox"
)

var (
	assetPatterns = []string{
		"info.plist",
		"*.png",
		"scripts/*",
		"icons/*.png",
		"LICENCE.txt",
		"README.md",
	}
	info *build.Info
)

func init() {
	var err error
	if info, err = build.NewInfo(); err != nil {
		// Alfred isn't installed/configured (e.g. on CI), which is fine for
		// everything except Link. Point at a throwaway preferences dir.
		os.Setenv("alfred_preferences", filepath.Join(os.TempDir(), "alfred-firefox", "Alfred.alfredpreferences"))
		os.Setenv("alfred_version", "5")
		if info, err = build.NewInfo(); err != nil {
			panic(err)
		}
	}
}

// Aliases are aliases for Mage commands.
var Aliases = map[string]interface{}{
	"b":   Build,
	"ext": Extension,
}

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Dist build & export a universal workflow in dist directory
func Dist() {
	mg.Deps(cleanBuild, Deps)
	mg.SerialDeps(Universal, Export)
}

// Test run go vet and go test
func Test() error {
	if err := sh.RunV("go", "vet", "./..."); err != nil {
		return err
	}
	return sh.RunV("go", "test", "./...")
}

// Extension build browser extension in dist directory
func Extension() error {
	if err := sh.Run("mkdir", "-p", distDir); err != nil {
		return err
	}
	if err := sh.Rm("dist/alfred-firefox-assistant.xpi"); err != nil {
		return err
	}
	return sh.Run("zsh", "-c",
		`cd extension && zip -r ../dist/alfred-firefox-assistant.xpi * -x '*.DS_Store'`)
}

// Export compile build directory into workflow in dist directory
func Export() error {
	fmt.Println("exporting workflow ...")
	path, err := build.Export(buildDir, distDir)
	if err != nil {
		return err
	}
	fmt.Println("exported workflow:", path)
	return nil
}

// Assets copy assets to build directory
func Assets() error {
	mg.Deps(cleanAssets)
	fmt.Println("linking assets to build directory ...")
	globs := build.Globs(assetPatterns...)

	if err := build.SymlinkGlobs(buildDir, globs...); err != nil {
		return fmt.Errorf("symlink %v to %q: %v", globs, buildDir, err)
	}
	return sh.Run("cp", "-vf", "server.sh", buildDir)
}

// Build build executable in build directory
func Build() error {
	mg.Deps(Assets)
	fmt.Println("building executable ...")
	if err := os.MkdirAll(buildDir, 0750); err != nil {
		return err
	}
	return sh.RunWith(info.Env(), "go", "build", "-o", buildDir+"/"+binaryName, ".")
}

// Universal build a universal (arm64 + amd64) executable in build directory
func Universal() error {
	mg.Deps(Assets)
	fmt.Println("building universal executable ...")
	if err := os.MkdirAll(buildDir, 0750); err != nil {
		return err
	}
	archs := []string{"arm64", "amd64"}
	parts := make([]string, 0, len(archs))
	for _, arch := range archs {
		env := map[string]string{}
		for k, v := range info.Env() {
			env[k] = v
		}
		env["GOOS"] = "darwin"
		env["GOARCH"] = arch
		env["CGO_ENABLED"] = "0"
		out := fmt.Sprintf("%s/%s-%s", buildDir, binaryName, arch)
		if err := sh.RunWith(env, "go", "build", "-o", out, "."); err != nil {
			return err
		}
		parts = append(parts, out)
	}
	args := append([]string{"-create", "-output", buildDir + "/" + binaryName}, parts...)
	if err := sh.Run("lipo", args...); err != nil {
		return err
	}
	for _, p := range parts {
		if err := os.Remove(p); err != nil {
			return err
		}
	}
	return nil
}

// Link symlink build directory to Alfred's workflow directory
func Link() error {
	mg.Deps(Build)
	fmt.Printf("linking workflow to %q ...\n", info.InstallDir)
	if err := sh.Rm(info.InstallDir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(info.InstallDir), 0750); err != nil {
		return err
	}
	if err := build.Symlink(info.InstallDir, buildDir, true); err != nil {
		return fmt.Errorf("symlink %q to %q: %v", buildDir, info.InstallDir, err)
	}
	return nil
}

// Deps tidy & download deps
func Deps() error {
	fmt.Println("verifying dependencies ...")
	if err := sh.RunWith(info.Env(), "go", "mod", "tidy", "-v"); err != nil {
		return err
	}
	return sh.RunWith(info.Env(), "go", "mod", "download")
}

// Clean remove built files
func Clean() {
	fmt.Println("cleaning ...")
	mg.Deps(cleanBuild, cleanMage)
}

// CleanDist delete exported workflows in dist directory
func CleanDist() error {
	fmt.Println("cleaning dist directory ...")
	return cleanDir(distDir, func(fi os.FileInfo) bool {
		return true
	})
}

func cleanBinary() error {
	return cleanDir(buildDir, func(fi os.FileInfo) bool {
		if fi.Name() == binaryName {
			return true
		}
		return false
	})
}

func cleanAssets() error {
	return cleanDir(buildDir, func(fi os.FileInfo) bool {
		if fi.IsDir() || filepath.Ext(fi.Name()) != "" {
			return true
		}
		return false
	})
}

func cleanBuild() {
	mg.Deps(cleanBinary, cleanAssets)
}

func cleanMage() error {
	fmt.Println("cleaning mage cache ...")
	return sh.Run("mage", "-clean")
}

func cleanDir(dirname string, match func(os.FileInfo) bool) error {
	var paths []string

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		return nil
	}

	err := filepath.Walk(dirname, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == dirname { // ignore root
			return nil
		}
		if match(fi) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, p := range paths {
		if err := os.RemoveAll(p); err != nil {
			return err
		}
		fmt.Println("deleted:", p)
	}

	return nil
}
