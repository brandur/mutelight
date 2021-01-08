+++
location = "San Francisco"
published_at = 2012-08-19T08:43:46-07:00
slug = "caps-lock"
title = "Caps Lock + Tmux"
+++

Nothing. But beyond that, without rebinding, the location and function of caps lock on a modern keyboard is more than inconvenient -- it's actively destructive. Enabling caps lock during a Vim session in normal code causes all hell to break loose, and what's going on isn't obvious unless you've seen the symptoms before.

I was pleasantly surprised to find that in OSX you can now disable caps lock out of the box if you don't intend to rebind it. This is accomplished via `System Preferences --> Keyboard --> Modifier Keys --> Caps Lock Key --> No Action`, and provides a measurable improvement over the system default.

Then I started to wonder whether I could put caps lock to good use by solving another problem caused by Apple's keyboard design, and the answer turned out to be yes.

Caps Lock as Tmux Prefix
------------------------

Tmux has moved beyond a terminal multiplexing tool and has become one of the most important tools in my kit by acting as the de facto window manager for all my important tools and sessions. As such, I hit my Tmux prefix shortcut `C-a` _a lot_, which is tremendously inconvenient because even in 2012 Apple is still jamming a `fn` key onto everything they make so that `ctrl` is harder to hit.

Switching to caps lock as a Tmux prefix solves this problem forever. Here's how to do it:

1. Download [PCKeyboardHack](http://pqrs.org/macosx/keyremap4macbook/pckeyboardhack.html.en), install, and restart.
2. From the new System Preferences pane, change the keycode under the `Change Caps Lock` entry to `109` (that's `F10`), and check its box.
3. In your `.tmux.conf`, change applicable settings to use `F10`:

    ``` bash
    # thanks to PCKeyboardHack, F10 is caps lock and caps lock is F10
    set-option -g prefix F10

    # go to last window by hitting caps lock two times in rapid succession
    bind-key F10 last-window
    ```
