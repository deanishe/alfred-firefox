<div align="center">
  <img src="https://github.com/deanishe/alfred-firefox/blob/master/icon.png" alt="Alfred-Firefox icon" title="Alfred-Firefox icon"/>
</div>

Firefox Assistant for Alfred
============================

Search and manipulate Firefox's bookmarks, history and tabs from Alfred.

![Animated demo of workflow in use][demo]

The workflow can be easily [extended with your own actions][scripts].

Installation
------------

The workflow supports Alfred 4+ and the extension works with (at least) Firefox, Firefox Nightly and Firefox Developer Edition.

1. Download and install the [latest version of the workflow][workflow].
2. Run `ffass` in Alfred and choose `Install Firefox Extension` to get [the Firefox extension][addon].

See [the setup documentation][setup] for more details.


Usage
-----

The basic usage is:

- `bm <query>` — Search bookmarks
- `bml <query>` — Search bookmarklets
- `hist <query>` — Search history
- `dl [<query>]` — Search downloads
- `tab [<query>]` — Search tabs
- `ffass [<query>]` — Workflow status and links

See [the usage documentation][usage] for full details.


Integration
-----------

The workflow can be used by other workflows to retrieve the title and URL of the active Firefox tab (in lieu of AppleScript, which Firefox doesn't support). See [the integration docs][integration] for details.


Documentation
-------------

See [the full documentation][docs] for detailed info on setting up, using and customising the workflow.


Licensing & thanks
------------------

This workflow and extension are released under the [MIT licence][licence].

It is written in [Go][go] and heavily based on the [AwGo library][awgo]. The icons are based on [Font Awesome][fontawesome].


[addon]: https://addons.mozilla.org/en-US/firefox/addon/alfred-launcher-integration/
[licence]: https://github.com/deanishe/alfred-firefox/blob/master/LICENCE.txt
[workflow]: https://github.com/deanishe/alfred-firefox/releases/latest
[demo]: https://github.com/deanishe/alfred-firefox/blob/master/demo.gif
[docs]: https://github.com/deanishe/alfred-firefox/blob/master/doc/index.md
[scripts]: https://github.com/deanishe/alfred-firefox/blob/master/doc/scripts.md
[integration]: https://github.com/deanishe/alfred-firefox/blob/master/doc/integration.md
[usage]: https://github.com/deanishe/alfred-firefox/blob/master/doc/usage.md
[setup]: https://github.com/deanishe/alfred-firefox/blob/master/doc/setup.md
[go]: https://golang.org
[awgo]: https://github.com/deanishe/awgo
[fontawesome]: https://fontawesome.com/

