// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import aw "github.com/deanishe/awgo"

var (
	iconBookmark        = &aw.Icon{Value: "icons/bookmark.png"}
	iconBookmarklet     = &aw.Icon{Value: "icons/bookmarklet.png"}
	iconDocs            = &aw.Icon{Value: "icons/docs.png"}
	iconError           = &aw.Icon{Value: "icons/error.png"}
	iconHistory         = &aw.Icon{Value: "icons/history.png"}
	iconIssue           = &aw.Icon{Value: "icons/issue.png"}
	iconInstall         = &aw.Icon{Value: "icons/install.png"}
	iconMore            = &aw.Icon{Value: "icons/more.png"}
	iconScript          = &aw.Icon{Value: "icons/script.png"}
	iconTab             = &aw.Icon{Value: "icons/tab.png"}
	iconUpdateAvailable = &aw.Icon{Value: "icons/update-available.png"}
	iconUpdateOK        = &aw.Icon{Value: "icons/update-ok.png"}
	iconURL             = &aw.Icon{Value: "icons/url.png"}
	iconWarning         = &aw.Icon{Value: "icons/warning.png"}

	// populated by loadURLActions
	scriptIcons = map[string]*aw.Icon{}

	imageExts = map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
		".icns": true,
	}
)

func init() {
	aw.IconError = iconError
	aw.IconWarning = iconWarning
}

// return custom icon or fallback
func actionIcon(name string, fallback *aw.Icon) *aw.Icon {
	if icon, ok := scriptIcons[name]; ok {
		return icon
	}
	return fallback
}
