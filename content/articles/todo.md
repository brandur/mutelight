+++
location = "San Francisco"
published_at = 2012-07-02T15:43:29-07:00
slug = "todo"
title = "A Simple @todo Pattern"
+++

Modern task tracking apps are both mindblowing in their sophistication, and somewhat ostentatious in their flashiness and sheer volume. We have apps like Remember the Milk that offer multiplatform support with cloud synchronization, and apps like Clear that provide such a beautiful interface and compelling experience that they beg to be used. The earliest apps on both the iPhone and iPad platforms were built for todo lists, and even Apple entered the game years later by introducing their Reminders app in iOS 5.

Despite this attractive selection, I use none of the above. Today, I wanted to share a very simple todo pattern that I've been using for months now with great results. Here it is in its entirety:

```
@todo
=====

  * Pick up the milk
  * h/Submit TPS report

Finished
--------

  * h/Order stationery

Defunct
-------

  * Submit talk

# vi: ts=2 sw=2 foldmethod=indent foldlevel=20
```

It's _that_ simple:

1. Current tasks go at the top.
2. Finished go under _finished_.
3. Tasks that were missed or have lapsed go under _defunct_.

The list stays open in Vim wrapped in Tmux pane at all times, and gets synced back to Dropbox. Finished items are transferred between lists using fast Vim bindings. If I think of something away from my computer, I add it to my phone, then transfer the task the next time I'm back.

The Vim hints at the end provide some nice folding behavior, which is useful when your finished list has become very long. Open and close individual lists using `zo` and `zc` respectively (the `foldlevel` hint at the end ensures that all lists are expanded when the file is first opened).
