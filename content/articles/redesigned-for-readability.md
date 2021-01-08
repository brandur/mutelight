+++
location = "Calgary"
published_at = 2011-04-28T20:19:00-06:00
slug = "redesigned-for-readability"
tiny_slug = "redesign"
title = "Redesigned for Readability"
+++

A few months ago I realized that I disliked the design of my blog enough that I'd usually run my own articles through [Readability](http://readability.com) before trying to read them. Despite being a design insensitive coder, even I realized that I had a case of really horrendous usability on my hands.

Here are a few features of the new design:

* **Google Web Fonts:** I use the [Google Web Fonts API](http://www.google.com/webfonts) to provide a consistent font experience across most platforms. The best part about Google Web Fonts though is that they look great on Linux browsers that would otherwise annoyingly render Linux Font (TM).
* **Text width:** articles should be approximately 100 characters wide for near optimum readability, but will shrink to fit smaller windows where necessary.
* **Text-shadow:** a CSS3 trick that makes text more readable, especially against a distracting background like the one you see here.
* **White space:** I thought that I was doing a pretty good job of using white space before, but was probably mistaken. I like it so much that from now on, I might just use it more often.

Although not quite as visible, I've moved away from **Nanoc** to a new backend:

* <del>**Ruby on Rails:** Mutelight now runs a custom built Rails-powered backend. Overall I find that the flexibility of having the full power of Rails available far outweighs the cost of writing a backend from scratch (working with Rails, writing a backend for a site like this really doesn't take very long). Moving away from a static compiler may come with performance concerns, but in the end Rails' full page caching comes out just as fast. I'll talk about this more in a future article.</del> (it still doesn't run on a static generator, but it's no longer powered by Rails)
* **Client-side syntax highlighter:** due to concerns with the speed of an EC2 micro instance, I've moved away from Pygments to [jQuery Syntax Highlighter](http://balupton.com/projects/jquery-syntaxhighlighter/). It does just as good of a job, and produces pristine markup as a bonus. After investigating available options extensively, I now feel that client-side syntax highlighters are the way forward these days.
* **Tiny URLs:** short links now look like `/a/redesign` instead of `/a/22`, providing more useful context.
* <del>**Formats for free:** all pages will now respond to JSON requests as well as the default HTML. Try accessing [http://mutelight.org/articles/redesigned-for-readability.json](http://mutelight.org/articles/redesigned-for-readability.json) for example.</del> (a subsequent redesign has disabled this feature; I no longer consider it to have much merit)

Another step that I've taken is to remove Disqus commenting (for now at least). I was hoping that these might be effective for correcting mistakes that I've made in articles, but in practice this doesn't seem to happen very often. If you discover a mistake, e-mail me at **brandur@mutelight.org** or send a pull request for a correction [on Github](https://github.com/brandur/mutelight), and I'll get around to correcting myself as soon as possible.

That said, I'm hugely thankful to the people who took the time to leave a comment on the old blog. I do hope you'll continue to provide feedback through other channels.

Lastly, I'll probably be making tweaks over the next few weeks. Enjoy!
