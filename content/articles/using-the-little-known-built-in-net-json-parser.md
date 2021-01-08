+++
location = "Calgary"
published_at = 2010-12-29T00:00:00-07:00
slug = "using-the-little-known-built-in-net-json-parser"
tiny_slug = "40"
title = "Using the Little-known Built-in .NET JSON Parser"
+++

You wouldn't know it, but an alternative to the multitude of JSON parsers built for the .NET platform is to use the one built into the core libraries of .NET itself. It feels like it flies under the radar a bit, which I suspect is due to it not being [listed on the JSON site](http://www.json.org/). The main advantage of preferring the built-in library is that there's one less dependency to track, but a bonus is that we can easily deserialize from multiple formats because the parser works by treating JSON like XML, as I'll discuss in a moment.

Use the JSON deserializer by building an `XmlReader` instance from the static class `JsonReaderWriterFactory`, then using your favorite XML querying technique to extract data (I use XPath below). JSON classes come from the `System.Runtime.Serialization.Json` namespace, accessible by adding a project reference to `System.Runtime.Serialization`.

``` cs
string filename = args[0];
byte[] buffer = File.ReadAllBytes(filename);
XmlReader reader = 
    JsonReaderWriterFactory
        .CreateJsonReader(buffer, new XmlDictionaryReaderQuotas());

XElement root = XElement.Load(reader);

// The fields we'd like to extract
XElement form   = root.XPathSelectElement("//form");
XElement status = root.XPathSelectElement("//status");
XElement type   = root.XPathSelectElement("//type");

// Field set
IEnumerable<XElement> fields = root.XPathSelectElements("//fields/*");
```

If any of the elements we selected didn't exist, the corresponding `XElement` will hold a `null`. This program will read a simple JSON document like the following:

``` js
{
    "form":   13, 
    "status": "NEW", 
    "type":   "FORM", 
    "fields": {
        "field": { "key": 50, "value": 1 } 
    }
}
```

Now for the interesting part! Notice in the code snippet above that we load our JSON document with an `XmlReader` and treat it almost exactly like XML. We can leverage this property by rewriting our simple parser to read either a JSON or an XML document.

``` cs
string filename = args[0];
byte[] buffer = File.ReadAllBytes(filename);
XmlReader reader = new FileInfo(filename).Extension.ToLower() == "xml"
    ? XmlReader.Create(new MemoryStream(buffer)) 
    : JsonReaderWriterFactory
        .CreateJsonReader(buffer, new XmlDictionaryReaderQuotas());

XElement root = XElement.Load(reader);

// The fields we'd like to extract
XElement form   = root.XPathSelectElement("//form");
XElement status = root.XPathSelectElement("//status");
XElement type   = root.XPathSelectElement("//type");

// Field set
IEnumerable<XElement> fields = root.XPathSelectElements("//fields/*");
```

Now our parser will work fine with an XML document that looks like this:

``` xml
<message>
    <form>13</form>
    <status>ACC</status>
    <type>NOOP</type>
    <fields>
        <field>
            <key>50</key>
            <value>1</value>
        </field>
    </fields>
</message>
```

Obviously the similarity of the two types of documents breaks down when we start using more complex features such as JSON's array support, or XML's namespaces; but it sure works great for simple cases.
