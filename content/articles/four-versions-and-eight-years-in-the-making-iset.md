+++
location = "Calgary"
published_at = 2009-06-23T00:00:00-06:00
slug = "four-versions-and-eight-years-in-the-making-iset"
tiny_slug = "13"
title = "Four Versions and Eight Years in the Making: ISet"
+++

Another addition to .NET 4.0 that I'd previously missed is an interface that most of us had long ago given up on ever seeing in as part of a Microsoft product: the `ISet<T>` interface. The [set documentation is on MSDN](http://msdn.microsoft.com/en-us/library/dd412081.aspx), and if you've previously used the `HashSet` class, it should look pretty familiar.

I say that `ISet` is four versions and eight years in the making, because practically speaking, a framework like .NET should've probably included set interfaces and implementations from the get-go. This article's original title was _four versions and eight years *late*_, which in all fairness, might be a more accurate description of the situation. Java, a language that C# borrowed from heavily (to say the least), has had its own `Set` interface since version 1.2, released in December 1998, years before C#'s first release in 2001.

Even giving Microsoft the benefit of the doubt by assuming that they had their reasons for not putting in a proper set until version 3.0, it's still pretty hard to see why they didn't see fit to put in `ISet` at the same time as `HashSet`. Another 4.0 addition, the `SortedSet`, is likely what finally forced Microsoft's hand. Having two separate set classes with no common interface probably felt a little wrong, even to them.

On a positive note, .NET 4.0 is really shaping up to be an outstanding release. With the addition of `ISet`, [tuples](/articles/finally-tuples-in-c-sharp.html), [optional/named parameters and co-variance](/articles/new-features-in-c-sharp-4.html), it seems like Microsoft is really focussing on fixing some of their previous omissions.

<span class="addendum">Addendeum --</span> Did anyone else encapsulate `Dictionary` to roll their own set implementation? I ended up using mine all the way up until `HashSet` walked onto the scene in 3.0.</p>
