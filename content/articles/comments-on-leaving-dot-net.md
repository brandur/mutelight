+++
location = "Calgary"
published_at = 2010-12-16T00:00:00-07:00
slug = "comments-on-leaving-dot-net"
tiny_slug = "36"
title = "Comments on \"Leaving .NET\""
+++

Yesterday I read [Dave Newman's article on leaving .NET](http://whatupdave.com/post/1170718843/leaving-net) which has recently been making its Internet rounds. I can't help but comment on how accurately he's managed to describe the .NET community and its long list of shortcomings. The article lists a number of problems:

* "Not invented here" syndrome (specifically MVC which Microsoft introduced as some ground-breaking new technology even though it came out years after the technology was already well-established in open source)
* Microsoft doesn't accept patches for code that it claims is open source
* There is generally no collaboration in the community on open source projects
* Codeplex isn't "real" open source, and the community needs to abandon it

Ironically, the very reason that Dave abandoned the Microsoft stack may be the same reason that big industry is so keen on it. To a big company, another big company is a safety net. Companies are more reliable than people: they don't leave the country, or leave projects half finished, and most importantly, they provide a nice stationary legal target. These are all good reasons that help explain why in general, big enterprise prefers leveraging frameworks built by Microsoft over those built by loosely organized groups of liberal hooligans.

One problem that he hit dead on is Microsoft's deep afflication by "not invented here" syndrome. Microsoft has never, and will never work with the community when it comes to open source projects. A few examples that I think about often:

* Releasing the Entity Framework when NHibernate was already widely used. Entity may have caught up, but in its early days, NHibernate was unquestionably superior.
* Microsoft Unit Testing Framework's came out years after a multitude of open source were used all over the place, including big industry.
* Reproducing the popular Subversion in the form of Team Foundation (see below).
* Evangelizing ASP.NET MVC when Rails and .NET frameworks based on Rails were already mature.
* The .NET framework itself! Let's face it, it's hard to refute that C# and .NET were _strongly_ "influenced" by Java and the JVM.

Microsoft has created quite a successful business model around an innovation strategy based on copying what they see out in the wild, and good for them, they're still one of the best possible examples of how to make big money in the software industry. It's that same successful model however, that also ensures that Microsoft is constantly lagging behind new features, sometimes by years, and many of those new features that are actually innovated in house fail to catch on.

Another problem with the copy model is that once in a while, revolutionary innovations in some technology comes at the wrong time (i.e. in the middle of development of a clone), and get missed. An example of this is Team Foundation, which was developed as a successor to SourceSafe and a response to Subversion creeping into .NET territory. By the time Team Foundation was released, the leading edge had switched gears completely and begun using distributed version control (i.e. Git, Mercurial, Darcs), a completely new paradigm, and the rest of industry had started to follow. As far as I know, due to poor timing, Microsoft still has no answer for modern source control.

Despite all this, we should be thankful that there's still a thriving open source community doing work in .NET. Efforts on the Mono project might be the best example of this, but NHibernate has also rebranded themselves recently (NHForge), and .NET projects on Github may even be becoming more relevant than those on Codeplex. Combined with its overwhelming use in the business world, .NET is still a good place to be.
