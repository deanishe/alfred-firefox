// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"github.com/peterbourgon/ff/ffcli"
)

var (
	bookmarksCmd = &ffcli.Command{
		Name:      "bookmarks",
		Usage:     "firefox -query <query> bookmarks",
		ShortHelp: "search bookmarks",
		LongHelp:  wrap(`Search Firefox bookmarks.`),
		Exec:      runBookmarks,
	}
	bookmarkletsCmd = &ffcli.Command{
		Name:      "bookmarklets",
		Usage:     "firefox -query <query> bookmarklets",
		ShortHelp: "search bookmarklets",
		LongHelp:  wrap(`Search Firefox bookmarklets and execute in frontmost tab.`),
		Exec:      runBookmarklets,
	}
	openURLCmd = &ffcli.Command{
		Name:      "open-url",
		Usage:     "firefox -url <url> open-url",
		ShortHelp: "open URL",
		LongHelp:  wrap(`Open specified URL`),
		Exec:      runOpenURL,
	}

	runBookmarkletCmd = &ffcli.Command{
		Name:      "run-bookmarklet",
		Usage:     "firefox -bookmark <id> run-bookmarklet",
		ShortHelp: "execute bookmarklet in the active tab",
		LongHelp:  wrap(`Executes the specified bookmarklet in a tab.`),
		Exec:      runBookmarklet,
	}

	tabsCmd = &ffcli.Command{
		Name:      "tabs",
		Usage:     "firefox [-query <query>] tabs",
		ShortHelp: "filter Firefox tabs",
		LongHelp:  wrap(`Filter Firefox tabs and perform actions on them.`),
		Exec:      runTabs,
	}

	currentTabCmd = &ffcli.Command{
		Name:      "current-tab",
		Usage:     "firefox [-query <query>] current-tab",
		ShortHelp: "actions for current tab",
		LongHelp:  wrap(`Filter and run actions for current tab`),
		Exec:      runCurrentTab,
	}

	tabCmd = &ffcli.Command{
		Name:      "tab",
		Usage:     "firefox -tab <id> -action <name> tab",
		ShortHelp: "execute tab action",
		LongHelp:  wrap(`Execute specified action on tab`),
		Exec:      runTabAction,
	}

	urlCmd = &ffcli.Command{
		Name:      "url",
		Usage:     "firefox -url <url> -action <name> url",
		ShortHelp: "execute URL action",
		LongHelp:  wrap(`Execute specified action on URL`),
		Exec:      runURLAction,
	}

	actionsCmd = &ffcli.Command{
		Name:      "actions",
		Usage:     "firefox [-tab <id>] [-url <url>] [-query <query>] tab",
		ShortHelp: "filter tab/URL actions",
		LongHelp:  wrap(`View/filter and execute tab/URL actions.`),
		Exec:      runActions,
	}
)

func runOpenURL(_ []string) error {
	log.Printf("opening URL %q ...", URL)
	_, err := util.RunCmd(exec.Command("open", URL))
	return err
}

func searchBookmarks(query string) ([]Bookmark, error) {
	c := mustClient()
	return c.Bookmarks(query)
}

func runBookmarks(_ []string) error {
	if len(query) < 3 {
		wf.Warn("Query Too Short", "Please enter at least 3 characters")
		return nil
	}

	log.Printf("searching bookmarks for %q ...", query)
	bookmarks, err := searchBookmarks(query)
	if err != nil {
		return err
	}

	custom := loadCustomActions()
	for _, bm := range bookmarks {
		if bm.IsBookmarklet() {
			continue
		}
		it := wf.NewItem(bm.Title).
			Subtitle(bm.URL).
			Arg(bm.URL).
			UID(bm.ID).
			Valid(true).
			Icon(iconBookmark).
			Var("CMD", "url").
			Var("ACTION", urlDefault).
			Var("URL", bm.URL).
			Var("TITLE", bm.Title)

		it.NewModifier(aw.ModCmd).
			Subtitle("Other Actions…").
			Arg("").
			Icon(iconMore).
			Var("CMD", "actions")

		custom.Add(it, false)
	}

	wf.WarnEmpty("No Results", "Try a different query?")
	wf.SendFeedback()
	return nil
}

func runBookmarklets(_ []string) error {
	if len(query) < 3 {
		wf.Warn("Query Too Short", "Please enter at least 3 characters")
		return nil
	}

	log.Printf("searching bookmarklets for %q ...", query)
	bookmarks, err := searchBookmarks(query)
	if err != nil {
		return err
	}

	for _, bm := range bookmarks {
		if !bm.IsBookmarklet() {
			continue
		}
		wf.NewItem(bm.Title).
			Subtitle("↩ to execute in current tab").
			UID(bm.ID).
			Copytext("bkm:"+bm.ID+","+bm.Title).
			Arg(bm.URL).
			Icon(iconBookmarklet).
			Valid(true).
			Var("CMD", "run-bookmarklet").
			Var("BOOKMARK", bm.ID)
	}

	wf.WarnEmpty("No Results", "Try a different query?")
	wf.SendFeedback()
	return nil
}

func runBookmarklet(_ []string) error {
	log.Printf("running bookmarklet %q in tab #%d ...", bookmarkID, tabID)
	c := mustClient()
	return c.RunBookmarklet(RunBookmarkletArg{BookmarkID: bookmarkID, TabID: tabID})
}

