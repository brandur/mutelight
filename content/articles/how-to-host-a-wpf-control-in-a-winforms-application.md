+++
location = "Calgary"
published_at = 2009-08-13T00:00:00-06:00
slug = "how-to-host-a-wpf-control-in-a-winforms-application"
tiny_slug = "17"
title = "How to Host a WPF Control in a WinForms Application"
+++

As people start making the jump from Windows Forms to WPF, many will hit an intermediate phase where it'll be useful to use WPF components in WinForms and vice-versa. Unfortunately, the two frameworks are not directly compatible because WPF uses `Control` under `System.Windows.Controls` while WinForms uses `Control` under `System.Windows.Forms` &mdash; and in terms of class hierarchies, the two `Control` classes are completely unrelated.

Luckily, Microsoft has had the good sense to give us a workaround. A WinForms application can contain a WPF control by putting the WPF control into an instance of `ElementHost`, a class deriving WinForms' `Control`, and then adding the host to a control collection normally:

``` csharp
/* our WPF control */
MapAnimationLayer animationLayer = new MapAnimationLayer();

/* container that will host our WPF control, we set it using 
 * the Child property */
ElementHost host = new ElementHost()
{
    BackColor = Color.Transparent, 
    Child     = animationLayer, 
    Dock      = DockStyle.Fill, 
};

/* now add the ElementHost to our controls collection 
 * normally */
Controls.Add(host);
```

The `ElementHost` class is part of a special namespace called `System.Windows.Forms.Integration`. Adding this namespace to your solution is slightly tricky because its DLL is not named as you'd expect; it's actually called `WindowsFormsIntegration`. From the Add Reference dialog window, scroll to the bottom of list in the .NET tab and you should see it.

WinForms in WPF
---------------

Similarly, another class called `WindowsFormsHost` allows you to insert WinForms controls into a WPF application. Here it is in XAML:

``` xml
<Window x:Class="WinFormsHostWindow"
    xmlns="http://schemas.microsoft.com/winfx/2006/xaml/presentation"
    xmlns:x="http://schemas.microsoft.com/winfx/2006/xaml"
    xmlns:wf="clr-namespace:System.Windows.Forms;assembly=System.Windows.Forms"
    >
    <Grid>
        <WindowsFormsHost>
            <!--
            the 'wf' namespace here refers to .NET's 
            System.Windows.Forms assembly, so this TextBox 
            is from WinForms
            -->
            <wf:TextBox x:Name="password" />
        </WindowsFormsHost>
    </Grid>
</Window>
```
