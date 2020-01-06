// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/peterbourgon/ff/ffcli"
)

var (
	// starts extension client & RPC server
	serveCmd = &ffcli.Command{
		Name:      "serve",
		Usage:     "firefox serve",
		ShortHelp: "start extension server (called by Firefox)",
		LongHelp: wrap(`
			Run extension server. This is called by the Firefox
			extension and provides and RPC server for the workflow
			to call into Firefox.
		`),
		Exec: runServer,
	}
)

// set up logging for the server.
// doesn't use the same log as the rest of the workflow, as this is
// a long-running process, and we don't want the log file it's using
// being rotated by another process
func initLogging() error {
	if fi, err := os.Stat(logfile); err == nil {
		if fi.Size() > int64(1048576) {
			_ = os.Rename(logfile, logfile+".1")
		}
	}
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	multi := io.MultiWriter(file, os.Stderr)
	log.SetOutput(multi)
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("[%d] ", os.Getpid()))
	return nil
}

// return PID of running server or 0.
func getPID() int {
	if data, err := ioutil.ReadFile(pidFile); err == nil {
		if pid, err := strconv.Atoi(string(data)); err == nil {
			return pid
		}
	}
	return 0
}

// return true if process with PID is running.
func processRunning(pid int) bool {
	if err := syscall.Kill(pid, 0); err == nil {
		return true
	}
	return false
}

// write PID file, terminating and waiting for existing server if it exists.
func writePID() error {
	pid := getPID()
	if pid != 0 {
		log.Printf("signalling existing server %d to stop ...", pid)
		_ = syscall.Kill(pid, syscall.SIGTERM)

		start := time.Now()
		for processRunning(pid) {
			if time.Now().Sub(start) > time.Second*2 {
				return fmt.Errorf("server already running")
			}
			time.Sleep(time.Millisecond * 100)
		}
	}

	return ioutil.WriteFile(pidFile, []byte(strconv.FormatInt(int64(os.Getpid()), 10)), 0600)
}

// start extension client and RPC server
func runServer(args []string) error {
	wf.Configure(aw.TextErrors(true))
	if err := writePID(); err != nil {
		return err
	}
	if err := initLogging(); err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(socketPath)
		_ = os.Remove(pidFile)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	f := newFirefox()
	go f.run()

	srv, err := newRPCService(socketPath, f)
	if err != nil {
		return err
	}
	go srv.run()

	var s string
	if err := srv.Ping("", &s); err != nil {
		log.Printf("[error] %v", err)
	} else {
		log.Printf("ping => %q", s)
	}

	/*
		var bookmarks []Bookmark
		for _, q := range []string{"", "haze", "github", "p2p"} {
			if err := srv.Bookmarks(q, &bookmarks); err != nil {
				log.Printf("[error] %v", err)
			} else {
				log.Printf("bookmarks(%q) => %d result(s)", q, len(bookmarks))
			}
		}

		var tabs []Tab
		if err := srv.Tabs("", &tabs); err != nil {
			log.Printf("[error] %v", err)
		} else {
			log.Printf("tabs => %d result(s)", len(tabs))
		}
	*/

	<-quit
	log.Print("shutting down ...")
	f.stop()
	srv.stop()
	return nil
}
