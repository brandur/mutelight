+++
location = "Calgary"
published_at = 2009-03-24T00:00:00-06:00
slug = "using-the-invoke-design-pattern-with-anonymous-methods"
tiny_slug = "8"
title = "Using the Invoke Design Pattern with Anonymous Methods"
+++

If you've done Windows forms development in .NET you've almost certainly seen the invoke design pattern before.

``` csharp
delegate void UpdateLabelTextCallback(string message);

void UpdateLabelText(string message)
{
    if (InvokeRequired)
        Invoke(
            new UpdateLabelTextCallback(UpdateLabelText), 
            message
        );
    else
        label1.Text = message;
}
```

This pattern is required on methods that update the UI in some way but which are called from a thread _other_ than the main/UI thread. Failing to use the invoke design pattern will usually result in an `InvalidOperationException` (cross-thread operation not valid) being thrown.

The code sample above shows the conventional implementation of the pattern using a delegate type which mirrors the method you're trying to callback and a call to `Invoke` with a reference back to the same method again. The `InvokeRequired` check decides whether an `Invoke` is necessary, and if it isn't, the method's normal logic is run.

Improving on This
-----------------

The conventional implementation is both long-winded and messy, and it can really start to get ugly when a method gets separated from its callback delegate type when the delegate needs to go in a different `#region`. A moderately complex form may need dozens of methods that follow the invoke pattern, so you'll end up with large sections of code dedicated to callback delegates for specific methods.

Luckily, using anonymous methods we can do a lot better. Let's rewrite our example to use an anonymous method instead of a delegate type.

``` csharp
void UpdateLabelText(string message)
{
    if (InvokeRequired)
        Invoke(new MethodInvoker(
            delegate { UpdateLabelText(message); }
        ));
    else
        label1.Text = message;
}
```

Much better! Now if we wanted to, we could take advantage of the fact that a call to `Invoke` on the main thread will have no detrimental effect on our program and make this method even shorter.

``` csharp
void UpdateLabelText(string message)
{
    Invoke(new MethodInvoker(
        delegate { label1.Text = message; }
    ));
}
```

It's possible to put any amount of logic into the anonymous method that is being invoked, even if it's more than one or two lines. This is legal of course, but starts to look pretty sloppy pretty fast.

Now with Lambdas
----------------

Anything we can do with an anonymous method we can also do with a lambda expression. For cases like this, it doesn't matter much which type of function you use, but if you're using lambdas everywhere else in your application, you can stick with them for invoke patterns too.

``` csharp
void UpdateLabelText(string message)
{
    Invoke(new MethodInvoker(
        () => label1.Text = message
    ));
}
```
