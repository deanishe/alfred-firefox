// Copyright (c) 2019 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

// Command firefox is an Alfred workflow to interact with Firefox.
package main

// TODO: implement setup (save native application manifest)
// TODO: package extension

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/deanishe/awgo/util"
	"github.com/mitchellh/go-wordwrap"
	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
)

const (
	maxCacheAge = time.Minute * 30
	wrapWidth   = 72
)

const (
	helpURL  = "https://github.com/deanishe/alfred-firefox/issues"
	docsURL  = "https://github.com/deanishe/alfred-firefox/blob/master/doc/index.md"
	addonURL = "https://addons.mozilla.org/en-US/firefox/addon/alfred-launcher-integration/"
	repo     = "deanishe/alfred-firefox"
)

// native application manifest
var (
	extensionID   = "alfredfirefox@deanishe.net"
	extensionName = "net.deanishe.alfred.firefox"
	manifestPath  = os.ExpandEnv("${HOME}/Library/Application Support/Mozilla/" +
		"NativeMessagingHosts/" + extensionName + ".json")
)

// workflow variables
var (
	wf = aw.New(
		aw.HelpURL(helpURL),
		update.GitHub(repo),
		aw.AddMagic(registerMagic{}),
	)

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
		Usage:     "alfred-firefox <command> [flags] [args...]",
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
	rootFlags.StringVar(&urlDefault, "url-default", "Open in Default Application",
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
		downloadsCmd,
		historyCmd,
		openCmd,
		revealCmd,
		runBookmarkletCmd,
		serveCmd,
		statusCmd,
		tabCmd,
		tabsCmd,
		urlCmd,
		updateCmd,
	}
	pidFile = filepath.Join(wf.CacheDir(), "server.pid")
	logfile = filepath.Join(wf.CacheDir(), fmt.Sprintf("%s.server.log", wf.BundleID()))
	u, _ := user.Current()
	socketPath = fmt.Sprintf("/tmp/alfred-firefox.%s.sock", u.Uid)
}

func run() {
	for _, dir := range scriptDirs {
		if err := os.MkdirAll(dir, 0700); err != nil {
			panic(err)
		}
	}

	if err := setup(false); err != nil {
		panic(err)
	}

	if err := loadURLActions(); err != nil {
		panic(err)
	}

	if err := rootCmd.Run(wf.Args()); err != nil {
		panic(err)
	}
}

func main() { wf.Run(run) }

// Magic Action to install native application manifest in Firefox
type registerMagic struct{}

func (m registerMagic) Keyword() string { return "register" }
func (m registerMagic) Description() string {
	return "Re-register workflow with Firefox"
}
func (m registerMagic) RunText() string {
	return "Registered. Re-open Firefox extension to connect."
}
func (m registerMagic) Run() error { return setup(true) }

var _ aw.MagicAction = registerMagic{}

func setup(force bool) error {
	if !force && util.PathExists(manifestPath) {
		return nil
	}

	path, err := filepath.Abs("./server.sh")
	if err != nil {
		return err
	}
	if path, err = filepath.EvalSymlinks(path); err != nil {
		return err
	}
	path = filepath.Clean(path)

	manifest := struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Path        string   `json:"path"`
		Type        string   `json:"type"`
		Allowed     []string `json:"allowed_extensions"`
	}{
		Name:        extensionName,
		Description: "Alfred plugin for Firefox",
		Path:        path,
		Type:        "stdio",
		Allowed:     []string{extensionID},
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(manifestPath, data, 0644); err != nil {
		return err
	}
	log.Printf("wrote native app manifest to %q", util.PrettyPath(manifestPath))
	log.Print("\n" + string(data))
	return nil
}

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
