// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import "net/rpc"

type rpcClient struct {
	client *rpc.Client
}

func newClient() (*rpcClient, error) {
	c, err := rpc.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}
	return &rpcClient{c}, nil
}

func (c *rpcClient) Bookmarks(query string) ([]Bookmark, error) {
	var bookmarks []Bookmark
	err := c.client.Call("Firefox.Bookmarks", query, &bookmarks)
	return bookmarks, err
}

func (c *rpcClient) Tabs() ([]Tab, error) {
	var tabs []Tab
	err := c.client.Call("Firefox.Tabs", "", &tabs)
	return tabs, err
}

func (c *rpcClient) CurrentTab() (Tab, error) {
	var tab Tab
	err := c.client.Call("Firefox.CurrentTab", "", &tab)
	return tab, err
}

func (c *rpcClient) ActivateTab(tabID int) error {
	return c.client.Call("Firefox.ActivateTab", tabID, nil)
}

func (c *rpcClient) CloseTabsLeft(tabID int) error {
	return c.client.Call("Firefox.CloseTabsLeft", tabID, nil)
}

func (c *rpcClient) CloseTabsRight(tabID int) error {
	return c.client.Call("Firefox.CloseTabsRight", tabID, nil)
}

func (c *rpcClient) CloseTabsOther(tabID int) error {
	return c.client.Call("Firefox.CloseTabsOther", tabID, nil)
}

func (c *rpcClient) RunJS(script string) error {
	return c.client.Call("Firefox.RunJS", script, nil)
}

func (c *rpcClient) RunBookmarklet(arg RunBookmarkletArg) error {
	return c.client.Call("Firefox.RunBookmarklet", arg, nil)
}
