+++
location = "Berlin"
published_at = 2012-06-10T15:21:51-06:00
slug = "bin-console"
title = "Your Ruby App Should Have a `bin/console`"
+++

Those of us who have worked or are working on Rails are somewhat spoiled by the ability to boot `script/rails console` and immediately start running commands from inside our projects. What may not be very well known is that this console isn't a piece of Rails black magic, and makes a nice pattern that extends well to any other type of non-Rails Ruby project.

Here's the basic pattern:

``` ruby
#!/usr/bin/env ruby

require "irb"
require "irb/completion" # easy tab completion

# require your libraries + basic initialization

IRB.start
```

With the right initialization, this will immediately drop you into a console with all your project's models, classes, and utilities available, and even with tab completion! It also translates easily over to cloud platforms, being only one `heroku run bin/console` away, so to speak.

I picked up the idea somewhere at Heroku where public opinion generally sways against heavy Rails-esque frameworks and towards more custom solutions built from the right set of lightweight components.

Here's a real world example for the [bin/console of Hekla](https://github.com/brandur/hekla/blob/master/bin/console), which runs this technical journal:

``` ruby
#!/usr/bin/env ruby

require "irb"
require "irb/completion"
require "bundler/setup"
Bundler.require

$: << "./lib"
require "hekla"

DB = Sequel.connect(Hekla::Config.database_url)

require_relative "../models/article"

# Sinatra actually has a hook on `at_exit` that activates whenever it's
# included. This setting will supress it.
set :run, false

IRB.start
```
