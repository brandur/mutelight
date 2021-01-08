+++
location = "Calgary"
published_at = 2011-06-17T14:22:00-06:00
slug = "simple-side-by-side-live-and-sandbox-rails-deployment-with-nginx-and-phusion-passenger"
tiny_slug = "simple-phusion"
title = "Simple Side-by-side Live and Sandbox Rails Deployment with Nginx and Phusion Passenger"
+++

Quite some time ago a little language called PHP built itself a little following and, quite suddenly, became one of the largest driving forces of the Internet. One of the key features that enabled this kind of forward momentum was the simplicity with which it could be deployed: simply enable an Apache extension, drop some PHP files into a web directory, and you're live. That was a fine to have back in the early days of the new millenia, but a common misconception that still survives to this day is that PHP is still easier than anything else to deploy.

These days, I'd prefer the simplicity of deploying a Rails app over a PHP app any day, and [Phusion Passenger](http://www.modrails.com/) is the technology that's made this possible. I'd like to present some simple Nginx and Passenger configuration to deploy side-by-side development and production Rails environments.


``` nginx
server {
    listen 80;
    server_name mutelight.org;

    root /home/http/www/mutelight.org/private/askja/public/;
    passenger_enabled on;

    # Keep a minimum around for fast responsiveness in production
    passenger_min_instances 1;
}

server {
    listen 80;
    server_name pre.mutelight.org;

    root /home/http/www/mutelight.org/private/pre/askja/public/;
    passenger_enabled on;

    rails_env development;
}
```

Here I have two [Askja](http://github.com/brandur/askja) repositories cloned at `mutelight.org/private/askja` and `mutelight.org/private/pre/askja`; Passenger makes these available at `mutelight.org` and `pre.mutelight.org` respectively. The key configuration line in both blocks is `passenger_enabled on;` which tells Passenger that the given roots are Rails apps. Also note that configured paths point to a Rails app's `public` directory, rather than Rails root.

Production at `mutelight.org` mandates that at least one Passenger instance should be kept alive at all times with `passenger_min_instances 1;`, thus keeping the site responsive when a client hits a dynamic page like [Sitemap](http://mutelight.org/sitemap.xml). This is more important for a site that performs a lot of full page caching, because cached pages will be served by Nginx directly and Passenger instances will shut themselves down with nothing to do.

Development at `pre.mutelight.org` specifies `rails_env development;` which is useful for serving better error messages, and avoiding page caching while working in the sandbox. Note that the sandbox has no `passenger_min_instances` directive, allowing all Passenger instances for the sandbox to be shut down while it's not in use.

After Nginx is up and running with the two Rails apps deployed, use the `passenger-status` command (as root) to see how many Passenger instances are running at any given time.
