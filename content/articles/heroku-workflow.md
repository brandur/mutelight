+++
location = "San Francisco"
published_at = 2012-07-09T17:58:52-07:00
slug = "heroku-workflow"
title = "The Heroku CLI's API Workflow"
+++

Ever wondered how the Heroku command line client accomplishes its work? We don't talk about it much, but the Heroku CLI isn't a black box, it's a fairly thin consumer of our own RESTful API, and that means everything you do in your day to day workflow on Heroku is available to work with in a programmatic fashion. The CLI even uses our own implementation of the API library called [heroku.rb](https://github.com/heroku/heroku.rb) (and by the way, heroku.rb is a great choice if you want to consume the API from Ruby).

A handy tool that we use here regularly is inspecting the CLI's workflow by telling Excon to send its output to standard out. Try it for yourself:

``` bash
EXCON_STANDARD_INSTRUMENTOR=true heroku list
```

Any calls that are implemented via heroku.rb make their requests using Excon, but a few of the older endpoints still use Restclient. If you run into one of these, you can do something very similar:

``` bash
RESTCLIENT_LOG=stdout heroku drains -a mutelight
```
