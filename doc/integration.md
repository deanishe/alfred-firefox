Integration
===========

Because Firefox lacks AppleScript support, it is rarely supported by workflows that retrieve the current page from browsers. This workflow provides an External Trigger than can be used by other workflows to get the title and URL of the current Firefox tab.

<!-- MarkdownTOC autolink="true" bracket="round" levels="2,3,4,5" autoanchor="true" -->

- [Variables](#variables)
- [Write tab info to disk](#write-tab-info-to-disk)
- [Pass tab info to External Trigger](#pass-tab-info-to-external-trigger)

<!-- /MarkdownTOC -->

The External Trigger is called `current-tab` and Firefox Assistant's bundle ID is `net.deanishe.alfred.firefox-assistant`, so the basic AppleScript code to call the trigger is:

```applescript
tell application id "com.runningwithcrayons.Alfred" to run trigger "current-tab" in workflow "net.deanishe.alfred.firefox-assistant"
```

The trigger works in two distinct ways: by writing the tab info to disk or by calling a second External Trigger in your workflow with the tab info in variables.

Check out [the demo workflow][demo] to see implementations of both methods.


<a id="variables"></a>
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


<a id="write-tab-info-to-disk"></a>
Write tab info to disk
----------------------

Pass the path of the target file as the argument to the External Trigger:

```applescript
tell application id "com.runningwithcrayons.Alfred" to run trigger "current-tab" in workflow "net.deanishe.alfred.firefox-assistant" with argument "/path/to/file.json"
```

This tells Firefox Assistant to write the tab info to `/path/to/file.json`.

**Note:** Because External Triggers are run asynchronously, creating a named pipe for the workflow to write to is the best solution. That way, you can immediately start reading the response from the pipe, and it will block until the workflow writes the tab data to it.

You can use this shell script in a Run Script action in your workflow to tell Firefox Assistant to write the current tab info to a file and then set the provided tab data as workflow variables:

```bash
#!/bin/zsh

set -e  # exit on error

# ensure workflow cache directory exists
mkdir -p "$alfred_workflow_cache"

# create a named pipe for the workflow to write to
pipe="${alfred_workflow_cache}/tab.$$.json"
mkfifo -m 0600 "$pipe"
# ensure pipe is deleted when script exits
trap "rm -f '$pipe'" EXIT INT TERM

# tell workflow to write tab deets (as Alfred JSON) to pipe
script="tell application id \"com.runningwithcrayons.Alfred\" to run trigger \"current-tab\" in workflow \"net.deanishe.alfred.firefox-assistant\" with argument \"$pipe\""
/usr/bin/osascript -e "$script" &!
# send JSON from pipe and to Alfred
cat "$pipe"
```

Downstream workflow elements can access the tab title, URL etc. via [the variables](#variables).

**Note:** If Firefox Assistant isn't installed, the script will fail and exit when it tries to call the External Trigger, so it won't block trying to read from the pipe.


<a id="pass-tab-info-to-external-trigger"></a>
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

