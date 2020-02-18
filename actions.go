// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
)

var (
	tabActions = map[string]tabAction{}
	urlActions = map[string]urlAction{}
)

type tabAction interface {
	Name() string
	Icon() *aw.Icon
	Run(tabID int) error
}

type urlAction interface {
	Name() string
	Icon() *aw.Icon
	Run(URL string) error
}

func init() {
	for _, a := range []tabAction{
		tAction{name: "Activate Tab", action: "activate", icon: iconTab},
		tAction{name: "Close Tabs to Left", action: "close-left", icon: iconTab},
		tAction{name: "Close Tabs to Right", action: "close-right", icon: iconTab},
		tAction{name: "Close Other Tabs", action: "close-other", icon: iconTab},
	} {
		tabActions[a.Name()] = a
	}

	a := openIncognito{}
	urlActions[a.Name()] = a
}

func loadURLActions() error {
	var (
		scripts = map[string]string{}
		infos   []os.FileInfo
		err     error
	)
	for _, dir := range scriptDirs {
		if infos, err = ioutil.ReadDir(dir); err != nil {
			return err
		}

		for _, fi := range infos {
			if fi.IsDir() {
				continue
			}

			var (
				path     = filepath.Join(dir, fi.Name())
				ext      = strings.ToLower(filepath.Ext(fi.Name()))
				name     = fi.Name()[0 : len(fi.Name())-len(ext)]
				_, known = util.DefaultInterpreters[ext]
				exe      = fi.Mode()&0111 != 0
			)
			if exe || known {
				scripts[name] = path
			}

			if imageExts[ext] {
				scriptIcons[name] = &aw.Icon{Value: path}
			}
		}
	}

	for name, path := range scripts {
		log.Printf("loaded URL action %q from %q", name, util.PrettyPath(path))
		a := uAction{
			name:   name,
			icon:   actionIcon(name, iconURL),
			script: path,
		}
		urlActions[name] = a
	}

	return nil
}

type tAction struct {
	name   string
	icon   *aw.Icon
	action string
}

func (a tAction) Name() string   { return a.name }
func (a tAction) Icon() *aw.Icon { return a.icon }
func (a tAction) Run(tabID int) error {
	c := mustClient()
	switch a.action {
	case "activate":
		_, err := util.RunAS(`tell application "Firefox" to activate`)
		if err != nil {
			return err
		}
		return c.ActivateTab(tabID)
	case "close-left":
		return c.CloseTabsLeft(tabID)
	case "close-right":
		return c.CloseTabsRight(tabID)
	case "close-other":
		return c.CloseTabsOther(tabID)
	default:
		return fmt.Errorf("unknown action %q", action)
	}
}

type uAction struct {
	name   string
	icon   *aw.Icon
	script string
}

func (a uAction) Name() string   { return a.name }
func (a uAction) Icon() *aw.Icon { return a.icon }
func (a uAction) Run(URL string) error {
	data, err := util.Run(a.script, URL)
	if err != nil {
		return err
	}
	s := string(data)
	if s != "" {
		log.Print(util.Pad(fmt.Sprintf(" output: %q ", a.name), "-", 50))
		log.Print(s)
	}
	return nil
}

// URL action to open a URL in a new incognito window
type openIncognito struct{}

func (a openIncognito) Name() string   { return "Open in Incognito Window" }
func (a openIncognito) Icon() *aw.Icon { return iconIncognito }
func (a openIncognito) Run(URL string) error {
	mustClient().OpenIncognito(URL)
	return nil
}

var (
	_ tabAction = (*tAction)(nil)
	_ urlAction = (*uAction)(nil)
	_ urlAction = openIncognito{}
)
