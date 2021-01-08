+++
location = "Calgary"
published_at = 2009-10-25T00:00:00-06:00
slug = "swallowing-an-extensible-exception-in-haskell"
tiny_slug = "19"
title = "Swallowing an Extensible Exception in Haskell"
+++

As of version 6.10.1 of the GHC, the `Control.Exception` namespace now uses extensible exceptions, the old style exceptions having been relocated to `Control.OldException`.

This is great news, but in some cases the new extensible exceptions are a little more fussy about how they're used than the old exceptions were, and can spit out some pretty cryptic error messages if handled improperly. For example, I found myself wanting to swallow any exceptions that occurred while trying to create a new CouchDB database. I started out with some straightforward code like this:

``` haskell
runCouchDB' (createDB "logs") `catch` \_ -> return ()
```

Here's the resulting error message:

```
Ambiguous type variable `e' in the constraint:
  `GHC.Exception.Exception e'
    arising from a use of `Control.Exception.Base.catch'
                 at <interactive>:1:0-71
Probable fix: add a type signature that fixes these type variable(s)
```

Being pretty new to Haskell, it wasn't immediately obvious to me what the problem was here, but after doing some research I realized that GHC wanted to be provided with the type of exception that I wanted to catch; the idea being that the compiler forces developers to think about what they're doing to some extent. In my case, this happens to be an `ErrorCall` exception instance, so I rewrote my code like this:

``` haskell
runCouchDB' (createDB "logs") `catch` \(_ :: ErrorCall) -> return ()
```

&hellip; and it compiles! If I had been doing something in my lambda function above, this type might have been inferred for me by the compiler, but I had to define it explicitly because I didn't use the result in any way. A list of exception instances is available in the [documentation for `Control.Exception`](http://www.haskell.org/ghc/docs/6.10-latest/html/libraries/base/Control-Exception.html#t%3AException).