func runTabs(_ []string) error {
	log.Printf("fetching tabs for query %q ...", query)

	var (
		c    = mustClient()
		tabs []Tab
		err  error
	)
	if tabs, err = c.Tabs(); err != nil {
		return err
	}

	custom := loadCustomActions()
	for _, t := range tabs {
		id := fmt.Sprintf("%d", t.ID)
		it := wf.NewItem(t.Title).
			Subtitle(t.URL).
			Arg(t.URL).
			UID(t.Title).
			Valid(true).
			Var("CMD", "tab").
			Var("ACTION", "Activate Tab").
			Var("TAB", id).
			Var("URL", t.URL).
			Var("TITLE", t.Title)

		it.NewModifier(aw.ModCmd).
			Subtitle("Other Actions").
			Arg("").
			Icon(iconMore).
			Var("CMD", "actions")

		custom.Add(it, true)
	}

	if query != "" {
		_ = wf.Filter(query)
	}

	wf.WarnEmpty("No Matching Tabs", "Try a different query?")
	wf.SendFeedback()
	return nil
}

func runTabAction(_ []string) error {
	log.Printf("running action %q on tab #%d ...", action, tabID)
	a, ok := tabActions[action]
	if !ok {
		return fmt.Errorf("unknown action %q", action)
	}
	return a.Run(tabID)
}

func runURLAction(_ []string) error {
	log.Printf("running action %q on URL %q ...", action, URL)
	a, ok := urlActions[action]
	if !ok {
		return fmt.Errorf("unknown action %q", action)
	}
	return a.Run(URL)
}

func runCurrentTab(_ []string) error {
	tab, err := mustClient().CurrentTab()
	if err != nil {
		return err
	}
	tabID = tab.ID
	URL = tab.URL
	return runActions([]string{})
}

func runActions(_ []string) error {
	if tabID != 0 {
		for _, a := range tabActions {
			wf.NewItem(a.Name()).
				UID(a.Name()).
				Copytext(a.Name()).
				Icon(a.Icon()).
				Valid(true).
				Var("CMD", "tab").
				Var("ACTION", a.Name()).
				Var("TAB", fmt.Sprintf("%d", tabID))
		}

		// add custom bookmarklet commands
		for _, a := range loadCustomActions() {
			if a.kind != "bookmarklet" {
				continue
			}
			wf.NewItem(a.name).
				UID(a.id).
				Copytext("bkm:"+a.id+","+a.name).
				Icon(iconBookmarklet).
				Valid(true).
				Var("CMD", "run-bookmarklet").
				Var("BOOKMARK", a.id).
				Var("TAB", fmt.Sprintf("%d", tabID))
		}
	}

	if URL != "" {
		for _, a := range urlActions {
			wf.NewItem(a.Name()).
				UID(a.Name()).
				Copytext(a.Name()).
				Icon(a.Icon()).
				Valid(true).
				Var("CMD", "url").
				Var("ACTION", a.Name()).
				Var("URL", URL)
		}
	}

	if query != "" {
		_ = wf.Filter(query)
	}

	wf.WarnEmpty("No Matching Actions", "Try a different query?")
	wf.SendFeedback()
	return nil
}

func mustClient() *rpcClient {
	c, err := newClient()
	if err != nil {
		panic(err)
	}
	return c
}

// actions mapped to hotkeys by user via environment variables.
type customActions []customAction

// add custom actions to a Bookmark/Tab item.
func (ca customActions) Add(it *aw.Item, tab bool) {
	for _, a := range ca {
		if a.kind == "tab" && !tab {
			continue
		}
		a.Add(it)
	}
}

// action defined via environment variable.
type customAction struct {
	kind string // "tab", "url" or "bookmarklet"
	id   string // only set on "bookmarklet" actions
	name string // human-readable name of action
	mods []aw.ModKey
}

func (ca customAction) Add(it *aw.Item) {
	m := it.NewModifier(ca.mods...).
		Subtitle(ca.name).
		Var("ACTION", ca.name)

	switch ca.kind {
	case "tab":
		m.Icon(iconTab).Var("CMD", "tab")
	case "url":
		m.Var("CMD", "url").Icon(iconURL)
		if icon, ok := scriptIcons[ca.name]; ok {
			m.Icon(icon)
		}
	case "bookmarklet":
		m.Var("CMD", "run-bookmarklet").Var("BOOKMARK", ca.id).
			Icon(iconBookmarklet)
	}

}

func loadCustomActions() customActions {
	var (
		actions customActions
		key     string
		id      string
		name    string
		kind    string
		mods    []aw.ModKey
	)
	for _, s := range os.Environ() {
		parts := strings.SplitN(s, "=", 2)
		key, name = strings.ToLower(parts[0]), parts[1]
		if !strings.HasPrefix(key, "url_") && !strings.HasPrefix(key, "tab_") {
			continue
		}
		if key == "url_default" {
			continue
		}

		kind = key[0:3]
		if mods = parseMods(key[4:]); len(mods) == 0 {
			log.Printf("[warning] invalid modifier for %q", key)
			// Add empty modifiers. The action won't be attached to any
			// modifier key, but will still show up in action menus
			// (if it's a bookmarklet). This allows adding bookmarklets to
			// action menus without assigning them a hotkey.
			mods = []aw.ModKey{""}
		}

		if strings.HasPrefix(name, "bkm:") {
			kind = "bookmarklet"
			parts := strings.SplitN(name[4:], ",", 2)
			id, name = parts[0], parts[1]
		}

		actions = append(actions, customAction{kind, id, name, mods})
	}
	return actions
}

func parseMods(s string) []aw.ModKey {
	var keys []aw.ModKey
	for _, v := range strings.Split(s, "_") {
		switch v {
		case "cmd":
			keys = append(keys, aw.ModCmd)
		case "opt", "alt":
			keys = append(keys, aw.ModOpt)
		case "ctrl":
			keys = append(keys, aw.ModCtrl)
		case "shift":
			keys = append(keys, aw.ModShift)
		default:
			log.Printf("[warning] unknown modifier: %q", v)
		}
	}
	return keys
}
