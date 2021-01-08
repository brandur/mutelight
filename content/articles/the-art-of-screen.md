+++
location = "Calgary"
published_at = 2009-11-18T00:00:00-07:00
slug = "the-art-of-screen"
tiny_slug = "22"
title = "The Art of Screen"
+++

During all the years I've used Linux, unquestionably the one greatest productivity improvement I ever experienced was the day I discovered [GNU Screen](http://www.gnu.org/software/screen/). Screen allows you to many things, the most important of which is to continue your terminal sessions from anywhere, and between different consoles.

Activity Monitoring
-------------------

The power of screen continues to amaze me: today I stumbled across a feature by accident that I'd been missing out on all this time (I probably should've given the [man page](http://www.manpagez.com/man/1/screen/) a good read a long time ago). It's called activity monitoring, and it lets you have screen notify you when something happens in a background window. If you're a heavy IM or IRC user like me, these notifications can be tremendously useful.

Turn activity monitoring on for your current window using `C-a M`. Using this same shortcut again will toggle activity monitoring back to off. Now when something happens in this window, Screen will flash a message across its status bar, and display an `@` next to its title.

<img src="https://d25zpof2afwnhk.cloudfront.net/the-art-of-screen/screen-status.png" alt="Notice the @ next to a modified window's title" />

The message displayed when something occurs can be changed using `activity <message>` in your configuration file. Change it to an empty string (`""`) to kill the message completely.

### Finch + Irssi

The problem with Screen's activity monitoring is that it will be trigged by any change whatsoever. This leads to problems for common command line apps like Finch and Irssi that update constantly.

In Finch, buddies are constantly coming online and going offline, which triggers the monitor. A partial solution for this is to just close your buddy list window (`M-c`), so that only changes in your chat windows go through to you.

The biggest problem in Irssi (besides the fact that IRC channel are being updated constantly) is that by default, its status bar has a clock that will flip once a minute, and trigger the monitor. The easiest thing to do here is to get rid of the clock: `/statusbar window remove time`. Irssi will save your configuration changes automatically.
