// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/util"
	"github.com/peterbourgon/ff/ffcli"
)

var (
	// search history
	historyCmd = &ffcli.Command{
		Name:      "history",
		Usage:     "alfred-firefox -query <query> history",
		ShortHelp: "search browsing history",
		LongHelp:  wrap(`Search browser history.`),
		Exec:      runHistory,
	}

	// search downloads
	downloadsCmd = &ffcli.Command{
		Name:      "downloads",
		Usage:     "alfred-firefox -query <query> downloads",
		ShortHelp: "search downloads",
		LongHelp:  wrap(`Search browser downloads.`),
		Exec:      runDownloads,
	}

	// search bookmarks
	bookmarksCmd = &ffcli.Command{
		Name:      "bookmarks",
		Usage:     "alfred-firefox -query <query> bookmarks",
		ShortHelp: "search bookmarks",
		LongHelp:  wrap(`Search browser bookmarks.`),
		Exec:      runBookmarks,
	}

	// search bookmarklets
	bookmarkletsCmd = &ffcli.Command{
		Name:      "bookmarklets",
		Usage:     "alfred-firefox -query <query> bookmarklets",
		ShortHelp: "search bookmarklets",
		LongHelp:  wrap(`Search bookmarklets and execute in frontmost tab.`),
		Exec:      runBookmarklets,
	}

	/*
		// open URL
		// TODO: is this used? can it be removed?
		openURLCmd = &ffcli.Command{
			Name:      "open-url",
			Usage:     "alfred-firefox -url <url> open-url",
			ShortHelp: "open URL",
			LongHelp:  wrap(`Open specified URL`),
			Exec:      runOpenURL,
		}
	*/

	// execute a bookmarklet in the specified tab
	runBookmarkletCmd = &ffcli.Command{
		Name:      "run-bookmarklet",
		Usage:     "alfred-firefox [-tab <id>] -bookmark <id> run-bookmarklet",
		ShortHelp: "execute bookmarklet in the specified tab",
		LongHelp: wrap(`
			Execute a bookmarklet in a tab. Bookmark ID is required.
			If no tab ID is specified, bookmarklet is run in the active tab.
		`),
		Exec: runBookmarklet,
	}

	// filter open tabs
	tabsCmd = &ffcli.Command{
		Name:      "tabs",
		Usage:     "alfred-firefox [-query <query>] tabs",
		ShortHelp: "filter browser tabs",
		LongHelp:  wrap(`Filter browser tabs and perform actions on them.`),
		Exec:      runTabs,
	}

	// filter tab & URL actions for current tab
	currentTabCmd = &ffcli.Command{
		Name:      "current-tab",
		Usage:     "alfred-firefox [-query <query>] current-tab",
		ShortHelp: "actions for current tab",
		LongHelp:  wrap(`Filter and run actions for current tab`),
		Exec:      runCurrentTab,
	}

	infoFlags = flag.NewFlagSet("tab-info", flag.ExitOnError)
	shellVars bool // export tab info as shell variables
	// export info for current tab
	currentTabInfoCmd = &ffcli.Command{
		Name:      "tab-info",
		Usage:     "alfred-firefox tab-info [-shell]",
		ShortHelp: "export current tab info",
		LongHelp:  wrap(`Export current tab info as variables`),
		FlagSet:   infoFlags,
		Exec:      runCurrentTabInfo,
	}

	// run a tab/URL action for the specified tab
	tabCmd = &ffcli.Command{
		Name:      "tab",
		Usage:     "alfred-firefox [-tab <id>] -action <name> tab",
		ShortHelp: "execute tab action",
		LongHelp: wrap(`
			Execute specified action on tab. Both URL and tab actions
			are available on tabs.
			`),
		Exec: runTabAction,
	}

	// inject JS into the specified tab
	injectCmd = &ffcli.Command{
		Name:      "inject",
		Usage:     "alfred-firefox [-tab <id>] inject <script>",
		ShortHelp: "inject JavaScript into tab",
		LongHelp: wrap(`
			Execute JavaScript in specifed tab and return result as JSON.
			`),
		Exec: runInject,
	}

	// run action for URL
	urlCmd = &ffcli.Command{
		Name:      "url",
		Usage:     "alfred-firefox [-url <url>] -action <name> url",
		ShortHelp: "execute URL action",
		LongHelp:  wrap(`Execute specified action on URL`),
		Exec:      runURLAction,
	}

	// filter URL (and tab) actions
	actionsCmd = &ffcli.Command{
		Name:      "actions",
		Usage:     "alfred-firefox [-tab <id>] [-url <url>] [-query <query>] actions",
		ShortHelp: "filter tab/URL actions",
		LongHelp:  wrap(`View/filter and execute tab/URL actions.`),
		Exec:      runActions,
	}

	// check for update
	updateCmd = &ffcli.Command{
		Name:      "update",
		Usage:     "alfred-firefox update",
		ShortHelp: "check for workflow update",
		LongHelp:  wrap(`Check if newer version of workflow is available.`),
		Exec:      runUpdate,
	}

	// show workflow status
	statusCmd = &ffcli.Command{
		Name:      "options",
		Usage:     "alfred-firefox [-query <query>] options",
		ShortHelp: "show workflow status & options",
		LongHelp:  wrap(`Show workflow status, info and options.`),
		Exec:      runStatus,
	}

	// open file in default application
	openCmd = &ffcli.Command{
		Name:      "open",
		Usage:     "alfred-firefox open <path>",
		ShortHelp: "open file in default application",
		LongHelp:  wrap(`Open file in default application.`),
		Exec:      runOpen,
	}

	// reveal file in Finder
	revealCmd = &ffcli.Command{
		Name:      "reveal",
		Usage:     "alfred-firefox reveal <path>",
		ShortHelp: "reveal file in Finder",
		LongHelp:  wrap(`Reveal file in Finder.`),
		Exec:      runReveal,
	}
)

