// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

// TODO: implement history search
// TODO: implement download search
// TODO: implement download scripts

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/mitchellh/go-wordwrap"
	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
)

const (
	maxCacheAge = time.Minute * 30
	wrapWidth   = 72
)

var (
	wf = aw.New()

	// Filepaths
	scriptDirs = []string{
		filepath.Join(wf.Dir(), "scripts"),
		filepath.Join(wf.DataDir(), "scripts"),
	}
	socketPath string
	pidFile    string
	logfile    string

	// CLI flags/environment variables
	URL        string
	urlDefault string
	tabID      int
	action     string
	bookmarkID string
	query      string

	rootFlags = flag.NewFlagSet("firefox", flag.ExitOnError)
	rootCmd   = &ffcli.Command{
		Usage:     "firefox <command> [flags] [args...]",
		ShortHelp: "Firefox workflow for Alfred",
		LongHelp: wrap(`
		Alfred workflow to control Firefox.
		You must also install the Firefox extension for this workflow to work.
		`),
		FlagSet: rootFlags,
		Options: []ff.Option{ff.WithEnvVarNoPrefix()},
		Exec: func(args []string) error {
			return flag.ErrHelp
		},
	}
)

func init() {
	rootFlags.StringVar(&URL, "url", "", "URL")
	rootFlags.StringVar(&urlDefault, "url-default", "Open in Firefox",
		"Default URL action")
	rootFlags.IntVar(&tabID, "tab", 0, "ID of tab")
	rootFlags.StringVar(&bookmarkID, "bookmark", "", "ID of bookmark")
	rootFlags.StringVar(&query, "query", "", "Search query")
	rootFlags.StringVar(&action, "action", "", "Action name")

	rootCmd.Subcommands = []*ffcli.Command{
		actionsCmd,
		bookmarkletsCmd,
		bookmarksCmd,
		currentTabCmd,
		openURLCmd,
		runBookmarkletCmd,
		serveCmd,
		tabCmd,
		tabsCmd,
		urlCmd,
	}
	pidFile = filepath.Join(wf.CacheDir(), "server.pid")
	logfile = filepath.Join(wf.CacheDir(), fmt.Sprintf("%s.server.log", wf.BundleID()))
	socketPath = filepath.Join(os.Getenv("HOME"), ".alfred-firefox.sock")
}

func run() {
	for _, dir := range scriptDirs {
		if err := os.MkdirAll(dir, 0700); err != nil {
			panic(err)
		}
	}
	if err := loadURLActions(); err != nil {
		panic(err)
	}

	if err := rootCmd.Run(wf.Args()); err != nil {
		panic(err)
	}
}

func main() { wf.Run(run) }

var rxPara = regexp.MustCompile(`\n\n+`)

func wrap(text string) string {
	paras := rxPara.Split(text, -1)

	for i, s := range paras {
		var b strings.Builder
		scanner := bufio.NewScanner(strings.NewReader(s))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			b.WriteString(line + " ")
		}
		paras[i] = wordwrap.WrapString(b.String(), wrapWidth)
	}

	return strings.TrimSpace(strings.Join(paras, "\n\n"))
}
