+++
location = "Berlin"
published_at = 2012-06-09T14:33:13-06:00
slug = "ruby-debug"
title = "A Cross-version Debug Pattern for Ruby"
+++

As circumstances would have it, I work with a few projects that need to be compatible with both Ruby 1.8 and 1.9. I also like to use a debugger, and getting one to work seemlessly under both versions isn't exactly intuitive, so I wanted to share a nice pattern that I've been using recently in my personal projects.

The traditional debugger `ruby-debug` has been known to be 1.9 incompatible for some time now, but more recently, its updated version `ruby-debug19` is no longer 1.9 compatible having been broken by 1.9.3 without a new release. Luckily, the awesome new [`debugger`](https://github.com/cldwalker/debugger) gem stepped in to fill the gap.

Include both debuggers in your `Gemfile` with platform conditionals:

``` ruby
group :development, :test do
  gem "debugger",   "~> 1.1.3",  :platforms => [:ruby_19]
  gem "ruby-debug", "~> 0.10.4", :platforms => [:ruby_18]
end
```

I debug pretty often, but don't like to type a lot, so I usually include a shortcut in my `test_helper.rb` to get a debugger invoked quickly regardless of the Ruby version that you're running:

``` ruby
def d
  begin
    require "debugger"
  rescue LoadError
    require "ruby-debug"
  end
  debugger
end
```

Now drop it into a file like so:

``` ruby
def requires_frequent_debugging
  risky_call rescue nil
  Singleton.manipulate_global_state
  d # the debugger will start on the next line
  Model.do_business_logic
  super
end
```

It might seem like the debugger would start in the `d` method rather than where you want to debug, forcing you to finish the stack frame before you could start debugging. Fortunately, that's not the case. The `d` method has returned by the time the debugger is invoked, leaving you exactly where you want to be.

In a classic case of open-source overkill, I've extracted the pattern described above into a trivial gem called [d2](https://github.com/brandur/d2). Throw it in your Gemfile, make sure that your project is either using `Bundler.setup` or including `require 'd2'` somewhere, then use `d2` somewhere to trigger the debugger.

<span class="addendum">Aside --</span> A slightly interesting Ruby tidbit related to the code above is that we use `rescue LoadError` because a generic `rescue` only catches `StandardError` exceptions. `LoadError` is derived from a different hierarchy headed by `ScriptError`.
