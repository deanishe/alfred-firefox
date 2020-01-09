Bookmarklets
============

The workflow can run [bookmarklets][bookmarklets] from your Firefox bookmarks in arbitrary tabs.

You can use the keyword `bml` to filter your bookmarklets and run them in the active tab, but you can also assign them to specific Hotkeys to run them on arbitrary tabs when filtering tabs (keyword `tab`) or add them to the `Other Actions…` list for tabs.

Bookmarklets are assigned to alternative hotkeys in a similar way to [scripts][scripts], but using the form `TAB_<KEY>` instead of `URL_<KEY>`, and the variables set in [Alfred's workflow configuration sheet][config-sheet] must follow a specific format:

| Variable name | Variable value | Effect |
|--|--|--|
| `TAB_OPT` | `bml:seoxED9MBuqi,Add to Pinboard` | Run your `Add to Pinboard` bookmarklet in the selected tab |
| `TAB_CMD_SHIFT` | `bml:uFJCA875bTEt,Open in GoDoc` | Run your `Open in GoDoc` bookmarklet in the selected tab |

The value for a bookmarklet *must* have the form `bml:<bookmark ID>,<name>`. You can get the appropriate value by hitting `⌘C` (copy) on a bookmarklet in the bookmarklet list (default keyword: `bml`).

You can change the name (what's after the first comma) to anything you want, and you can assign a custom icon to the bookmarklet by adding an image file with the same basename (i.e. without file extension) as the name you've assigned to the bookmarklet to the `scripts` directory. So for example, you could set a custom icon for the above `Open in GoDoc` bookmarklet by putting a file `Open in GoDoc.png` in the `scripts` subdirectory of the workflow's data directory. See [Script icons][script-icons] for details.


Bookmarklets without hotkeys
----------------------------

You can also add bookmarklets to the `Other Actions…` list without assigning them a hotkey. To do this add a workflow variable with a name of the form `TAB_<BLAH>` with a valid value for a bookmarklet (as described above). If `<BLAH>` is not a valid combination of hotkeys, the bookmarklet will be shown in the `Other Actions…` list, but will not be assigned to an alternative hotkey:

| Variable name | Variable value | Effect |
|--|--|--|
| `TAB_PINBOARD` | `bml:seoxED9MBuqi,Add to Pinboard` | Add `Add to Pinboard` action to `Other Actions…` list |
| `TAB_GODOC` | `bml:uFJCA875bTEt,Open in GoDoc` | Add `Open in GoDoc` action to `Other Actions…` list |


---

[^ Documentation index](index.md).


[config-sheet]: https://www.alfredapp.com/help/workflows/advanced/variables/#environment
[bookmarklets]: https://en.wikipedia.org/wiki/Bookmarklet
[scripts]: scripts.md
[script-icons]: scripts.md#script-icons

