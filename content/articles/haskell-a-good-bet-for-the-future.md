+++
location = "Calgary"
published_at = 2009-09-11T00:00:00-06:00
slug = "haskell-a-good-bet-for-the-future"
tiny_slug = "18"
title = "Haskell: A Good Bet for the Future"
+++

Everyone is pushing concurrency these days. Recently functional programming has re-emerged into the spotlight because the paradigm is particularly conducive to parallel programming. Three languages that have been getting a special amount of attention are Haskell, Erlang, and Clojure (and all three with good reason!).

Microsoft will be releasing a major concurrency framework called [Parallel Extensions](http://en.wikipedia.org/wiki/Parallel_Extensions) along with .NET 4.0, which will include parallel LINQ (PLINQ) and the parallel tasks (AKA futures) library. As usual, they're well behind the curve, but at least when Microsoft gets into something, a lot of corporate programmers will start taking the technology seriously.

My problem recently has been picking one technology to focus on learning. Naturally, a large part of my time will go towards learning the new .NET 4.0 concurrency, but as a language for personal use, Haskell has won my attention. Instead of offering an amateur explanation of its advantages over other languages, I'm just going to link to some of the particularly interesting articles I've read on it recently:

* [No monadic headaches: multi-core concurrency is easy in Haskell](http://cgi.cse.unsw.edu.au/~dons/blog/2007/11/26#no-headaches) &mdash; article on how Haskell's concurrency can act as a direct replacement for Erlang's (but with much nicer syntax), using message passing channels
* [Multicore Programming in Haskell Now!](http://donsbot.wordpress.com/2009/09/05/defun-2009-multicore-programming-in-haskell-now/) &mdash; summary of Haskell's wide range of parallel programming models
* [Supercompilation for Haskell](http://community.haskell.org/~ndm/downloads/slides-supercompilation_for_haskell-03_mar_2009.pdf) (PDF warning) &mdash; how to compile Haskell to rival C's speed
* [Haskell Arrays, Accelerated Using GPUs](http://www.scribd.com/doc/19637022/Haskell-Arrays-Accelerated-with-GPUs) &mdash; executing massively parallel operations on the GPU from Haskell
