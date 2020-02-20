// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"

	"github.com/deanishe/awgo/util"
)

// rpcServer provides the RPC API. It passes commands and responses between
// RPC clients and the Firefox extension.
type rpcServer struct {
	ff       *firefox // native application client run by FF extension
	sock     string   // path to UNIX socket for RPC
	listener net.Listener
	server   *rpc.Server
}

// create new RPC server on socket specified by filepath addr
func newRPCService(addr string, client *firefox) (*rpcServer, error) {
	var err error
	s := &rpcServer{
		ff:     client,
		sock:   addr,
		server: rpc.NewServer(),
	}

	if err = s.server.RegisterName("Firefox", s); err != nil {
		return nil, err
	}

	if s.listener, err = net.Listen("unix", s.sock); err != nil {
		return nil, err
	}

	if err = os.Chmod(addr, 0600); err != nil {
		return nil, err
	}

	return s, nil
}

// AppName returns name of browser running the server.
func (s *rpcServer) AppName(_ string, name *string) error {
	*name = browserName
	return nil
}

// Ping checks connection to Firefox extension. Extension responds with "pong".
func (s *rpcServer) Ping(_ string, result *string) error {
	defer util.Timed(time.Now(), "ping")
	var r responseString
	if err := s.ff.call("ping", nil, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	*result = r.String
	return nil
}

// func (s *rpcServer) Windows(_ string, windows *[]Window) error {
// 	defer util.Timed(time.Now(), "get windows")
// 	var r responseWindows
// 	if err := s.ff.call("all-windows", nil, &r); err != nil {
// 		return err
// 	}
// 	*windows = r.Windows
// 	return nil
// }

// Tabs returns all Firefox tabs.
func (s *rpcServer) Tabs(_ string, tabs *[]Tab) error {
	defer util.Timed(time.Now(), "get tabs")
	var r responseTabs
	if err := s.ff.call("all-tabs", nil, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	*tabs = r.Tabs
	return nil
}

// ActivateTab brings the specified tab to the front.
func (s *rpcServer) ActivateTab(tabID int, _ *struct{}) error {
	defer util.Timed(time.Now(), "activate tab")
	var r responseNone
	if err := s.ff.call("activate-tab", tabID, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}

// Tab returns the specified tab. If tabID is 0, returns the active tab.
func (s *rpcServer) Tab(tabID int, tab *Tab) error {
	defer util.Timed(time.Now(), "get tab")
	var r responseTab
	if err := s.ff.call("tab", tabID, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	*tab = r.Tab
	return nil
}

/*
// CurrentTab returns the currently-active tab.
func (s *rpcServer) CurrentTab(_ string, tab *Tab) error {
	defer util.Timed(time.Now(), "get current tab")
	var r responseTab
	if err := s.ff.call("tab", 0, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	*tab = r.Tab
	return nil
}
*/

// CloseTabsLeft closes tabs to the left of specified tab.
func (s *rpcServer) CloseTabsLeft(tabID int, _ *struct{}) error {
	defer util.Timed(time.Now(), "close tabs to left")
	var r responseNone
	if err := s.ff.call("close-tabs-left", tabID, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}

// CloseTabsRight closes tabs to the right of specified tab.
func (s *rpcServer) CloseTabsRight(tabID int, _ *struct{}) error {
	defer util.Timed(time.Now(), "close tabs to right")
	var r responseNone
	if err := s.ff.call("close-tabs-right", tabID, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}

// CloseTabsOther closes other tabs in same window as the specified one.
func (s *rpcServer) CloseTabsOther(tabID int, _ *struct{}) error {
	defer util.Timed(time.Now(), "close other tabs")
	var r responseNone
	if err := s.ff.call("close-tabs-other", tabID, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}

// Bookmarks returns all Firefox bookmarks matching query.
func (s *rpcServer) Bookmarks(query string, bookmarks *[]Bookmark) error {
	defer util.Timed(time.Now(), fmt.Sprintf("search bookmarks for %q", query))
	var (
		r   responseBookmarks
		err error
	)
	if query == "" {
		err = s.ff.call("all-bookmarks", nil, &r)
	} else {
		err = s.ff.call("search-bookmarks", query, &r)
	}
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	*bookmarks = r.Bookmarks
	return nil
}

// History searches Firefox browsing history.
func (s *rpcServer) History(query string, history *[]History) error {
	defer util.Timed(time.Now(), fmt.Sprintf("search history for %q", query))
	var (
		r   responseHistory
		err error
	)
	err = s.ff.call("search-history", query, &r)
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	*history = r.Entries
	return nil
}

// Downloads searches Firefox downloads.
func (s *rpcServer) Downloads(query string, downloads *[]Download) error {
	defer util.Timed(time.Now(), fmt.Sprintf("search download for %q", query))
	var (
		r   responseDownload
		err error
	)
	err = s.ff.call("search-downloads", query, &r)
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	*downloads = r.Downloads
	return nil
}

// OpenIncognito opens URL in a new incognito window.
func (s *rpcServer) OpenIncognito(URL string, _ *struct{}) error {
	defer util.Timed(time.Now(), "open incognito")
	var r responseNone
	if err := s.ff.call("open-incognito", URL, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}

// RunJSArg is the arguments required for RunJS call. TabID may be 0, in which
// case the JavaScript is executed in the active tab.
type RunJSArg struct {
	TabID int    `json:"tabId"`
	JS    string `json:"js"`
}

// RunJS executes JavaScript in the specified tab. If tabID is 0, the
// script is executed in the current tab.
func (s *rpcServer) RunJS(arg RunJSArg, JSON *string) error {
	defer util.Timed(time.Now(), "execute JS")
	var r responseString
	if err := s.ff.call("execute-js", arg, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	*JSON = r.String
	return nil
}

// RunBookmarkletArg is the arguments required for RunBookmarklet call.
// TabID may be 0, in which case the bookmarklet is executed in the active tab.
type RunBookmarkletArg struct {
	TabID      int    `json:"tabId"`
	BookmarkID string `json:"bookmarkId"`
}

// RunBookmarklet executes a given bookmarklet in a given tab.
func (s *rpcServer) RunBookmarklet(arg RunBookmarkletArg, _ *struct{}) error {
	defer util.Timed(time.Now(), "run bookmarklet")
	var r responseNone
	if err := s.ff.call("run-bookmarklet", arg, &r); err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return nil
}

func (s *rpcServer) run() {
	log.Printf("serving RPC on %q ...", s.sock)
	s.server.Accept(s.listener)
}

func (s *rpcServer) stop() error {
	return s.listener.Close()
}

type responseString struct {
	String string `json:"payload"`
	Error  string `json:"error"`
}

// type responseWindows struct {
// 	Windows []Window `json:"payload"`
// }

type responseTabs struct {
	Tabs  []Tab  `json:"payload"`
	Error string `json:"error"`
}

type responseTab struct {
	Tab   Tab    `json:"payload"`
	Error string `json:"error"`
}

type responseHistory struct {
	Entries []History `json:"payload"`
	Error   string    `json:"error"`
}

type responseTabCurrent struct {
	Tab   Tab    `json:"payload"`
	Error string `json:"error"`
}

type responseBookmarks struct {
	Bookmarks []Bookmark `json:"payload"`
	Error     string     `json:"error"`
}

type responseDownload struct {
	Downloads []Download `json:"payload"`
	Error     string     `json:"error"`
}

type responseBool struct {
	OK    bool   `json:"payload"`
	Error string `json:"error"`
}

type responseNone struct {
	Error string `json:"error"`
}
