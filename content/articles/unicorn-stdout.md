+++
location = "San Francisco"
published_at = 2012-10-21T20:12:24-07:00
slug = "unicorn-stdout"
title = "Have Unicorn Log to $stdout"
+++

A strange quirk of Unicorn is that by default it will write all its logging output to `$stderr`. Even a relatively harmless operation like a restart will result in noise written to your system error log:

``` sh
executing ["/home/core/.bundle/gems/ruby/1.8/bin/unicorn_rails", "-c", "config/normalized_unicorn.rb"] (in /home/core)
forked child re-executing...
I, [2012-10-17T09:00:35.029145 #12322]  INFO -- : inherited addr=/tmp/core.sock fd=4
I, [2012-10-17T09:00:35.029885 #12322]  INFO -- : Refreshing Gem list
reaped #<Process::Status: pid=2784,exited(0)> worker=1
reaped #<Process::Status: pid=2785,exited(0)> worker=2
reaped #<Process::Status: pid=2783,exited(0)> worker=0
master complete
master process ready
worker=1 ready
worker=2 ready
worker=0 ready
```

Simply redefining Unicorn's logger to one pointing to `$stdout` will fix the problem:

``` ruby
# by default, Unicorn will log to $stderr; go to $stdout instead
logger Logger.new($stdout)
```
