+++
location = "Calgary"
published_at = 2009-06-24T00:00:00-06:00
slug = "a-look-at-less"
tiny_slug = "14"
title = "A Look at Less"
+++

A few weeks ago, I read about a new project called [Less](http://lesscss.org/), a CSS framework that extends CSS by adding a few simple, but extremely useful features. Less' killer concept is it uses existing CSS syntax so that CSS files are automatically valid Less files as well.

I tend to feel that CSS is complex enough as it is, and the last thing I need when building a website is an expansive CSS framework, or a separate build process for CSS files that needs its own detailed configuration. Less being what it is, I was immediately sold on it and as an experiment, converted this website's CSS to Less, a process that took about half an hour. The [finished product](http://github.com/fyrerise/mutelight/blob/5a36d9ff48dd39f74abbf80d4f0bb94d274a6301/css/screen.less) is on GitHub.

Installation
------------

The install process couldn't be easier. Assuming that you've installed Ruby and RubyGems from your package manager, just run the following command:

```
gem install less
```

After installing Less, a `.less` file is converted to `.css` using:

```
lessc source [destination]
```

Usage
-----

### Nested Rules

Less' best feature (IMO) is its support for nested rules. Where in normal CSS you'd have to mark CSS rules with long selectors, you can actually nest selectors inside of other selectors in Less. I tend to write long selectors in CSS because I find that it helps keep display logic out of HTML files, so nesting rules with Less was a natural extension for me.

``` css
/* before */
div.content
{
    ...
}

div.content ol
{
    padding: 0 0 0 20px;
}

div.content ol li
{
    text-align: justify;
    margin: 0 0 8px 0;
}
```

``` css
/* after */
div.content
{
    ...

    ol
    {
        padding: 0 0 0 20px;

        li
        {
            text-align: justify;
            margin: 0 0 8px 0;
        }
    }
}
```

Notice how well the nesting complements lists! A nice side effect of Less is that when you nest _ALL_ your document's CSS rules (like [I've done](http://github.com/fyrerise/mutelight/blob/5a36d9ff48dd39f74abbf80d4f0bb94d274a6301/css/screen.less) with mine), the layout of your Less file starts to resemble the layout of your markup file, making rules easy to find.

### Variables

Less introduces variables to CSS, a welcome addition to CSS. Using variables is very simple: define them somewhere in your Less file with `@<var>: <value>;` then refer to them elsewhere with `@<var>`.
 
``` css
/* primary text color, for main text body */
@primary_color: #a9a9a8;

/* secondary text color, for links, alternate text, etc. */
@secondary_color: #7d7b88;

...

a
{
    color: @secondary_color;
}
```

One limitation to variables is that they can't be used with strings, a problem which I describe in more detail below.

### Mixins and Operations

Less also supports mixins and operations, two features I didn't end up putting to use, but which are described very succinctly on [the Less website](http://lesscss.org/).

Limitations
-----------

While switching over to Less, I did encounter a few problems that appear to be unsolvable as far as the framework stands today.

### Nested Selector Groups

When writing CSS, it's common practice to group selectors together when they share properties. For example, if we wanted both our normal links and visited links to be blue, we'd write something like the following.

``` css
a:link, a:visited
{
    color: #00f;
}
```

Naturally, you might want to extend this to Less' nested rules as well.

``` css
body
{
    /* nested groups */
    a:link, a:visited
    {
        color: #00f;
    }
}
```

Unfortunately, if you tried to do so, Less would not output your CSS how you intended. Instead, it would write something like this:

``` css
/* this doesn't work as we would have expected */
body a:link, a:visited
{
    color: #00f;
}
```

For these groups to work as we intended, both groups would have to the full CSS selector like `body a:link, body a:visited`.

### String Variables

With the addition of variables, a logical step would be start using them to store string fragments that could be joined with other fragments to build something useful, like a URL:

``` css
@my_images: "http://mutelight.org/images/";

body
{
    background: #000 url(@my_images + "back.png");
}
```

Less doesn't seem to like this right now, but hopefully it'll make it into a future version.

Syntax Highlighting in Vim
--------------------------

Of course, Less just wouldn't be as good without some sort of syntax highlighting in Vim! Fortunately, highlighting Less files as CSS works well. Just add the following line to your `~/.vimrc`:

```
au BufRead,BufNewFile *.less setfiletype css
```
