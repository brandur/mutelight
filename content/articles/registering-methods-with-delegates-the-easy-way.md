+++
location = "Calgary"
published_at = 2009-06-25T00:00:00-06:00
slug = "registering-methods-with-delegates-the-easy-way"
tiny_slug = "15"
title = "Registering Methods with Delegates the Easy Way"
+++

Every C# developer is familiar with the syntax for registering a method to a delegate or event:

``` csharp
_myButton.Click += new EventHandler(myButton_Click);
```

I'm finding that a lot of people don't seem to know about the alternate "shorthand" syntax for doing the same thing that was introduced in .NET 2.0:

``` csharp
_myButton.Click += myButton_Click;
```

I mention this because I feel that the shorthand syntax is preferable in a number of ways, starting with the fact that it's shorter, therefore eliminating some of the normal boilerplate cruft and leaving you with more readable code. It's also easier to understand for developers new to C#. Consider the following:

``` csharp
_myButton.Click += new EventHandler(myButton_Click);
_myButton.Click -= new EventHandler(myButton_Click);
```

Every C# developer eventually learns that even though you created two separate delegate objects in this example, the `myButton_Click` method still gets unregistered. The reason is that when the delegate is looking for a delegate to unregister, it doesn't compare delegates by reference, but instead by the underlying method(s) they will call. I'll be the first to admit that I found this behaviour a little counterintuitive when I started coding C#.

Now, compare the previous example to the same code, but using shorthand notation:

``` csharp
_myButton.Click += myButton_Click;
_myButton.Click -= myButton_Click;
```

Regardless of how much C# you now, it's pretty obvious here that `myButton_Click` is going to be unregistered by the time this snippet finishes executing.

<span class="addendum">Addendum &mdash;</span> While looking around for delegate information, I came across a really [great article on weak event listeners in C#](http://www.codeproject.com/KB/cs/WeakEvents.aspx). This is an important concept because despite all the memory safety C# provides to you, it's still easy to get memory leaks from objects dangling by a single delegate reference that wasn't property unregistered.
