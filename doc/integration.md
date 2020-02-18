Integration
===========

Because Firefox lacks AppleScript support, it is rarely supported by workflows that retrieve the current page from browsers. This workflow provides an External Trigger than can be used by other workflows to get the title and URL of the current Firefox tab.

<!-- vim-markdown-toc GFM -->

* [Variables](#variables)
* [Write tab info to a file](#write-tab-info-to-a-file)
* [Pass tab info to External Trigger](#pass-tab-info-to-external-trigger)

<!-- vim-markdown-toc -->

The External Trigger is called `current-tab` and Firefox Assistant's bundle ID is `net.deanishe.alfred.firefox-assistant`, so the basic AppleScript code to call the trigger is:

```applescript
tell application id "com.runningwithcrayons.Alfred" to run trigger "current-tab" in workflow "net.deanishe.alfred.firefox-assistant"
```

The trigger works in two distinct ways: by writing the tab info to a file or by calling a second External Trigger in your workflow with the tab info in variables.

Check out [the demo workflow][demo] to see implementations of both methods.


Variables
---------

The variables set by the trigger are:

| Variable    | Description                            |
| ----------- | -------------------------------------- |
| `FF_TITLE`  | Tab title                              |
| `FF_URL`    | Tab URL                                |
| `FF_TAB`    | Tab ID                                 |
| `FF_WINDOW` | Window ID                              |
| `FF_INDEX`  | Index of tab in window (0 = first tab) |


Write tab info to a file
------------------------

Pass the path of the target file as the argument to the External Trigger:

```applescript
tell application id "com.runningwithcrayons.Alfred" to run trigger "current-tab" in workflow "net.deanishe.alfred.firefox-assistant" with argument "/path/to/file.json"
```

This tells Firefox Assistant to write the tab info to `/path/to/file.json`.

**Note:** External Triggers are run asynchronously, so your calling code must wait some time (a second should suffice) for the workflow to write the file.

You can use this AppleScript in a Run Script action in your workflow to tell Firefox Assistant to write the current tab info to a file, wait for the file to appear, and finally set the tab data as workflow variables:

```applescript
(*
	Tell Firefox Assistant to write tab info as JSON to a file.
	The output is written as workflow JSON, so writing it to STDOUT will tell Alfred
	to set the tab info as workflow variables:

	FF_URL   - URL of active tab
	FF_TITLE - Title of active tab
	etc.
*)

-- Write to file in current workflow's cache directory
property _dir: (system attribute "alfred_workflow_cache")
property _path: _dir & "/firefox_current_tab.json"
property _file: POSIX file _path

on run argv
	-- Ensure directory exists
	do shell script "mkdir -p " & quoted form of _dir
	-- Ensure JSON file does not
	do shell script "rm -f " & quoted form of _path

	-- Tell Firefox Assistant to write tab info to _path
	tell application id "com.runningwithcrayons.Alfred"
		run trigger "current-tab" in workflow "net.deanishe.alfred.firefox-assistant" with argument _path
	end tell

	-- External Trigger calls are asynchronous, so wait up to one second for file to appear,
	-- then return contents of file (which are written to STDOUT)
	tell application "Finder"
		repeat 100 times
			try
				if exists _file then
					set _content to (read _file)
					-- Remove JSON file
					do shell script "rm -f " & quoted form of _path
					-- Return Alfred JSON to set tab info as workflow variables
					return _content
				end if
			end try
			delay 0.01
		end repeat
	end tell

	-- Failed to get tab info (in time)
	error "tab info not saved by Firefox Assistant"
end run
```

Downstream workflow elements can access the tab title, URL etc. via [the variables](#variables).

**Note:** Although the script waits for up to a second for the file to appear, this shouldn't slow your workflow down if Firefox Assistant isn't installed, as the script will fail when it tries to call the External Trigger.


Pass tab info to External Trigger
---------------------------------

Alternatively, the External Trigger can be told to pass the tab info back to a second External Trigger in your (or a third-party) workflow.

To do this, you must set the workflow variables `BUNDLE_ID` and `TRIGGER` before calling the External Trigger **and you must not specify an arg**.

`BUNDLE_ID` should be the bundle ID of your workflow and `TRIGGER` is the name of the External Trigger in your workflow.

Firefox Assistant will pass [the variables](#variables) back to your workflow's External Trigger. This method is faster than writing to a file, but often not as convenient.

See [the demo workflow][demo] for an example implementation.


---

[^ Documentation index](index.md)


[demo]: https://github.com/deanishe/alfred-firefox/raw/master/doc/Firefox%20Trigger%20Demo.alfredworkflow

