+++
location = "Calgary"
published_at = 2009-11-10T00:00:00-07:00
slug = "new-net-4-feature-enum-hasflag"
tiny_slug = "21"
title = "New .NET 4.0 Feature: Enum.HasFlag()"
+++

A minor inclusion to the new .NET 4.0 framework is the addition of the `Enum.HasFlag()` method. It behaves much like you'd expect, allowing you to check whether an enum instance contains a given flag.

``` cs
[Flags]
public enum PacketOptions
{
    None       = 0x00, 
    All        = 0xff, 
    Compressed = 0x01, 
    Encrypted  = 0x02, 
}

PacketOptions opts = PacketOptions.All;
Assert.That(opts.HasFlag(PacketOptions.Compressed), Is.True);
```

Previously, flag testing was done using `(opts & PacketOptions.Compressed) != 0`, which was a little awkward because the `!=` operator has higher precedence than `&`, hence the extraneous parenthesis.

If you'd attempted to implement your own `HasFlag` extension method, you'd have realized that it was not possible due to the way C# handles its type constraints (yep, I've been there).

Discovered via [Reed Copsey](http://reedcopsey.com/?p=77).
