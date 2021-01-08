+++
location = "Calgary"
published_at = 2011-07-14T10:41:00-06:00
slug = "the-objective-c-retain-property-pattern"
tiny_slug = "retain"
title = "The Objective-C Retain Property Pattern"
+++

Being pretty new to the whole Objective-C thing, I was initially a little confused about the proper way to implement properties that use the `retain` keyword. After doing some research on the subject, I've put together a complete sample based on what's considered the "correct" pattern when using them.

``` objc
@interface FactsViewController : NSObject
{
}

@property (nonatomic, retain) NSArray* facts;

@end

@implementation FactsViewController

@synthesize facts = __facts;

- (void) dealloc
{
    [__facts release];
    [super dealloc];
}

@end
```

In this straightforward example, the class interface defines the `facts` property, which is synthesized in the implementation with a backing field by the name of `__facts`.

The `retain` keyword indicates that the compiler should generate a property setter that retains a reference to a value set to it (i.e. increment the reference count by one). It also ensures that when the property is being set, if it had a previous value, that old value is properly released. The generated setter looks something like this:

``` objc
- (void) setFacts:(NSArray *)value
{
    [value retain];
    [__facts release];
    __facts = value;
}
```

The one other important aspect of the pattern is that `__facts` gets released in `dealloc`. Although a `retain` property handles a release when a new value is being set to it, it does not release automatically when an object is being deallocated&mdash;forcing us to release all properties explicitly.

Other resources online may suggest that rather than releasing the backing field with `[__facts release]`, it's cleaner just to clear it with `self.facts = nil` (which will implicitly release the old object). This is generally considered bad practice in case another object observing that property is triggered from inside the destructor, and tries to access the source object while it's in an inconsistent state.
