+++
location = "Calgary"
published_at = 2009-02-21T00:00:00-07:00
slug = "dot-net-3-sequence-methods-moved-to-enumerable"
tiny_slug = "1"
title = ".NET 3.5 Sequence Methods Moved to Enumerable"
+++

I still see a lot of C# code samples around online that use `Sequence.Range` to generate a numerical range and `Sequence.Repeat` to repeat a number. This is a little confusing, because these methods may have been around in the early days of .NET 3.0, but have since moved to `System.Linq.Enumerable`.

Shown below is today's correct usage for `Range` and `Repeat`.

``` csharp
using System.Linq;

// Generate a range of integers
// This will enumerate 3, 4, 5, 6, 7, 8, 9
IEnumerable<int> nums = Enumerable.Range(3, 7);

// Repeat any value a given number of times
// This will enumerate 5, 5, 5
IEnumerable<int> fives = Enumerable.Repeat(5, 3);
// This will enumerate "string", "string"
IEnumerable<string> strings = Enumerable.Repeat("string", 2);
```