func init() {
	infoFlags.BoolVar(&shellVars, "shell", false, "export shell variables")
}

// func runOpenURL(_ []string) error {
// 	wf.Configure(aw.TextErrors(true))
// 	log.Printf("opening URL %q ...", URL)
// 	_, err := util.RunCmd(exec.Command("open", URL))
// 	return err
// }

// search Firefox history
func runHistory(_ []string) error {
	checkForUpdate()
	if len(query) < 3 {
		wf.Warn("Query Too Short", "Please enter at least 3 characters")
		return nil
	}

	log.Printf("searching bookmarks for %q ...", query)
	history, err := mustClient().History(query)
	if err != nil {
		return err
	}

	custom := loadCustomActions()
	for _, h := range history {
		it := wf.NewItem(h.Title).
			Subtitle(h.URL).
			Arg(h.URL).
			UID(h.ID).
			Valid(true).
			Icon(iconHistory).
			Var("CMD", "url").
			Var("ACTION", urlDefault).
			Var("URL", h.URL).
			Var("TITLE", h.Title)

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

// search Firefox bookmarks
func runBookmarks(_ []string) error {
	checkForUpdate()
	if len(query) < 3 {
		wf.Warn("Query Too Short", "Please enter at least 3 characters")
		return nil
	}

	log.Printf("searching bookmarks for %q ...", query)
	bookmarks, err := mustClient().Bookmarks(query)
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

// search Firefox bookmarklets
func runBookmarklets(_ []string) error {
	checkForUpdate()
	if len(query) < 3 {
		wf.Warn("Query Too Short", "Please enter at least 3 characters")
		return nil
	}

	log.Printf("searching bookmarklets for %q ...", query)
	bookmarks, err := mustClient().Bookmarks(query)
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
			Copytext("bml:"+bm.ID+","+bm.Title).
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

// execute a bookmarklet in a tab
func runBookmarklet(_ []string) error {
	wf.Configure(aw.TextErrors(true))
	log.Printf("running bookmarklet %q in tab #%d ...", bookmarkID, tabID)

	return mustClient().
		RunBookmarklet(RunBookmarkletArg{BookmarkID: bookmarkID, TabID: tabID})
}

// filter open Firefox tabs
func runTabs(_ []string) error {
	log.Printf("fetching tabs for query %q ...", query)
	checkForUpdate()

	var (
		tabs []Tab
		err  error
	)
	if tabs, err = mustClient().Tabs(); err != nil {
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
			Icon(iconTab).
			Var("CMD", "tab").
			Var("ACTION", "Activate Tab").
			Var("TAB", id).
			Var("URL", t.URL).
			Var("TITLE", t.Title)

		it.NewModifier(aw.ModCmd).
			Subtitle("Other Actions…").
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

// execute a tab or URL action on the given tab
func runTabAction(_ []string) error {
	wf.Configure(aw.TextErrors(true))
	// load tab info so we can also run URL actions
	tab, err := mustClient().Tab(tabID)
	if err != nil {
		return err
	}

	log.Printf("running action %q on tab #%d ...", action, tab.ID)
	if a, ok := tabActions[action]; ok {
		return a.Run(tabID)
	}
	if a, ok := urlActions[action]; ok {
		return a.Run(tab.URL)
	}
	return fmt.Errorf("unknown action %q", action)
}

// run an action on a URL
func runURLAction(_ []string) error {
	_ = wf.Configure(aw.TextErrors(true))
	if URL == "" {
		tab, err := mustClient().Tab(0)
		if err != nil {
			return err
		}
		URL = tab.URL
	}
	log.Printf("running action %q on URL %q ...", action, URL)
	a, ok := urlActions[action]
	if !ok {
		return fmt.Errorf("unknown action %q", action)
	}
	return a.Run(URL)
}

// export variables containing info for currently-active tab
func runCurrentTabInfo(_ []string) error {
	_ = wf.Configure(aw.TextErrors(true))
	tab, err := mustClient().Tab(0)
	if err != nil {
		return err
	}
	if shellVars {
		fmt.Printf("export FF_TAB=%d\n", tab.ID)
		fmt.Printf("export FF_WINDOW=%d\n", tab.WindowID)
		fmt.Printf("export FF_INDEX=%d\n", tab.Index)
		fmt.Printf("export FF_TITLE=\"%s\"\n", tab.Title)
		fmt.Printf("export FF_URL=\"%s\"\n", tab.URL)
		return nil
	}
	av := aw.NewArgVars().
		Var("FF_TAB", fmt.Sprintf("%d", tab.ID)).
		Var("FF_WINDOW", fmt.Sprintf("%d", tab.WindowID)).
		Var("FF_INDEX", fmt.Sprintf("%d", tab.Index)).
		Var("FF_TITLE", tab.Title).
		Var("FF_URL", tab.URL)
	return av.Send()
}

// show actions for currently-active tab
func runCurrentTab(_ []string) error {
	tab, err := mustClient().Tab(0)
	if err != nil {
		return err
	}
	tabID = tab.ID
	URL = tab.URL
	return runActions([]string{})
}

// inject JavaScript into specified tab. If tabID is 0, JS in injected
// into the active tab.
func runInject(args []string) error {
	_ = wf.Configure(aw.TextErrors(true))
	if len(args) != 1 {
		return fmt.Errorf("inject command takes 1 argument, not %d", len(args))
	}
	js, err := mustClient().RunJS(RunJSArg{TabID: tabID, JS: args[0]})
	if err != nil {
		return err
	}
	fmt.Print(js)
	return nil
}

// filter actions for tab or URL
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
				Copytext("bml:"+a.id+","+a.name).
				Icon(actionIcon(a.name, iconBookmarklet)).
				Valid(true).
				Var("CMD", "run-bookmarklet").
				Var("BOOKMARK", a.id).
				Var("TAB", fmt.Sprintf("%d", tabID))
		}
	}

	if URL != "" {
		for _, a := range urlActions {
			if a.Name() == urlDefault {
				continue
			}
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

// check if a newer version of workflow is available
func runUpdate(_ []string) error {
	wf.Configure(aw.TextErrors(true))
	log.Print("checking for update ...")
	if err := wf.CheckForUpdate(); err != nil {
		return err
	}
	if wf.UpdateAvailable() {
		log.Println("a newer version of the workflow is available")
	}
	return nil
}

// show workflow status and options
func runStatus(_ []string) error {
	if c, err := newClient(); err != nil {
		wf.NewItem("No Connection to Browser").
			Subtitle(err.Error()).
			Icon(iconError)
	} else {
		if err := c.Ping(); err != nil {
			wf.NewItem("No Connection to Browser").
				Subtitle(err.Error()).
				Icon(iconError)

		} else {
			wf.NewItem("Connected to Browser").
				Subtitle("Extension is installed and running")
		}
	}

	wf.NewItem("Register Workflow with Browser").
		Subtitle("Use if you've updated or moved the workflow and it isn't working").
		Autocomplete("workflow:register").
		Icon(iconInstall).
		Valid(false)

	wf.NewItem("Install Browser Extension").
		Subtitle("Get the browser extension to integrate this workflow with Firefox").
		Arg(addonURL).
		Valid(true).
		Icon(iconAddon).
		Var("CMD", "url").
		Var("ACTION", urlDefault).
		Var("URL", addonURL)

	if wf.UpdateAvailable() {
		wf.NewItem("Update Available").
			Subtitle("↩ or ⇥ to install new version").
			Autocomplete("workflow:update").
			Icon(iconUpdateAvailable).
			Valid(false)
	} else {
		wf.NewItem("Workflow is Up to Date").
			Icon(iconUpdateOK).
			Valid(false)
	}

	dir := filepath.Join(wf.DataDir(), "scripts")
	wf.NewItem("Open Scripts Directory").
		Subtitle("Open custom scripts directory in Finder").
		Arg(dir).
		Valid(true).
		Icon(iconScript).
		Var("CMD", "url").
		Var("ACTION", "Open in Default Application").
		Var("URL", dir)

	wf.NewItem("Documentation").
		Subtitle("Open documentation in your browser").
		Arg(helpURL).
		Valid(true).
		Icon(iconDocs).
		Var("CMD", "url").
		Var("ACTION", urlDefault).
		Var("URL", docsURL)

	wf.NewItem("Report Issue").
		Subtitle("Open issue tracker in your browser").
		Arg(helpURL).
		Valid(true).
		Icon(iconIssue).
		Var("CMD", "url").
		Var("ACTION", urlDefault).
		Var("URL", helpURL)

	if query != "" {
		wf.Filter(query)
	}

	wf.WarnEmpty("No Matching Items", "Try a different query?")
	wf.SendFeedback()
	return nil
}

func runDownloads(_ []string) error {
	log.Printf("searching downloads for %q ...", query)
	downloads, err := mustClient().Downloads(query)
	if err != nil {
		return err
	}

	for _, dl := range downloads {
		wf.NewItem(filepath.Base(dl.Path)).
			Subtitle(util.PrettyPath(dl.Path)).
			Arg(dl.Path).
			UID(dl.Path).
			IsFile(true).
			Icon(&aw.Icon{Value: dl.Path, Type: aw.IconTypeFileIcon}).
			Valid(true).
			Var("CMD", "open").
			NewModifier(aw.ModCmd).
			Subtitle("Reveal in Finder").
			Var("CMD", "reveal")
	}

	wf.WarnEmpty("Nothing Found", "Try a different query?")
	wf.SendFeedback()
	return nil
}

// open file in default application
func runOpen(args []string) error {
	path := args[0]
	log.Printf("opening file %q ...", util.PrettyPath(path))
	return exec.Command("/usr/bin/open", path).Run()
}

// reveal file in Finder
func runReveal(args []string) error {
	path := args[0]
	log.Printf("revealing file %q in Finder ...", util.PrettyPath(path))
	return exec.Command("/usr/bin/open", "-R", path).Run()
}

// run update check in background
func checkForUpdate() {
	if wf.UpdateCheckDue() && !wf.IsRunning("update") {
		wf.RunInBackground("update", exec.Command(os.Args[0], "update"))
	}
	// TODO: show "update available" message
}
