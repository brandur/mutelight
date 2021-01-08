+++
location = "Berlin"
published_at = 2012-06-08T11:16:35-06:00
slug = "authbind"
title = "Authbind with a Simple Test"
+++

Having recently deployed an [ELB](http://aws.amazon.com/elasticloadbalancing/) in front of the production instances running the core Ruby application in our ecosystem, we started experimenting with the idea of removing Nginx from the box's HTTP stack. Why? Having been devised by people smarter than myself, I couldn't understand this idea initially, so let me explain a little further.

We use Unicorn because of its nice [restarting trick](https://github.com/blog/517-unicorn) that enables deploys with minimal complexity, and with no dropped connections. A side effect of the mechanism Unicorn uses to provide this feature is that the Unicorn master process runs on a single port which accepts connections, then delegates to one of the worker processes running on the box. That single port bound to by the master process makes a very nice target for the ELB, removing the need for a reverse proxy local to the box. One less component in the HTTP stack is one less piece that can fail, and reduces the incumbent knowledge required to properly manage our stack.

The fact that Unicorn was designed with the expectation of being run behind Nginx to buffer incoming requests and handle slower connections (it's right there on [the Unicorn philosophy page](http://unicorn.bogomips.org/PHILOSOPHY.html)) is another discussion, but we generally found that Unicorn runs pretty well on its own for our purposes. That is except when it's behind an ELB in HTTPS mode, but those findings deserve an article of their own.

Authbind
--------

Assuming that you want to deploy Unicorn on port 80, the very first challenge you'd run into is that on a typical Linux box, root privileges are required to bind to any ports below 1024. A great way to work around this is by using Authbind, start by installing it via your favorite package manager:

``` bash
aptitude install authbind
```

Authbind's permissions are managed with a special set of files in `/etc/authbind`. Create a file telling Authbind that binding to port 80 should be allowed:

``` bash
touch /etc/authbind/byport/80
```

Authbind determines that a user is allowed to bind an application to port 80 if they have access to execute this file. Change ownership of the file to the user your web server runs under (assumed to be `http` here) and make sure it has executable (`x`) permissions. Alternatively, we could accomplish the same thing using groups.

``` bash
# as root
chown http /etc/authbind/byport/80
chmod 500 /etc/authbind/byport/80
```

Test the setup using Python's built-in HTTP server:

``` bash
# as http user
authbind python -m SimpleHTTPServer 80 # Python 2.x
authbind python -m http.server 80      # Python 3.x
```

That's it! Notice that the web server command here should be prefixed by the `authbind` command for this to be allowed. Another Authbind invocation worth mentioning is `authbind --deep` which enables port binding permissions for the program being executed, as well as any other child programs spawned from it.
