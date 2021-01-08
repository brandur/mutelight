+++
location = "Calgary"
published_at = 2009-05-29T00:00:00-06:00
slug = "finally-tuples-in-c-sharp"
tiny_slug = "9"
title = "FINALLY: Tuples in C#"
+++

It's taken until .NET 4.0, but Microsoft has finally decided to add tuples to the C# language. Here are their prototypes:

``` csharp
/* all tuples are defined under the System namespace */

[SerializableAttribute]
public class Tuple<T1> : IStructuralEquatable, 
    IStructuralComparable, 
    IComparable

[SerializableAttribute]
public class Tuple<T1, T2> : IStructuralEquatable,
    IStructuralComparable,
    IComparable

[SerializableAttribute]
public class Tuple<T1, T2, T3> : IStructuralEquatable,
    IStructuralComparable,
    IComparable
    
/* and all the way up to ... */

[SerializableAttribute]
public class Tuple<T1, T2, T3, T4, T5, T6, T7, TRest> : IStructuralEquatable,
    IStructuralComparable,
    IComparable
```

It's a bit unfortunate for those of us who'd already added our own tuple implementations to our solutions because we'd given up hoping for built-in types years ago. Better late than never though.

Tuples can be created using either the constructors for your type's corresponding arity, or more easily using the static `Tuple.Create` method:

``` csharp
Tuple<int, string> who1 = new Tuple<int, string>(1, "Bob");
Tuple<int, string> who2 = Tuple.Create(2, "Bella");
```

Tuple elements are accessed using the readonly `Item1`, `Item2`, `Item3`, ... properties:

``` csharp
string name1 = who1.Item2; /* "Bob" */
string id2   = who2.Item1; /* 2 */
```

In case you're wondering, you can build a really big tuple using the last type arity: `Tuple<T1, T2, ..., TRest>`. The final `TRest` parameter is used to hold a second tuple inside the first one, and you could even put a third tuple inside the second and so on. As you might expect though, the syntax is clumsy:

``` csharp
Tuple<int, int, int, int, int, int, int, Tuple<int, int, int>> bigTuple = 
    Tuple.Create(
        /* T1 through T7 */
        1, 2, 3, 4, 5, 6, 7, 
        /* TRest */
        Tuple.Create(8, 9, 10)
    );
```

So, we're not going to see anymore `out` parameters right guys?
