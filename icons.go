// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import aw "github.com/deanishe/awgo"

var (
	iconBookmark        = &aw.Icon{Value: "icons/bookmark.png"}
	iconBookmarklet     = &aw.Icon{Value: "icons/bookmarklet.png"}
	iconError           = &aw.Icon{Value: "icons/error.png"}
	iconHistory         = &aw.Icon{Value: "icons/history.png"}
	iconMore            = &aw.Icon{Value: "icons/more.png"}
	iconTab             = &aw.Icon{Value: "icons/tab.png"}
	iconUpdateAvailable = &aw.Icon{Value: "icons/update-available.png"}
	iconUpdateOK        = &aw.Icon{Value: "icons/update-ok.png"}
	iconURL             = &aw.Icon{Value: "icons/url.png"}
	iconWarning         = &aw.Icon{Value: "icons/warning.png"}
)

func init() {
	aw.IconError = iconError
	aw.IconWarning = iconWarning
}
