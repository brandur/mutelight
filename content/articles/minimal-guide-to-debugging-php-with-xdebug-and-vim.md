+++
location = "Calgary"
published_at = 2011-04-29T11:13:00-06:00
slug = "minimal-guide-to-debugging-php-with-xdebug-and-vim"
tiny_slug = "xdebug"
title = "Minimal Guide to Debugging PHP with XDebug and Vim"
+++

For those of us moving to PHP from modern IDEs like Visual Studio, `var_dump` and `die` are simply not enough! Luckily, I can report a generally positive experience using Vim and XDebug as an effective debugger. Here's a short guide to installation and usage.

Installation
------------

### Server-side

Install XDebug via your system's package manager, PECL/PEAR, or by compiling it from source. See [XDebug's installation instructions](http://www.xdebug.org/docs/install) for up-to-date information on this. Keep in mind though, if you already have something like Cachegrind installed on your server, you may already have and use XDebug.

Configure your PHP installation so that it knows where to find XDebug's `.so` and that debugging should be enabled. This might be done in either `php.ini` or `/etc/php5/conf.d/xdebug.ini`.

```
; point this to wherever your xdebug.so is located
zend_extension=/usr/lib/php5/20090626/xdebug.so
xdebug.remote_enable=1
xdebug.remote_host=localhost
xdebug.remote_port=9000
```

**Warning:** for XDebug to work properly, you may have to disable your Zend debugger if you have one installed (may be found in `/etc/php5/conf.d/zend.ini`). With Zend enabled, I could connect to an XDebug session, but the next debug command I tried to run (step into, step over, run) would terminate the session.

### Client-side

Install the [DBGp remote debugger interface for Vim](http://www.vim.org/scripts/script.php?script_id=2508). This plugin will allow Vim to interact with Xdebug.

This script requires that your Vim be compiled with Python and signs support. Type `vim --version` and look for `+python` and `+signs` in the features list. If you don't have these features, try installing your distribution's **gVim** package, often it comes along with a command executable that includes them.

Usage
-----

To initialize a debugging session, XDebug will attempt to make a connection to the remote host and port that you specified above (`localhost:9000`), so you need to make sure that your Vim is reachable at that address. Unless you're running Vim on the same server as PHP, this may involve building an SSH tunnel back to your development box.

Start debugging by opening Vim and hitting `F5`, Vim will say:

```
waiting for a new connection on port 9000 for 10 seconds...
```

You have ten seconds to open a web browser and navigate to your running PHP app. It's important to know however, that the debug session will only start if you include the parameter `XDEBUG_SESSION_START=1` with the request (see the _optimization_ section below to remove this requirement). An example URL:

```
http://mysite.com/?XDEBUG_SESSION_START=1
```

Your Vim should now enter debug mode with the debug shortcuts should be listed in a buffer on the right. For example, use `F2` to step into, `F3` to step over, and `F12` to inspect the property under your cursor.

Typing `:Bp` in command mode will toggle a breakpoint on the current line, and `F5` will run the program to hit it. `:Up` and `:Dn` move up and down the stack trace.

`,e` will allow you to evaluate arbitrary code. To fully inspect a property, I recommend running `print_r($var, true)` from eval mode. A word of warning here though: evaluating invalid code may kill your debug session, but your client (Vim) won't realize that it's defunct.

Command line PHP can also be debugged as long as XDebug is enabled in the PHP CLI `php.ini`. This is especially useful for debugging unit tests.

Optimization
------------

The requirement of including `XDEBUG_SESSION_START=1` as a request parameter can be waived by adding the following code to your `php.ini`/`xdebug.ini` and restarting your server. With this setting enabled, XDebug will attempt to start a debugging session with every request.

```
xdebug.remote_autostart=1
```
