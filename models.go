// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"net/url"
	"strings"
)

type Window struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Active bool   `json:"active"`
	Tabs   []Tab  `json:"tabs"`
}

func (w Window) String() string {
	return fmt.Sprintf("Window(id=%d, title=%q, active=%v)", w.ID, w.Title, w.Active)
}

type Tab struct {
	ID       int    `json:"id"`
	WindowID int    `json:"windowId"`
	Index    int    `json:"index"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Type     string `json:"type"`
	Active   bool   `json:"active"`
}

func (t Tab) String() string {
	return fmt.Sprintf("Tab(id=%d, title=%q, url=%q, active=%v)", t.ID, t.Title, t.URL, t.Active)
}

type Bookmark struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	ParentID string `json:"parentId"`
	Index    int    `json:"index"`
}

func (bm Bookmark) String() string {
	return fmt.Sprintf("Bookmark(id=%q, title=%q, url=%q)", bm.ID, bm.Title, bm.URL)
}

func (bm Bookmark) IsBookmarklet() bool {
	return strings.HasPrefix(bm.URL, "javascript:")
}

func (bm Bookmark) JavaScript() string {
	if !bm.IsBookmarklet() {
		return ""
	}
	s := strings.TrimPrefix(bm.URL, "javascript:")
	s, _ = url.PathUnescape(s)
	return s
}
