+++
location = "San Francisco"
published_at = 2012-06-19T19:40:09-07:00
slug = "netrc"
title = "Building an API with Netrc"
+++

If you track the progress of the [Heroku client](https://github.com/heroku/heroku), you may have noticed that in the last few months we've [switched the way that your credentials are stored](https://github.com/heroku/heroku/blob/master/CHANGELOG#L278) over from a custom format in `~/.heroku` to an older and more normalized storage standard, `.netrc`. This isn't an isolated event either, you may have noticed that GitHub has recently changed the recommended [clone method on new repositories to https](https://github.com/brandur/dummy), which has the side-effect of bypassing your standard access with `~/.ssh/id_rsa`. How do you get back to not being prompted for your credentials every time you push to the repository? Netrc.

`.netrc` is an old standard that dates all the way back to the days of FTP, that romantic wild west era of the Internet where the concept of "passive mode" kind of made sense. Its job is to store a user's credentials for accessing remote machines in a simple and consistent format:

    machine brandur.org
      login brandur@mutelight.org
      password my-very-secure-personal-password
    machine mutelight.org
      login brandur@mutelight.org
      password my-even-secure-password-with-a-number-on-the-end-7

Although originally intended for FTP, its use has since expanded to a other network clients including Git, Curl, and of course Heroku.

A common pattern that I've run into while building API's over the last few months is to protect APIs with HTTP basic authentication. This isn't necessarily the best solution in the long term, passing tokens provisioned with OAuth2 may be better, but it's a mechanism that can be set up quickly and easily.

Take this Sinatra app as an example:

``` ruby
# run with:
#   gem install sinatra
#   ruby -rubygems api.rb

require "sinatra"

set :port, 5000

helpers do
  def auth
    @auth ||= Rack::Auth::Basic::Request.new(request.env)
  end

  def auth_credentials
    auth.provided? && auth.basic? ? auth.credentials : nil
  end

  def authorized?
    auth_credentials == [ "", "my-secret-api-key" ]
  end

  def authorized!
    halt 401, "Forbidden" unless authorized?
  end
end

put "/private" do
  authorized!
  200
end
```

After running it, we can test our new API with Curl:

```
curl -i -u ":my-secret-api-key" -X PUT http://localhost:5000/private

HTTP/1.1 200 OK
X-Frame-Options: sameorigin
X-XSS-Protection: 1; mode=block
Content-Type: text/html;charset=utf-8
Content-Length: 0
Connection: keep-alive
Server: thin 1.3.1 codename Triple Espresso
```

Now here's the interesting part. Add the following lines to your `.netrc`:

```
machine localhost
  password my-secret-api-key
```

Try the same Curl command again but using the `-n` (for `--netrc`) flag:

```
curl -i -n -X PUT http://localhost:5000/private
```

Voilà! The speed of being able to run ad-hoc queries against an API you're building rather than drudging up your API key every time turns out to be a huge win practically, and it's a pattern that I now use regularly during development.

A limitation that's hinted at above is that you can only have a single entry for `localhost`. Generally, I find that this isn't a huge problem because most of the APIs I want to hit are deployed in a staging or production environment with a named URL.

Heroku
------

Now onto a nice real-world example. Are you a Heroku user? Have you updated your Gem since February 2012? If the answer to both these questions is **yes!**, try this from a console:

```
curl -n https://api.heroku.com/apps
```

Security
--------

A glaring problem with `.netrc` is that it keeps a large number of your extremely confidential credentials out in the open in plain text. Presumably, the file is `chmod`'ed to `600` and you're using full-disk encryption, but that's still probably not enough (say someone happens to find your computer unlocked).

The [netrc](https://rubygems.org/gems/netrc) gem used by the Heroku client will try to find a GnuPG encrypted file at `~/.netrc.gpg` before falling back to the plain text version. Although this convention is far from a standard, it's still recommended security practice.
