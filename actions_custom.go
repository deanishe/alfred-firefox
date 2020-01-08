// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"log"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
)

// actions defined by user via environment variables.
type customActions []customAction

// Add custom actions to a Bookmark/Tab item. If "tab" is false,
// only URL actions are added.
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

	// Modifier keys. If empty, the action isn't bound to
	// any keyboard shortcut, but is still shown in action lists.
	mods []aw.ModKey
}

// add action to Bookmark/Tab item
func (ca customAction) Add(it *aw.Item) {
	if len(ca.mods) == 0 { // list-only action
		return
	}
	m := it.NewModifier(ca.mods...).
		Subtitle(ca.name).
		Var("ACTION", ca.name)

	switch ca.kind {
	case "tab":
		m.Icon(iconTab).Var("CMD", "tab")
	case "url":
		m.Var("CMD", "url").Icon(actionIcon(ca.name, iconURL))
	case "bookmarklet":
		m.Var("CMD", "run-bookmarklet").Var("BOOKMARK", ca.id).
			Icon(iconBookmarklet)
	}

}

// return custom actions set by user
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
			// Warn about, but allow, empty modifiers.
			// The action won't be bound to any modifier key, but
			// it will still show up in tab action lists (if it's
			// a bookmarklet).
			log.Printf("invalid modifier for %q", key)
		}

		if strings.HasPrefix(name, "bml:") {
			kind = "bookmarklet"
			parts := strings.SplitN(name[4:], ",", 2)
			id, name = parts[0], parts[1]
		}

		actions = append(actions, customAction{kind, id, name, mods})
	}
	return actions
}

// parse string of form "cmd_opt_shift" into slice of ModKeys
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
			// default:
			// 	log.Printf("[warning] unknown modifier: %q", v)
		}
	}
	return keys
}
