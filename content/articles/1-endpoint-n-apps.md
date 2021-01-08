+++
location = "San Francisco"
published_at = 2012-09-03T00:19:42-07:00
slug = "1-endpoint-n-apps"
title = "1 SSL Endpoint, N Apps"
+++

A few months ago, SSL was rolled out as an included feature for every app on the Heroku platform by enabling secure connection to each app's **heroku/herokuapp.com** domain, so even if developers prefer to use a custom domain, they'll at least have an SSL option for the components of their app where a secure connection is critical.

The platform's answer for developers requiring SSL on a custom domain is the use of the [SSL Endpoint](https://devcenter.heroku.com/articles/ssl-endpoint) addon, priced at $20 a month (the dark days of $100/mo. ssl:ip are finally over!). After adding SSL Endpoint to an app, a developer uploads their cert and an endpoint is created with a name like `mie-6498.herokussl.com`. He or she then CNAMEs their domain to the endpoint and secure requests are routed through with no app changes necessary.

And just a final bit of background: any given request on the Heroku platform enters through the [routing mesh](https://devcenter.heroku.com/articles/http-routing). The tl;dr is that it finds an appropriate runtime where the an app is deployed and forwards its requests through.

One Endpoint
------------

In case the $20/mo. per app for a custom domain seems a steep price to pay, you may be happy to find out that in many cases a single SSL Endpoint can be shared between many apps.

Requests coming through an SSL Endpoint follow the same rules as the rest of the platform&mdash;a request may enter through an endpoint but from there is routed through the mesh normally. Therefore, it's not an SSL Endpoint's associated app that decides where a request goes, but rather the incoming domain that's been CNAME'd to the endpoint.

A savvy developer can take advantage of this behavior to allow a single SSL Endpoint to route to any number of Heroku apps. For the connection to stay secure, the cert uploaded to the endpoint needs to be signed for any domains that you intended for use for it, but even a [free cert from StartCom](http://www.startcom.org/) allows two domains to be included without any special verification. A wildcard certificate (i.e. `*.mutelight.org`) will secure an entire stack of apps deployed into the Heroku cloud.

Below is a simple example demonstrating how a single endpoint is shared for both [brandur.org](https://brandur.org) and [facts.brandur.org](https://facts.brandur.org):

``` bash
#
# the app brandur-org below has ssl:endpoint
# the app facts-web does not
#

$ heroku addons -a brandur-org
ssl:endpoint

$ heroku addons -a facts-web
No addons installed

#
# both www.brandur.org (entry point for the app brandur-org) and
# facts.brandur.org (app facts-web) are CNAME'd to mie-6498
#

$ host www.brandur.org
www.brandur.org is an alias for mie-6498.herokussl.com.

$ host facts.brandur.org
facts.brandur.org is an alias for mie-6498.herokussl.com.

#
# both apps get a secure connection because brandur-org's cert includes both
# domains
#

$ heroku certs -a brandur-org
Endpoint                Common Name(s)           Expires               Trusted
----------------------  -----------------------  --------------------  -------
mie-6498.herokussl.com  facts.brandur.org,       2013-07-21 03:31 UTC  True
                        www.brandur.org
```
