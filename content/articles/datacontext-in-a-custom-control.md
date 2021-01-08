+++
location = "Calgary"
published_at = 2010-11-26T00:00:00-07:00
slug = "datacontext-in-a-custom-control"
tiny_slug = "32"
title = "DataContext in a Custom Control"
+++

In WPF or Silverlight, new developers will often stumble across the concept of _data context_ when learning about data binding. Data context is a mechanism that allows a framework element to define a data source for itself and its children to use when they participate in data binding. For example, a group box containing a number of associated data fields might specify its data context as some model containing a property for each of those fields. Each child control of that group box only needs to specify a relative binding path containing the name of the property it would like to bind to. Linking to a model object containing that property isn't needed because that model was already set as the data context on a parent element and is inherited.

When building a custom control, a very common technique is set that control's data context back to itself so that child controls can be bound to properties in that same control's code-behind.

``` xml
<UserControl x:Class="MyControl" 
    DataContext="{Binding RelativeSource={RelativeSource self}}">

    <!-- Only a relative path is needed because data -->
    <!-- context was set at a higher level           -->
    <TextBlock Text="{Binding Title}" />

</UserControl>
```

Did you see the bug in the code above? It's not immediately apparent and takes a little experience to realize that a user control should never specify its own data context in its definition. But why?

The answer becomes more apparent when we try to combine the data context concept with our new control elsewhere.

``` xml
<Grid DataContext="{StaticResource ViewModel}">

    <!-- Here we'd expect this control to be bound to -->
    <!-- MyContent on our ViewModel resource          -->
    <my:MyControl Content="{Binding MyContent}" />

</Grid>
```

In the above example, we'd expect `MyControl` to behave like any other framework element and bind its `Content` to `MyContent` on our `ViewModel` resource. Unfortunately, for anyone trying to use this control, it's actually binding its `Content` property to `MyContent` on itself (which probably doesn't even exist). The reason is that we already hard-coded a data context into the control's definition, which will take precedent in this case.

Fortunately, there's an easy way to solve this problem. Instead of specifying our data context on the root of the control itself, we should specify the data context on the control's top-level child element. This is often a `Grid` called `LayoutRoot` that Visual Studio generates automatically.

``` xml
<UserControl x:Class="MyControl">

    <Grid x:Name="LayoutRoot"
        DataContext="{Binding RelativeSource={RelativeSource self}}">

        <!-- This child element will still bind the -->
        <!-- same way!                              -->
        <TextBlock Text="{Binding Title}" />

    </Grid>

</UserControl>
```

Framework elements within the control itself can bind using a relative path in the same way, and data context in classes using the control will not be polluted, which will help prevent unexpected side effects.
