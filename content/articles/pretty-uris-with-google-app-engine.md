+++
location = "Calgary"
published_at = 2009-03-09T00:00:00-06:00
slug = "pretty-uris-with-google-app-engine"
tiny_slug = "4"
title = "Pretty URIs with Google App Engine"
+++

Getting pretty URI/URLs in Google app engine is easy, but the first time I looked I wasn't able to figure out how to do it from Google's documentation, so I've saved some instructions here.

The first step is to add in some regex capture groups when sending your URIs to `WSGIApplication`. The code below specifies that `ScoreHandler` can receive either one argument like `/game/<game_id>/score` or two arguments like `/game/<game_id>/score/<score_id>`.

``` python
from google.appengine.ext   import webapp
from handlers.score_handler import ScoreHandler
from wsgiref.handlers       import CGIHandler

def main():
    application = webapp.WSGIApplication([
        (r'/game/(.+)/score/(.+)', ScoreHandler), 
        (r'/game/(.+)/score',      ScoreHandler), 
    ], debug=True)
    CGIHandler().run(application)

if __name__ == '__main__':
    main()
```

When one of these URI expressions is hit, `WSGIApplication` will expect to find a function in the corresponding handler with a number of arguments equal to the number of regex groups that were matched.

``` python
from google.appengine.ext import webapp

class ScoreHandler(webapp.RequestHandler):

    # Receives one or two arguments, depending on which URI 
    # was accessed
    def get(self, game_id, score_id=None):
        self.response.out.write(
            "game_id = %s, score_id = %s" % (game_id, score_id)
        )
```

As shown above, we can also map multiple URIs with different numbers of regex capture groups to the same handler class by providing a function with default arguments. In this example, if `/game/2/score` was hit, `game_id` would have a value of 2 while `score_id` would be `None`. If `/game/2/score/3` was hit, `game_id` would be 2 and `score_id` would be 3.
