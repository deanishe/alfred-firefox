Scripts
=======

The workflow can do arbitrary things with URLs via scripts. Most of the built-in URL Actions are implemented via scripts in the internal `scripts` directory, and you can add your own scripts (with optional icons) to extend the workflow's functionality.

Place your custom scripts in the `scripts` subdirectory of the workflow's data directory (which can be quickly accessed via the `ffass` keyword and `Open Scripts Directory` item). **Do not add your own scripts to the workflow's internal `scripts` directory**: they'll be removed when you update the workflow.


<!-- vim-markdown-toc GFM -->

* [How URL scripts work](#how-url-scripts-work)
  * [Current browser](#current-browser)
  * [Script icons](#script-icons)
* [Advanced scripting](#advanced-scripting)
  * [Running actions](#running-actions)
  * [Getting tab information](#getting-tab-information)
  * [Injecting JavaScript](#injecting-javascript)
* [Bookmarklets](#bookmarklets)

<!-- vim-markdown-toc -->


How URL scripts work
--------------------

Scripts are called with the URL (from a bookmark, tab etc.) as the first command-line argument (i.e. `$1`). A script can be an executable file or any of the script types (determined via file extension) [supported by AwGo][script-types].

Scripts can be run via the `Other Actions…` menu (`⌘↩` on a bookmark/tab) or directly via alternative hotkeys (e.g. `^↩`) by setting an appropriate variable [in the workflow's configuration sheet][config-sheet]. To assign a hotkey to a script, create a workflow variable of the form `URL_<KEY>` with the script's name (without extension) as the value, e.g.:

| Variable name | Variable value | Effect |
|---------------|----------------|--------|
| `URL_CTRL`    | `Open in Safari`| Open URL in Safari when `^↩` is pressed on a result |
| `URL_OPT_CTRL`    | `Open in Chrome`| Open URL in Google Chrome when `^⌥↩` is pressed on a result |

Modifier keys (`OPT`, `CMD`, `SHIFT`, `CTRL`) can be arbitrarily combined by joining them with underscores, e.g.: `URL_OPT_SHIFT_CMD`, `URL_CTRL_OPT_SHIFT` etc.).

You can quickly grab the name of a script by using `CMD+C` (copy) on an action in the `Other Actions…` list, which will copy the script's name to the clipboard.


### Current browser ###

As the workflow supports different versions of Firefox (Firefox, Firefox Nightly, Firefox Developer Edition), the name of the application it's currently connected to will be specifed in the `BROWSER` environment variable. The default `Open in Firefox.sh` script, for example, uses the command

```bash
/usr/bin/open -a "$BROWSER" "$1"
```

to open the URL in the browser the URL came from.


### Script icons ###

You can optionally assign a custom icon to a script by putting an image file with the same basename (i.e. excluding extension) in the `scripts` directory. Icons of type PNG, GIF, JPG and ICNS are supported.

For example, if your script is called `Add to Pinboard.py`, you can assign it a custom icon by putting a file called `Add to Pinboard.png` (or `Add to Pinboard.icns` etc.) in the `scripts` directory.


Advanced scripting
------------------

In addition to URL scripts, you can also call the `alfred-firefox` binary directly in order to run a script by name, get info for the current tab or [inject arbitrary JavaScript](#injecting-javascript) into the current tab.

As the workflow directory is the working directory, you can call the workflow binary with `./alfred-firefox`.

You can also call it from within your own URL action scripts.

There are a couple of example actions in the workflow, one showing how to run a bookmarklet via a Hotkey and another showing how to run an action via a Hotkey.


### Running actions ###

You can execute an action by name by running:

```bash
./alfred-firefox -action "Name of Action" tab
```

The name of the action should be as shown in the workflow UI, i.e. the base name of the script file without the file extension. You can run any action known to the workflow (tab or URL action, built-in or user-added).


### Getting tab information ###

The `alfred-firefox tab-info` command outputs data about the active tab (title, URL, tab ID etc.). By default, it emits Alfred JSON to set workflow variables (suitable for executing in a Run Script action). Passing the `-shell` flag emits shell script suitable for `eval`:

```bash
eval $( ./firefox-alfred tab-info -shell )
echo $FF_TITLE
echo $FF_URL
```

The exported variables are:

| Variable    | Description                            |
| ----------- | -------------------------------------- |
| `FF_TITLE`  | Tab title                              |
| `FF_URL`    | Tab URL                                |
| `FF_TAB`    | Tab ID                                 |
| `FF_WINDOW` | Window ID                              |
| `FF_INDEX`  | Index of tab in window (0 = first tab) |


### Injecting JavaScript ###

You can also inject JS into a tab by running:

```bash
./alfred-firefox inject '... your JS code here ...'
```

For example, add a script called `Alert.sh` to the scripts directory (keyword `ffass`, choose `Open Scripts Directory`) with the following contents:

```bash
./alfred-firefox inject 'alert("Hello!")'
```

Now you have a new `Alert` action available that will pop up an alert in the selected tab.

**Note:** Such scripts will also be available in bookmark menus, but running the script from there will always inject the JS into the active tab.

The `inject` command outputs the return value of the executed JavaScript as JSON, so you can use this feature to extract information from webpages:

```bash
#!/bin/bash
./alfred-firefox inject 'Array.from(document.querySelectorAll("h2")).map(el => el.innerText)'
```

Outputs (for [this page][execute-docs]):

```json
[["On this Page","Syntax","Examples","Browser compatibility","Thank you!","Tell us what’s wrong with this table","Learn the best of web development","MDN","Mozilla"]]
```

**Note:** The result is *always* an array. See [the documentation for the `tabs.executeScript` API][execute-docs] for more information.


Bookmarklets
------------

The workflow can also run bookmarklets (read from your Firefox bookmarks), but these work slightly differently to scripts. See [bookmarklets][bookmarklets] for details.


---

[^ Documentation index](index.md)


[script-types]: https://godoc.org/github.com/deanishe/awgo/util#Runner
[config-sheet]: https://www.alfredapp.com/help/workflows/advanced/variables/#environment
[bookmarklets]: bookmarklets.md
[execute-docs]: https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/tabs/executeScript#Return_value
