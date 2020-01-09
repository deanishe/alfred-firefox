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
		"icons/*.png",
		"LICENCE.txt",
	}
	info *build.Info
)

func init() {
	var err error
	if info, err = build.NewInfo(); err != nil {
		panic(err)
	}
}

// Aliases are aliases for Mage commands.
var Aliases = map[string]interface{}{
	"b": Build,
}

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Dist build & export workflow
func Dist() {
	mg.Deps(cleanBuild, Deps)
	mg.SerialDeps(Build, Export)
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
