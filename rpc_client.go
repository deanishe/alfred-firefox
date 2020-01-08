// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"log"
	"net/rpc"
)

// RPC client used by workflow to execute extension actions.
type rpcClient struct {
	client *rpc.Client
}

// Create new RPC client. Returns an error if connection to server fails.
func newClient() (*rpcClient, error) {
	c, err := rpc.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}
	return &rpcClient{c}, nil
}

// return new RPC client, panicking if it can't connect to server
func mustClient() *rpcClient {
	c, err := newClient()
	if err != nil {
		log.Printf("[ERROR] %v", err)
		panic("Cannot Connect to Extension")
	}
	return c
}

// Ping checks connection to Firefox extension.
func (c *rpcClient) Ping() error {
	var s string
	return c.client.Call("Firefox.Ping", "", &s)
}

// Bookmarks returns all Firefox bookmarks matching query.
func (c *rpcClient) Bookmarks(query string) ([]Bookmark, error) {
	var bookmarks []Bookmark
	err := c.client.Call("Firefox.Bookmarks", query, &bookmarks)
	return bookmarks, err
}

// History searches Firefox browsing history.
func (c *rpcClient) History(query string) ([]History, error) {
	var history []History
	err := c.client.Call("Firefox.History", query, &history)
	return history, err
}

// Downloads searches Firefox downloads.
func (c *rpcClient) Downloads(query string) ([]Download, error) {
	var downloads []Download
	err := c.client.Call("Firefox.Downloads", query, &downloads)
	return downloads, err
}

// Tabs returns all Firefox tabs.
func (c *rpcClient) Tabs() ([]Tab, error) {
	var tabs []Tab
	err := c.client.Call("Firefox.Tabs", "", &tabs)
	return tabs, err
}

// CurrentTab returns the currently-active tab.
func (c *rpcClient) CurrentTab() (Tab, error) {
	var tab Tab
	err := c.client.Call("Firefox.CurrentTab", "", &tab)
	return tab, err
}

// ActivateTab brings the specified tab to the front.
func (c *rpcClient) ActivateTab(tabID int) error {
	return c.client.Call("Firefox.ActivateTab", tabID, nil)
}

// CloseTabsLeft closes tabs to the left of specified tab.
func (c *rpcClient) CloseTabsLeft(tabID int) error {
	return c.client.Call("Firefox.CloseTabsLeft", tabID, nil)
}

// CloseTabsRight closes tabs to the right of specified tab.
func (c *rpcClient) CloseTabsRight(tabID int) error {
	return c.client.Call("Firefox.CloseTabsRight", tabID, nil)
}

// CloseTabsOther closes other tabs in same window as the specified one.
func (c *rpcClient) CloseTabsOther(tabID int) error {
	return c.client.Call("Firefox.CloseTabsOther", tabID, nil)
}

// func (c *rpcClient) RunJS(script string) error {
// 	return c.client.Call("Firefox.RunJS", script, nil)
// }

// RunBookmarklet executes a given bookmarklet in a given tab.
func (c *rpcClient) RunBookmarklet(arg RunBookmarkletArg) error {
	return c.client.Call("Firefox.RunBookmarklet", arg, nil)
}
