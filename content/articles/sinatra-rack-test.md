+++
location = "Calgary"
published_at = 2012-12-29T11:32:30-08:00
slug = "sinatra-rack-test"
title = "Testing Sinatra With Rack-test"
+++

A common problem when starting out with Sinatra and trying to exercise what you've built with `rack-test` is that by default, Sinatra will swallow your errors and spit them out as a big HTML page in the response body. Trying to debug your tests by inspecting an HTML backtrace from `last_response.body` is a harrowing experience (take it from someone who's tried).

The solution is to tell Sinatra to raise errors back to you instead of burying them in HTML. Here's the proper combination of options to accomplish that:

``` ruby
set :raise_errors, true
set :show_exceptions, false
```

Here's a more complete example:

``` ruby
# app.rb
class App < Sinatra::Base
  configure do
    set :raise_errors, true
    set :show_exceptions, false
  end

  get "/" do
    raise "error!"
  end
end
```

``` ruby
# app_test.rb
describe App do
  include Rack::Test::Methods

  it "shows an error" do
    get "/"
  end
end
```
