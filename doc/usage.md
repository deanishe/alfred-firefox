Usage
=====

The workflow has the following keywords:

- `bm <query>` — Search Firefox bookmarks
  - `↩` — Open URL using default action
  - `⌘↩` — Show all URL actions
  - `...` — Run user-defined actions
- `bml <query>` — Search Firefox bookmarklets
  - `↩` — Run selected bookmarklet in active tab
  - `⌘C` — Copy bookmarklet ID & name to clipboard to set up a [custom tab action](bookmarklets.md)
- `tab [<query>]` — Filter tabs
  - `↩` — Activate tab
  - `⌘↩` — Show all tab & URL actions
  - `...` — Run user-defined action or bookmarklet
- `hist <query>` — Search Firefox history
  - `↩` — Open URL using default action
  - `⌘↩` — Show all URL actions
  - `...` — Run user-defined actions
- `dl [<query>]` — Search downloads
  - `↩` — Open downloaded file
  - `⌘↩` — Reveal downloaded file in Finder
  - `<your hotkey here>` — Show tab actions for active tab. You must assign your own Hotkey to use this very useful function.
- `ffass [<query>]` — Workflow status & setup
  - `Connected to Firefox` / `No Connection to Firefox` — Whether workflow can connect to Firefox
  - `Register Workflow with Firefox` — Install native app manifest so Firefox knows where to find the workflow. This should only be necessary if you've moved the workflow.
  - `Install Firefox Extension` — Get the Firefox extension for this workflow.
  - `Workflow is Up to Date` / `Update Available` — Whether a newer version of the workflow is available.
  - `Open Scripts Directory` — Open the custom scripts directory. You can add your own URL actions and/or icons here. See [Scripts](scripts.md) for details.
  - `Documentation` — Open these help pages in your browser.
  - `Report Issue` — Open the workflow's issue tracker in your browser.


See [Scripts](scripts.md) for more information on assigning custom hotkeys to URL actions and adding your own actions and icons.

See [Bookmarklets](bookmarklets.md) for more information on assigning custom hotkeys and icons to bookmarklets.

---

[^ Documentation index](index.md)

