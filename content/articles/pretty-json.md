+++
location = "San Francisco"
published_at = 2012-09-02T23:02:34-07:00
slug = "pretty-json"
title = "Prettifying JSON for Curl Development"
+++

In the same vein as my post from a few weeks ago on [developing an API with Curl and Netrc](/netrc), here's a handy trick for sending pre-prettified JSON output back to clients, but only those that are identifying as Curl. The reasoning being that prettified JSON isn't useful most of the time, but it's a nice touch while using developing or testing against an API with Curl.

Bundle MultiJson in your `Gemfile`:

``` ruby
gem "multi_json"
```

Now define a helper for identifying Curl clients, and use it wherever encoding JSON:

``` ruby
# sample Sinatra app

helpers do
  def curl?
    !!(request.user_agent =~ /curl/)
  end
end

get "/articles" do
  articles = Article.all
  [200, MultiJson.encode(articles, pretty: curl?)]
end
```
