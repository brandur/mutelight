+++
location = "Calgary"
published_at = 2010-02-15T00:00:00-07:00
slug = "validate-xml-according-to-schema-location-hints"
tiny_slug = "24"
title = "Validate XML According to Schema Location Hints"
+++

For those of us unfortunate enough to still be writing XML documents, it's fairly common practice to "hint" at the location of an XSD schema for a document using attributes in its root node:

``` xml
<?xml version="1.0" encoding="utf-8" ?>

<root xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
      xsi:noNamespaceSchemaLocation="file:Schema.xsd">
    ...
</root>
```

While working with .NET's `XPathDocument` in C#, I noticed that when using this class with just a URI string argument, it will not validate against your schema. The following code throws no exception even if your XSD indicates an invalid file:

``` cs
using System.Xml.XPath;

XPathDocument d = new XPathDocument("path/to/file.xml");
```

Properly validating your XML involves giving the `XPathDocument` an `XmlReader` instance that's been setup to respect schema location hints:

``` cs
using System.Xml;
using System.Xml.Schema;
using System.Xml.XPath;

XmlReader r = XmlReader.Create(
    "path/to/file.xml", 
    new XmlReaderSettings() {
        ValidationType  = ValidationType.Schema, 
        ValidationFlags = 
            XmlSchemaValidationFlags.ProcessSchemaLocation, 
    }
);
XPathDocument d = new XPathDocument(r);
```
