+++
location = "Calgary"
published_at = 2011-05-31T10:34:00-06:00
slug = "what-i-learned-about-javascript-by-breaking-a-top-200-website"
tiny_slug = "ie-js"
title = "What I Learned About JavaScript by Breaking a Top 200 Website"
+++

Let no one contest the fact that I learn my lessons the hard way. Here's the story of how I broke the search engine on one of the Internet's largest websites for IE 6/7 due to a quirk in the JavaScript parser of these older Microsoft browsers.

The bug that I introduced was subtle enough that it made it past two levels of code reviews and through QA as well, which is a partial indicator of how much mind share IE 6/7 users get even within a company that officially supports them. A day after it went live, we'd already received five support calls reporting the problem; pretty good considering that most website users aren't particularly disposed to picking up a phone when they run into a problem.

The bug itself looks decidedly unimpressive considering that it probably impacted thousands of visitors:

``` js
filterTypes = {
    FILTER_A:        translations.filter_a, 
    FILTER_B:        translations.filter_b, 
    /* the comma on this last line will break IE 6/7 */
    FILTER_C:        translations.filter_c, 
};
```

Catch that? Any experienced JavaScript developer will probably see it easily, but most of us won't considering that the above code is perfectly valid under any modern day browser or JavaScript interpreter.

In case it wasn't obvious, the compilation problem is the trailing comma after `translations.filter_c` which dumb browsers like IE 6/7 will interpret as an indication that the array/object contains a trailing `undefined` element at the end, effectively like the following:

``` js
filterTypes = {
    FILTER_A:        translations.filter_a, 
    FILTER_B:        translations.filter_b, 
    FILTER_C:        translations.filter_c, 
    undefined
};
```

With the addition of the `undefined` element the code above is no longer valid JavaScript, and even a modern interpreter such as the one found in Node.JS will balk at it:

```
SyntaxError: Unexpected token }
```

Recourses
---------

After causing such an embarrassing problem, my obvious next reaction was to implement safeguards to ensure that it never happens again. Here are a few techniques to mitigate the risk of writing unsafe JS.

### CoffeeScript

There is a growing number of forward thinkers out there who believe that writing JS is akin to writing CIL for the CLI or bytecode for the JVM; in short that it's unnecessarily low-level. In the near future, languages like CoffeeScript will be leveraged to write compilable, efficient, and safe JavaScript. The right package for your platform will make transparently deploying CoffeeScript just as easy as if you were writing static JS files directly.

In fact, the future is now. As of [Rails 3.1](http://weblog.rubyonrails.org/2011/5/22/rails-3-1-release-candidate), CoffeeScript now takes the place of JavaScript by default. Personally, CoffeeScript has also been my preferred method of implementing JS for over a year now due to its added language features, and its concise and expressive syntax.

The equivalent CoffeeScript in this case looks very similar to the JS:

``` coffee
filterTypes = {
  FILTER_A:        translations.filter_a, 
  FILTER_B:        translations.filter_b, 
  FILTER_C:        translations.filter_c, 
}
```

The difference is that the Coffee compiler takes care of that trailing comma. The produced JS will look like this:

``` js
(function() {
  var filterTypes;
  filterTypes = {
    FILTER_A: translations.filter_a,
    FILTER_B: translations.filter_b,
    FILTER_C: translations.filter_c
  };
}).call(this);
```

Note also that another common pitfall (global name clashing) is avoided automatically by adding a function wrapper around the compiled code.

### JSLint

Don't have enough pull at your company to change its entire JavaScript backend to CoffeeScript? Well then perhaps JSLint is more to your liking. It's a JS code quality tool that can be integrated with editors like Vim to verify your syntax automatically when you open or save a `.js` file, and will make problems like the trailing array/object comma a thing of the past.

Vim users should try the [jslint.vim plugin](https://github.com/hallettj/jslint.vim). Keep in mind that as well as installing the plugin itself, you'll also need a JavaScript interpreter like Spidermonkey, Rhino, or Node installed on your system.
