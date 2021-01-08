+++
location = "Calgary"
published_at = 2009-08-11T00:00:00-06:00
slug = "application-dot-crawl"
tiny_slug = "16"
title = "Application.Crawl()"
+++

Most .NET developers are familiar with the standard entry point of a WinForms application:

``` csharp
Application.Run(new MainForm());
```

I've been finding recently that the name of this method is rapidly becoming outdated as modern Microsoft technology finds its way into more applications. Developers like myself can get the wrong idea about what this method actually does while skimming source code. Therefore, I respectfully propose that Microsoft obsoletes the `Run` method, and replaces it with something a little more true to form:

``` csharp
Application.Crawl(new MainForm());
```

This insight comes from my recent experience installing Expression Blend 3. I'd previously noticed that Microsoft applications seemed to be getting slower after installing the beta for Visual Studio 2010<sup class="footnote" id="fnr1"><a href="#fn1">1</a></sup>, but Microsoft has really taken the Blend installer to the next level.

Close to 100% of a single core was eaten up for the entire half-hour install process of Blend. The slowdown was so extreme that my mouse cursor skipped around the screen in 20 pixel increments as I made futile attempts to continue my other work. For the record, I'm running a Core 2 E8400 @ 3.00 GHz with 4 GB of memory.

<div class="figure">
    <a href="https://d25zpof2afwnhk.cloudfront.net/application-dot-crawl/cpu-eater-3.png" title="Link to full-size image"><img src="https://d25zpof2afwnhk.cloudfront.net/application-dot-crawl/cpu-eater-3-small.png" alt="Sus microprocessorius, more widely known as the common CPU hog, in its natural habitat" /></a>
    <p><strong>Fig. 1:</strong> <em>Sus microprocessorius,</em> more widely known as the common CPU hog, in its natural habitat</p>
</div>

The good news is that Expression Blend itself seems to be very usable. Next week I promise to write about something more constructive.

<p class="footnote" id="fn1"><a href="#fnr1"><sup>1</sup></a> VS 2010 is a beta release. Some performance improvements are bound to make it into the final product.</p>
