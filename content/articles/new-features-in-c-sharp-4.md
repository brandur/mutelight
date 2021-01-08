+++
location = "Calgary"
published_at = 2009-03-11T00:00:00-06:00
slug = "new-features-in-c-sharp-4"
tiny_slug = "5"
title = "New Features in C# 4.0"
+++

Admittedly a few months late, here's my recap of the new C# language features being introduced in 4.0, summarized from the [PDC 2008 talk _the Future of C#_ presented by chief architect Anders Hejlsberg](http://channel9.msdn.com/pdc2008/TL16/). Note that this post doesn't cover any of the new parallel features of 4.0.

Dynamically Typed Objects
-------------------------

C# 4.0's new dynamic features are built on top of .NET's Dynamic Language Runtime (DLR), an extension of the CLR. IronPython and IronRuby were built on top of the DLR, and now C# and VB will be using its features as well.

Conventionally, if we received an object from somewhere without a class or interface that we could reference from our code, we'd have to store it as an arbitrary `object`. To make matters worse, if we wanted to call a method on that object, we'd have to use reflection, and we'd end up with some pretty nasty-looking code.

``` csharp
// Prior to 4.0, if we received an object with no static 
// type that we could reference, we'd call a method on it 
// like this
object calc = someObject.GetCalculator();
object result = calc.GetType().InvokeMember(
    "Add", 
    BindingFlags.InvokeMethod, 
    null, 
    new object[] { 1, 1 }
);
int sum = Convert.ToInt32(result);
```

In 4.0, we can improve this situation by using the `dynamic` keyword. We can store our object as dynamic, and from then on we can make dynamic calls on that object using the well-known dot operator.

``` csharp
// After 4.0, we can statically type the object we received 
// to be dynamic (in the words of Anders Hejlsberg). We can 
// then call methods dynamically on it using just our normal 
// dot operator.
dynamic calc = someObject.GetCalculator();
// Note that when we assign the result of our dynamic call 
// to a variable, that result is converted dynamically
int sum = calc.Add(1, 1);
```

Unfortunately, the `dynamic` keyword will come with some classic drawbacks of dynamically-typed languages; the most noticeable of which to VS users may be that statement completion on dynamic objects will not be possible. Other problems will be that fewer errors can be caught at compile time, and that runtime performance will not be as good as it is for static objects.

Anders points out that the objective of these new dynamic features is not to make C# a dynamically-typed language, but to make it less painful to talk to dynamic portions of your code when doing so is necessary (for example, talking to JavaScript from Silverlight).

``` python
def GetCalculator():
    return Calculator()

# This is the calculator class referenced above written in 
# Python
class Calculator():
    def Add(self, x, y):
        return x + y
```

The Python code shown above is a possible implementation for the arbitrary calculator object referenced in the previous code samples. The code shown below illustrates how this Python code could be used from a C# program.

``` csharp
dynamic pythonCalc = 
    Python.CreateRuntime().UseFile("Calculator.py");

// Our simple example from above
dynamic calc = pythonCalc.GetCalculator();
int sum = calc.Add(1, 1);

// Since Python is purely dynamic, this will work on any 
// type that implements the '+' operator
TimeSpan twoDays = calc.Add(
    TimeSpan.FromDays(1), 
    TimeSpan.FromDays(1)
);
```

Optional and Named Parameters
-----------------------------

Reinforcing their long-standing reputation of being late to the game, Microsoft has finally decided to implement optional and named parameters in C#. Below shows two `TimeSpan` constructors and a possible new revision of them using optional parameters.

``` csharp
// Methods that previously had to have multiple overloads 
// like this
public TimeSpan(int hours, int minutes, int seconds)
: self(0, hours, minutes, seconds)
{ ... }

public TimeSpan(int days, int hours, int minutes, 
                int seconds)
{ ... }

// Could be changed to something more concise like this
public TimeSpan(int seconds = 0, int minutes = 0, 
                int hours = 0, int days = 0)
// (Note this constructor does not exist, and likely 
// never will. I've re-arranged the parameters to show 
// how it probably would have been designed in the start 
// had this language feature existed then.)
{ ... }
```

Methods written this way can be called in a wide range of different ways as long as all their required parameters are present. Named parameters are written like _parameter name_: _value_.

``` csharp
// 10 seconds
var t1 = new TimeSpan(10);

// 24 hours, 20 minutes, 10 seconds
var t2 = new TimeSpan(10, 20, 24);

// 7 days, 24 hours, 20 minutes, 10 seconds
var t3 = new TimeSpan(10, 20, 24, 7);

// 7 days, 10 seconds
var t4 = new TimeSpan(10, days: 7);

// 7 days, 24 hours, 20 minutes, 10 seconds
var t5 = new TimeSpan(days: 7, hours: 25, minutes: 20, 
                      seconds: 10);
```

Improved COM Interoperability
-----------------------------

Up until 4.0, many COM interop calls were messy ordeals involving many unused parameters, all of which had to be using `ref`.

``` csharp
object fileName = "My Document.docx";
object missing  = System.Reflection.Missing.Value;

doc.SaveAs(ref fileName, 
           ref missing, ref missing, ref missing, 
           ref missing, ref missing, ref missing, 
           ref missing, ref missing, ref missing, 
           ref missing, ref missing, ref missing, 
           ref missing, ref missing, ref missing);
```

C# 4.0 will finally resolve this issue with the introduction of optional parameters and making `ref` optional.

``` csharp
doc.SaveAs("My Document.docx");
```

This improvement has been made possible by:

* **Optional and named parameters**: unused parameters can safely be omitted
* **Optional `ref` modifier**: `ref` need not be used (but only for interop calls, you still need `ref` elsewhere)
* **Interop type embedding**: interop methods that previously took arguments of type `object` now take `dynamic` instead

Co- and Contra-variance
-----------------------

.NET arrays have been co-variant since their introduction, however, they are not _safely_ co-variant. Notice how in the example below we can assign any object back to the string array, and cause a runtime exception.

Remember that objects which are co-variant can be treated as less derived.

``` csharp
string[] strings = GetStringArray();
ProcessUnsafely(strings);

void ProcessUnsafely(object[] objects)
{
    objects[0] = "Okay";
    objects[1] = new TimeSpan(); // Not okay
}
```

C# generics have always been invariant, but in C# 4.0 they will be safely co-variant.

``` csharp
// This code would previously not have compiled
List<string> strings = GetStringList();
ProcessSafely(strings);

void ProcessSafely(IEnumerable<object> objects)
{
    // Safety comes from the fact that IEnumerable cannot 
    // be changed
}
```

Co-variance is implemented on an interface using the `out` keyword on any co-variant type parameters.

``` csharp
public interface IEnumerable<out T>
{
    IEnumerable<T> GetEnumerator();
}

public interface IEnumerator<out T>
{
    T Current { get; }
    bool MoveNext();
}
```

Contra-variance can also be implemented using the `in` keyword. Remember that objects which are contra-variant can be treated as more derived.

``` csharp
public interface IComparer<in T>
{
    int Compare(T x, T y);
}

// Legal
IComparer<object> objComparer = GetComparer();
IComparer<string> strComparer = objComparer;
```

A few notes on co- and contra-variance to remember:

* `in` and `out` are only supported for interfaces and delegate types
* Value types are always invariant (`IEnumerable<int>` cannot be stored to `IEnumerable<object>`)
* `ref` and `out` parameters are always invariant
