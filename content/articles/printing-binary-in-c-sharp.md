+++
location = "Calgary"
published_at = 2009-06-04T00:00:00-06:00
slug = "printing-binary-in-c-sharp"
tiny_slug = "11"
title = "Printing Binary in C#"
+++

I've been looking on and off for the last couple days for a way to easily convert binary data to a hex representation in C# for printing. As it turns out, the `BitConverter` class is normally used to do this:

``` csharp
string s = System.BitConverter.ToString(
    new byte[] { 1, 15, 17 }
);
// s contains: "01-0F-11"
```

It's also handy for packing and unpacking binary data:

``` csharp
// Given an integer, the converter always returns a 4-byte 
// array. The length of this array will depend on the data 
// type it's given.
byte[] bytes = System.BitConverter.GetBytes(42);
// bytes contains: { 42, 0, 0, 0 }

// ToInt32 always reads 4 bytes
int num = System.BitConverter.ToInt32(bytes, 0);
// num is: 42
```
