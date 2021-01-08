+++
location = "Silver Star"
published_at = 2010-02-11T00:00:00-07:00
slug = "generating-a-permalink-slug-in-haskell"
tiny_slug = "23"
title = "Generating a Permalink Slug in Haskell"
+++

One of the basic features of any blogging software is a function for generating "slugs" for your articles. A slug is normally a URL-friendly version of your article's title that's had any special characters and spaces stripped out, for example "My Awesome Article!!" might become "my-awesome-article".

Generating a slug is a pretty straightforward exercise with regular expressions, but might not be quite so obvious in Haskell because of the rather arcane way that the regex libraries are provided.

Regex
-----

Regular expressions in Haskell are *weird*. A `Text.Regex.Base` module provides an interface to a [large variety of possible backends](http://www.haskell.org/haskellwiki/Regular_expressions) that do the heavy-lifting. The commonly-used functions in base are `=~` and `=~~`, both polymorphic, meaning they behave differently depending on the type signature we specify for them. Read [this beautifully written Haskell regex tutorial](http://www.serpentine.com/blog/2007/02/27/a-haskell-regular-expression-tutorial/) to understand this in more depth.

I've chosen to use the PCRE backend for Regex; a very fast module that is particularly suited to working with bytestrings. I'd install it on Archlinux like so (the cabal command is suited for any system):

```
sudo pacman -Sy pcre
cabal install regex-base regex-pcre
```

Another caveat in Haskell's regex libraries is that no regex replace function is provided. Here's one borrowed from spookylukey's Haskell blog:

``` haskell
import Text.Regex.PCRE ( (=~~) )
import qualified Data.ByteString.Lazy.Char8 as B

{- | Replace using a regular expression. ByteString version. -}
regexReplace ::
    B.ByteString          -- ^ regular expression
    -> B.ByteString       -- ^ replacement text
    -> B.ByteString       -- ^ text to operate on
    -> B.ByteString
regexReplace regex replacement text = go text []
    where go str res =
              if B.null str
              then B.concat . reverse $ res
              else case (str =~~ regex) :: Maybe (B.ByteString, B.ByteString, B.ByteString) of
                     Nothing -> B.concat . reverse $ (str:res)
                     Just (bef, _ , aft) -> go aft (replacement:bef:res)
```

Bytestrings
-----------

It's generally recommended for fast code to use the `ByteString` type rather than the Haskell's built-in string, because the built-in string is considered slow. A Google search will reveal that massive memory usage and speed improvements in real-world applications are attributed to switching from strings to bytestrings.

Normally to use a bytestring in our code we'd have to write something like:

``` haskell
B.pack "my string"
```

However, if you have a reasonably recent version of GHC, you can specify the `-XOverloadedStrings` language option flag to define bytestrings the same way you define strings.

Generating Slugs
----------------

Now that we've got the basics in place, we can write a simple function to generate slugs for a blog:

``` haskell
import GHC.Unicode ( toLower )
import qualified Data.ByteString.Lazy.Char8 as B

makeSlug :: B.ByteString -> B.ByteString
makeSlug = regexReplace "[ _]" "-" 
         . regexReplace "[^a-z0-9_ ]+" "" 
         . B.map toLower
```

Save this to a file called `Slug.hs` and load up a GHCI prompt to test it out (remember to use `-XOverloadedStrings`!):

```
$ ghci -XOverloadedStrings
Prelude> :m Data.ByteString.Lazy.Char8
Prelude Data.ByteString.Lazy.Char8> :load Slug.hs
*Main Data.ByteString.Lazy.Char8> :set prompt "Prelude> "
Prelude> unpack $ makeSlug "My Awesome Article!!"
"my-awesome-article"
Prelude> unpack $ makeSlug "My Awesome w/ Numbers 789_2"
"my-awesome-w-numbers-789-2"
```

That's the gist of it. The slug function itself can now be tweaked as desired.
