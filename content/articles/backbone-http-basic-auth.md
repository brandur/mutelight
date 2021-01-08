+++
location = "San Francisco"
published_at = 2012-09-02T20:18:29-07:00
slug = "backbone-http-basic-auth"
title = "HTTP Basic Authentication with Backbone"
+++

I recently ran across a situation where I needed to call down to an HTTP basic authenticated API in a Backbone app. Short of doing some hacking of the source, Backbone doesn't provide an easy way to accomplish this, so here's one simple option for your perusal.

Backbone relies on the inclusion of jQuery or Zepto in your project to provide the underlying infrastructure for making AJAX calls. If you're using jQuery, there's a function called `$.ajaxSetup` that will set options before every AJAX call. Use it to set the `Authorization` header (_warning: CoffeeScript_):

``` coffee
$.ajaxSetup
  headers:
    Authorization: "Basic #{toBase64(":secret-api-password")}"
```

Under HTTP basic, both the user and password need to be base64 encoded before being sent along to the server. JavaScript doesn't provide utilities to handle that out of the box, so the `toBase64` function above needs to be implemented to get this example running.

A nice option is [CryptoJS](http://code.google.com/p/crypto-js/). Download the package and include the following files in your project:

* `core.js`
* `enc-base64.js`

Now you're ready to implement `toBase64` and complete this example:

``` coffee
toBase64 = (str) ->
  words = CryptoJS.enc.Latin1.parse(str)
  CryptoJS.enc.Base64.stringify(words)
```
