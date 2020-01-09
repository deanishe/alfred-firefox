Scripts
=======

The workflow can do arbitrary things with URLs via scripts. Most of the built-in URL Actions are implemented via scripts in the internal `scripts` directory, and you can add your own scripts (with optional icons) to extend the workflow's functionality.

Place your custom scripts in the `scripts` subdirectory of the workflow's data directory (which can be quickly accessed via the `ffass` keyword and `Open Scripts Directory` item). **Do not add your own scripts to the workflow's internal `scripts` directory**: they'll be removed when you update the workflow.


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


Script icons
------------

You can optionally assign a custom icon to a script by putting an image file with the same basename (i.e. excluding extension) in the `scripts` directory. Icons of type PNG, GIF, JPG and ICNS are supported.

For example, if your script is called `Add to Pinboard.py`, you can assign it a custom icon by putting a file called `Add to Pinboard.png` (or `Add to Pinboard.icns` etc.) in the `scripts` directory.


Bookmarklets
------------

The workflow can also run bookmarklets (read from your Firefox bookmarks), but these work slightly differently to scripts. See [bookmarklets][bookmarklets] for details.


---

[^ Documentation index](index.md).


[script-types]: https://godoc.org/github.com/deanishe/awgo/util#Runner
[config-sheet]: https://www.alfredapp.com/help/workflows/advanced/variables/#environment
[bookmarklets]: bookmarklets.md

