+++
location = "San Francisco"
published_at = 2013-06-05T07:18:13-07:00
slug = "params"
title = "Discriminating Input"
+++

While designing our [V3 platform API](https://devcenter.heroku.com/articles/platform-api-reference), we made the decision to make the formatting in our requests and responses as symmetric as possible. Although common for an API to return JSON, it's not quite as common to take it as input, but is our recommended usage for all incoming `PATCH`/`POST`/`PUT` requests.

Largely reasons largely for developer convenience, we decided to allow fall back to form-encoded parameters as well (for the time being at least), so we put together a helper method that allows us to handle these in a generic fashion. It looks something like this:

``` ruby
class API < Sinatra::Base
  post "/resources" do
    params = parse_params
    Resource.create(name: params[:name])
    201
  end

  private

  def parse_params
    if request.content_type == "application/json"
      indifferent_params(MultiJson.decode(request.body.read))
      request.body.rewind
    else
      params
    end
  end
end
```

By specifying `Content-Type: application/json`, JSON-encoded data can be sent to and read by the API:

```
curl -X POST https://api.example.com/resources \
    -d '{"name":"my-resource"}' -H "Content-Type: application/json"
```

The more traditional method for encoding POSTs is to use the `application/x-www-form-urlencoded` MIME type which looks like `company=heroku&num_founders=3` and is sent in directly as part of the request body. Rack will decode form-encoded bodies by default and add them to the `params` hash, so our API easily falls back to this:

```
curl -X POST https://api.example.com/resources -d "name=my-resource"
```

(Note that Curl will send `Content-Type: application/x-www-form-urlencoded` by default.)

Good so far, but a side-effect that we hadn't intended is that our API will also read standard query parameters:

```
curl -X POST https://api.example.com/resources?name=my-resource
```

On closer examination of the Rack source code, it's easy to see that Rack is trying to simplify its users lives by blending all incoming parameters into one giant input hash:

``` ruby
def params
  @params ||= self.GET.merge(self.POST)
rescue EOFError
  self.GET.dup
end
```

While not a problem per se, this does widen the available options for use of API to cases beyond what we considered to be reasonable. We cringed to think about seeing technically correct, but somewhat indiscriminate usage examples:

```
curl -X POST https://api.heroku.com/apps?region=eu -d "name=my-app"
```

By re-implementing the helper above to ignore `params`, the catch-all set of parameters, and instead use `request.POST`, which contains only form-encoded input, we an exclude query input:


``` ruby
  def parse_params
    if request.content_type == "application/json"
      indifferent_params(MultiJson.decode(request.body.read))
      request.body.rewind
    elsif request.form_data?
      indifferent_params(request.POST)
    else
      {}
    end
  end
```

## rack-test

As an addendum, it's worth mentioning that `rack-test` also sends `application/x-www-form-urlencoded` by default (and always will unless you explicitly override `Content-Type` to a non-nil value), and that's what's going on when you do this:

``` ruby
it "creates a resource" do
  post "/resources", name: "my-resource"
end
```

We found that it was worthwhile writing our tests to check the primary input path foremost, so most look closer to the following:

``` ruby
it "creates a resource" do
  header "Content-Type", "application/json"
  post "/resources", MultiJson.encode({ name: "my-resource" })
end
```
