+++
location = "Calgary"
published_at = 2011-03-28T00:00:00-06:00
slug = "practical-tmux"
tiny_slug = "42"
title = "Practical Tmux"
+++

Today I switched over completely from GNU Screen to the more modern BSD-licensed alternative, **tmux**. After making sure that tmux had replacements for all Screen's key features, I took the plunge, and haven't looked back. The [project's webpage](http://tmux.sourceforge.net/) has a complete list of features available under tmux, but as an everyday user of Screen, here are the major reasons I switched:

* **Better redraw model:** I use Awesome WM, a tiling window manager, and terminals containing Screen sessions would glitch out regularly; spewing all kinds of artifacts into their windows. Sufficed to say, tmux doesn't do this.

* **Screen contents persisted through full-screen programs:** in Screen, you lose your terminal's previous contents after leaving a full-screen program like an editor. Tmux doesn't have this problem.

* **Rational configuration:** I once tried to configure my screen's status line, and eventually just gave up. In comparison, tmux's lines like `set -g status-right "#[fg=green]#H` are almost a little _too_ easy. This goes for other configuration values as well.

* **Visual bell that works:** one of the only things under Linux that's come close to driving me completely crazy is the line _Wuff -- Wuff!!_. I mean, I'm all good with programmers having a sense of humour, but this is just too much. Even after disabling the visual bell and doing away with this default message, it's tragically not possible to remove the visual bell from Screen completely.

* **Automatic window renaming:** windows are renamed automatically to the command running in them unless their name has been manually changed using `C-a ,`.

* **Vertical splits:** it's always been a mystery that Screen can do horizontal screen splits but not vertical without fancy patches. Tmux does both out of the box.

* **VI key bindings in copy mode:** VI or Emacs keys are available upon entering tmux's copy mode.

* **Runtime configuration:** you can easily open a prompt in tmux to apply configuration to a running session.

Common Problems and Solutions
-----------------------------

All that said, unfortunately tmux isn't completely convention over configuration, and it took a bit of work to get running exactly how I wanted. The purpose of this post is to go over some common tmux problems and their solutions.

### C-b

Now there's a prefix that you have to reach for! Fix this in `~/.tmux.conf` by changing it to `C-a`:

```
set-option -g prefix C-a
```

### C-a C-a for the Last Active Window

This was a feature in Screen that was great enough to keep around. Add the following to `~/.tmux.conf`:

```
bind-key C-a last-window
```

### Command Sequence for Nested Tmux Sessions

Often I'll run a multiplexer inside another multiplexer and need a command sequence to send things to the inner session. In Screen, this could be accomplished using `C-a a <command>`. This doesn't work out of the box in tmux, but can be fixed with a little configuration.

```
bind-key a send-prefix
```

### Start Window Numbering at 1

Zero-based indexing is sure great in programming languages, but not so much in terminal multiplexers where that zero is all the way on the other side of the keyboard.

```
set -g base-index 1
```

### Faster Command Sequences

Upon starting to use tmux, I noticed that I had to add a noticeable delay between two characters in a command sequence for it to recognize the command, for example between the `C-a` and `n` in `C-a n`. This is because tmux is waiting for an escape sequence. Fix that by setting escape time to zero.

```
set -s escape-time 0
```

### Aggressive Resize

By default, all windows in a session are constrained to the size of the smallest client connected to that session, even if both clients are looking at different windows. It seems that in this particular case, Screen has the better default where a window is only constrained in size if a smaller client is actively looking at it. This behaviour can be fixed by setting tmux's `aggressive-resize` option.

```
setw -g aggressive-resize on
```

### Multiple Clients Sharing One Session

Screen and tmux's behaviour for when multiple clients are attached to one session differs slightly. In Screen, each client can be connected to the session but view different windows within it, but in tmux, all clients connected to one session _must_ view the same window.

This problem can be solved in tmux by spawning two separate sessions and synchronizing the second one to the windows of the first. This is accomplished by first issuing a new session:

```
tmux new -s <base session>
```

Then pointing a second new session to the first:

```
tmux new-session -t <base session> -s <new session>
```

However, this usage of tmux results in the problem that detaching from these mirrored sessions will start to litter your system with defunct sessions which can only be cleaned up with some pretty extreme micromanagement. I wrote a script to solve this problem, call it `tmx` and use it simply with `tmx <base session name>`.

``` bash
#!/bin/bash

#
# Modified TMUX start script from:
#     http://forums.gentoo.org/viewtopic-t-836006-start-0.html
#
# Store it to `~/bin/tmx` and issue `chmod +x`.
#

# Works because bash automatically trims by assigning to variables and by 
# passing arguments
trim() { echo $1; }

if [[ -z "$1" ]]; then
    echo "Specify session name as the first argument"
    exit
fi

# Only because I often issue `ls` to this script by accident
if [[ "$1" == "ls" ]]; then
    tmux ls
    exit
fi

base_session="$1"
# This actually works without the trim() on all systems except OSX
tmux_nb=$(trim `tmux ls | grep "^$base_session" | wc -l`)
if [[ "$tmux_nb" == "0" ]]; then
    echo "Launching tmux base session $base_session ..."
    tmux new-session -s $base_session
else
    # Make sure we are not already in a tmux session
    if [[ -z "$TMUX" ]]; then
        # Kill defunct sessions first
        old_sessions=$(tmux ls 2>/dev/null | egrep "^[0-9]{14}.*[0-9]+\)$" | cut -f 1 -d:)
        for old_session_id in $old_sessions; do
            tmux kill-session -t $old_session_id
        done

        echo "Launching copy of base session $base_session ..."
        # Session is is date and time to prevent conflict
        session_id=`date +%Y%m%d%H%M%S`
        # Create a new session (without attaching it) and link to base session 
        # to share windows
        tmux new-session -d -t $base_session -s $session_id
        # Create a new window in that session
        #tmux new-window
        # Attach to the new session
        tmux attach-session -t $session_id
        # When we detach from it, kill the session
        tmux kill-session -t $session_id
    fi
fi 
```

<span class="addendum">Edit (2011/04/01) --</span> added new script logic so that defunct sessions are killed before starting a new one. Defunct sessions are left behind when tmux isn't quit explicitly.

<span class="addendum">Edit (2011/07/21) --</span> added configuration and the `tmx` script to my [tmux-extra repository](https://github.com/brandur/tmux-extra) on GitHub for more convenient access.

### Complete .tmux.conf

Here's my complete `.tmux.conf` for reference.

```
# C-b is not acceptable -- Vim uses it
set-option -g prefix C-a
bind-key C-a last-window

# Start numbering at 1
set -g base-index 1

# Allows for faster key repetition
set -s escape-time 0

# Set status bar
set -g status-bg black
set -g status-fg white
set -g status-left ""
set -g status-right "#[fg=green]#H"

# Rather than constraining window size to the maximum size of any client 
# connected to the *session*, constrain window size to the maximum size of any 
# client connected to *that window*. Much more reasonable.
setw -g aggressive-resize on

# Allows us to use C-a a <command> to send commands to a TMUX session inside 
# another TMUX session
bind-key a send-prefix

# Activity monitoring
#setw -g monitor-activity on
#set -g visual-activity on

# Example of using a shell command in the status line
#set -g status-right "#[fg=yellow]#(uptime | cut -d ',' -f 2-)"

# Highlight active window
set-window-option -g window-status-current-bg red
```
