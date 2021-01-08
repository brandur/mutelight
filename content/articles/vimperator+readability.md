+++
location = "Calgary"
published_at = 2010-07-16T00:00:00-06:00
slug = "vimperator+readability"
tiny_slug = "26"
title = "Vimperator + Readability"
+++

* *[Vimperator](http://vimperator.org/)* &mdash; gives Firefox Vim bindings and enables keyboard-based web browsing.

* *[Readability](http://lab.arc90.com/experiments/readability/)* &mdash; reformats hard-to-read web pages (all too common these days) for sane consumption.

Normally, Readability is dragged to your browser's bookmarks toolbar and is used from there. In Vimperator, bookmarks are not too convenient to access, so it's useful to add a quickmark so that you can activate Readability in three quick keystrokes.

1. Go to [Readability](http://lab.arc90.com/experiments/readability/) and select your settings

2. Copy the Javascript used to create your bookmarklet by copying the target of the big _Readability_ button (`;y<link number`> in Vimperator)

3. Open `~/vimperator/info/default/quickmarks`

4. Add `,"r":"<copied bookmarklet code>"` to the end of the hash enclosed by `{...}` (replace `r` with the character of your choice)

5. Restart Firefox

6. Find a badly formatted page and type `gor`
