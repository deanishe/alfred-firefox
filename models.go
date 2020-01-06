// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"net/url"
	"strings"
)

/*
type Window struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Active bool   `json:"active"`
	Tabs   []Tab  `json:"tabs"`
}

func (w Window) String() string {
	return fmt.Sprintf("Window(id=%d, title=%q, active=%v)", w.ID, w.Title, w.Active)
}
*/

// Tab represents a Firefox tab. It contains a subset of the properties
// of the tab.Tab object from Firefox's extensions API.
// https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/tabs/Tab
type Tab struct {
	ID       int    `json:"id"`       // unique ID of tab
	WindowID int    `json:"windowId"` // unique ID of window tab belongs to
	Index    int    `json:"index"`    // position of tab in window
	Title    string `json:"title"`    // tab's title
	URL      string `json:"url"`      // tab's URL
	Active   bool   `json:"active"`   // whether tab is the active tab in its window
}

func (t Tab) String() string {
	return fmt.Sprintf("Tab(id=%d, title=%q, url=%q, active=%v)", t.ID, t.Title, t.URL, t.Active)
}

// Bookmark represents a Firefox bookmark. It contains a subset of the properties
// of the bookmarks.BookmarkTreeNode object from the extensions API.
// https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/bookmarks/BookmarkTreeNode
type Bookmark struct {
	ID       string `json:"id"`       // unique ID
	Title    string `json:"title"`    // bookmark title
	Type     string `json:"type"`     // "bookmark" or "folder"
	URL      string `json:"url"`      // only present for type "bookmark"
	ParentID string `json:"parentId"` // ID of folder bookmark belongs to
	Index    int    `json:"index"`    // position in containing folder
}

func (bm Bookmark) String() string {
	return fmt.Sprintf("Bookmark(id=%q, title=%q, url=%q)", bm.ID, bm.Title, bm.URL)
}

// IsBookmarklet returns true of bookmark URL starts with "javascript:"
func (bm Bookmark) IsBookmarklet() bool {
	return strings.HasPrefix(bm.URL, "javascript:")
}

// JavaScript extracts JS code from a bookmarklet's URL. Returns an empty string
// if Bookmark is not a bookmarklet.
func (bm Bookmark) JavaScript() string {
	if !bm.IsBookmarklet() {
		return ""
	}
	s := strings.TrimPrefix(bm.URL, "javascript:")
	s, _ = url.PathUnescape(s)
	return s
}

// History is an entry from the browser history.
type History struct {
	ID    string `json:"id"`    // unique ID
	Title string `json:"title"` // page title
	URL   string `json:"url"`   // page URL
}

func (h History) String() string {
	return fmt.Sprintf("History(id=%q, title=%q, url=%q)", h.ID, h.Title, h.URL)
}
