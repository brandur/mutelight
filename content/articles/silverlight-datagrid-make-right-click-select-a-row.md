+++
location = "Calgary"
published_at = 2010-09-01T00:00:00-06:00
slug = "silverlight-datagrid-make-right-click-select-a-row"
tiny_slug = "28"
title = "Silverlight DataGrid: Make Right-click Select a Row"
+++

A slightly irking characteristic of many data grid libraries is that there seems to be a disconnect between expected usability in a grid, and the default behaviour of those same grids. I've found this to be a problem in all .NET grids that I've used including Infragistics and Microsoft's WPF/Silverlight implementations.

For example, one thing I expect to be able to do is right-click on a grid's row, and have it pull up a context menu with a list of actions to perform on it. Now that context menus have been added to Silverlight 4, half the battle is won, but you still can't get the "correct" right-click behavior out of the box.

Here's how to coax the grid to behave a bit more properly with `VisualTreeHelper`.

```  xml
<!-- XAML -->
<UserControl x:Class="MyControl" 
    xmlns:grid="clr-namespace:System.Windows.Controls;assembly=System.Windows.Controls.Data"
    xmlns:toolkit="http://schemas.microsoft.com/winfx/2006/xaml/presentation/toolkit">
    
    <grid:DataGrid MouseRightButtonDown="_grid_MouseRightButtonDown">
        <!-- Define a context menu for our grid here -->
        <toolkit:ContextMenuService.ContextMenu>
            <toolkit:ContextMenu>
                <toolkit:MenuItem Header="Action!" />
            </toolkit:ContextMenu>
        </toolkit:ContextMenuService.ContextMenu>
    </grid:DataGrid>
</UserControl>
```

``` cs
// Code-behind
private void _grid_MouseRightButtonDown(object sender, MouseButtonEventArgs e)
{
    IEnumerable<UIElement> elementsUnderMouse = 
        VisualTreeHelper
            .FindElementsInHostCoordinates(e.GetPosition(null), this);
    DataGridRow row = 
        elementsUnderMouse
            .Where(uie => uie is DataGridRow)
            .Cast<DataGridRow>()
            .FirstOrDefault();
    if (row != null)
        _grid.SelectedItem = row.DataContext;
}
```

The code above uses `VisualTreeHelper` to find any `UIElement`s under the mouse when the user performs a right-click. A set of `UIElement`s is returned including grid rows and cells as well as a good number of other elements. We filter this set for an instance of `DataGridRow`, get is bound data item, then tell the grid to select that object.
