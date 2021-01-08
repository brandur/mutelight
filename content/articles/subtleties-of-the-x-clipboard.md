+++
location = "Calgary"
published_at = 2011-08-03T07:50:00-06:00
slug = "subtleties-of-the-x-clipboard"
tiny_slug = "x-clipboard"
title = "Subtleties of the X Clipboard"
+++

A few months back, I started running the lean Archlinux build that I've been using through to today. I elected not to use a display manager, instead preferring booting [Awesome](http://awesome.naquadah.org/) (an aptly named tiling window manager) via the `startx` script. After getting used to Awesome's key bindings, and throwing Luakit, Urxvt and Tmux into the mix, I got about as close to an optimized Linux build as I was likely to get.

Everything was perfect ... except for one aspect: the clipboard. Its behavior was utterly perplexing: I could select text and middle-click (or `Shift-Insert`) it most places I wanted, but I could only copy _out_ of Chromium; while pasting it seemed to only respect text that had been copied from itself. Vim was even worse, even with `set clipboard=unnamed` it didn't seem to play nice with anything else.

This was pretty frustrating&mdash;the clipboard's importance in the everyday workflow really can't be understated. So what was the problem? To understand, we have to know a little more about the X clipboard.

The X Clipboard
---------------

In X10, **cut-buffers** were introduced. The concept behind this now obsolete mechanism was that when text was selected in a window, the window owner would copy it to a property of the root window called `CUT_BUFFER1`, a universally owned bucket available to every application on the system. General consensus on cut-buffers was that they were the absolute salt of the Earth, so a new system was devised.

Thus **selections** came about. Rather than applications copying data to a global bucket, they request ownership of a selection. When paste is called from another application, it requests data from the client that currently owns the selection. Aside from being much more versatile and less volatile than cut-buffers, selections can also be faster because no data has to be sent on a copy (only on paste). This is especially advantageous when there's a slow connection to the X server, but this strength is also a weakness because data made available by an application disappears when it closes.

Three selections are defined in the  <acronym title="Inter-Client Communication Conventions Manual">ICCCM</acronym>: `CLIPBOARD`, `PRIMARY`, and `SECONDARY`, each of which behaves like a clipboard in its own right:

* `CLIPBOARD`: traditionally used when text is copied and pasted from the edit menu, or via the `Ctrl+C` and `Ctrl+V` shortcuts in applications that support them.
* `PRIMARY`: traditionally used when a mouse selection is made, and pasted with middle-click or `Shift-Insert`.
* `SECONDARY`: ill-defined secondary selection. Most applications don't use it.

The heart of the problem for me is that I expected the X clipboard to behave like the clipboard on Windows or Mac OS X, but in fact X's architecture is fundamentally different with two separate, yet equally important, clipboards in use.

### Vim

Naturally, I had to know how Vim interacts with the X clipboard and was pleased to discover that it has some really great documentation on the subject (see for yourself with `:help x11-selection`). When running a GUI or X11-aware version of Vim, it has two registers that interact with X:

* `*` (as in `"*yy`): is the `PRIMARY` selection. `:set clipboard=unnamed` aliases it to the unnamed register.
* `+`: is the `CLIPBOARD` selection. `:set clipboard=unnamedplus` aliases it to the unnamed register.

Vim does not interact with the `SECONDARY` selection.

In Practice
-----------

I'm a Linux person at heart, but for me the two equal and separate selections remain an unfortunate usability problem.  Luckily for anyone with the same disposition, [Autocutsel](http://www.nongnu.org/autocutsel/) can help make X's behavior more logical and intuitive. It's a great little program that synchronizes the cut-buffer with `CLIPBOARD`, or both the cut-buffer and `CLIPBOARD` with `PRIMARY` as well.

Install Autocutsel (`pacman -S autocutsel` on Arch) and put the following two lines into your `.xinitrc` (or just run them from a terminal to immediately observe the effects):

```
autocutsel -fork &
autocutsel -selection PRIMARY -fork &
```

Now, no matter where you copy and paste from, be it `Ctrl+C` in Chrome, `p` in Vim, or through text selection in X, your clipboard is consistent across the entire system.
