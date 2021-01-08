+++
location = "Calgary"
published_at = 2010-12-15T00:00:00-07:00
slug = "simultaneous-oracle-and-sql-server-support-in-entity-framework-with-designer-generated-objects"
tiny_slug = "35"
title = "Simultaneous Oracle and SQL Server Support in Entity Framework with Designer Generated Objects"
+++

The ADO.NET Entity Framework is a great product. Without going into too much background, it's Microsoft's attempt at an <acronym title="Object-relational Mapper">ORM</acronym> to make working with relational datastores more convenient. In their own words, it solves the _impedance mismatch across various data representations (for example objects and relational stores)_, and _empowers developers to focus on the needs of the application as opposed to the complexities of bridging disparate data source_. I've heard that problems in earlier versions of Entity made it's use quite a burden, but today's version seems to be quite a solid product.

Perhaps Entity's greatest feature it's LINQ-enabled so that it can transform a strongly typed LINQ statement into an SQL expression. These LINQ transformations even support advanced features like joins, aggregates, typecasts, and nullable types. For example, in the following method _all_ the logic gets offloaded to the database. That is, unless you swap out your database provider for an in-memory provider for testing, in that case, all the work would be done by the CLR instead.

``` cs
public IQueryable<LoginStatsViewModel> GetLoginStatsByTimeRange(
    DateTime startDate,
    DateTime endDate
)
{
    return
        from l in RepositoryFor<Login>().AsQueryable()
        where l.LoginDate >= startDate
            && l.LoginDate < endDate
        group l by l.UserKey into g
        join u in RepositoryFor<User>().AsQueryable()
            on g.Key equals u.UserKey
        let t = g.Sum(l => l.TimeSpentOnSite)
        select new LoginStatsViewModel()
        {
            Key                  = Guid.NewGuid(),
            UserKey              = g.Key,
            UserName             = u.UserName,
            NumLogins            = g.Count(),
            LastLoginDate        = (int)g.Max(l => l.LoginDate),
            UserTypeName         = u.UserType.UserTypeName,
            TotalNumCaptchaTries = (int)g.Sum(l => l.NumCaptchaTries),
            TotalTimeSpentOnSite = t.HasValue ? (double)t.Value : 0.0,
        };
}
```

Some other interesting things about this LINQ expression:

* `Guid.NewGuid()` will be translated by Entity into an SQL equivalent
* Explicit join on `User` as well as the implicit join on `UserType` using a navigation property
* Sums on both non-nullable and nullable columns
* The ternary operator (`?:`) works!
* This method is lazy, so the `IQueryable` object returned by it can be queried further before a final SQL statement is generated

Data objects such as `Login` and `User` can be generated automatically by giving the Entity designer a connection string for your SQL Server database, written by hand in a manner similar to what you might do with NHibernate (these are called <acronym title="Plain Old CLR Objects">POCO</acronym> entities), or both.

On the surface, Entity seems like a completely outstanding product that no .NET shop could afford to overlook, and it might be, save for one serious drawback.

Data Providers
--------------

ADO.NET is based on a data provider model which provides a common managed interface that can implemented by data providers to enable interaction with various data stores. A default provider for SQL Server is included in the framework by default, and implementations for any third party databases are left up to, well, third parties. Entity Framework being fairly mature, a [number of data providers](http://msdn.microsoft.com/en-us/data/dd363565.aspx) are already available, and as is common in the Microsoft ecosystem, most of them for a price.

The problem with this arrangement is twofold:

1. SQL Server is nice, and many companies working with .NET try to use it if possible, largely on account of Microsoft's excellent tools support. In practice however, Oracle is extremely pervasive in big business, and many shops require backend support for it.

2. If enabling a data provider was as simple as just dropping in a few DLLs, it would be worth paying another company to save your time. However, the reality of the situation is that getting your existing Entity infrastructure to support a new database will require non-trivial redesign and debugging.

Compare this situation to the open-source alternative for a moment. By changing a single line of configuration in a Ruby on Rails project, I can get support for MySQL, PostgreSQL, SQLite, DB2, SQL Server, Oracle, and other relational databases. By switching out my `ActiveRecord` implementation, I can bring in support for the world of non-relational data stores: CouchDB, MongoDB, Redis, and a host of others. Obviously Ruby's dynamic nature helps a bit here, but the main point is that Entity's multi-platform support has huge room for improvement.

Our goal was to maintain an Entity infrastructure based around Microsoft's default SQL Server provider so that we could leverage its nice Visual Studio support, easy integration of new features released down the road, generation of data objects given a database, and most of all because life is usually made easier in the Microsoft ecosystem by preferring a Microsoft product over a third party alternative. In addition, Oracle support was a requirement due to its popularity among our client base. We strongly preferred an arrangement that would allow us to maintain a single set of data objects for all databases so that duplicated infrastructure wouldn't be necessary.

Devart
------

We tried a number of possible solutions and in the end concluded that the best we could do was to use a backend called [dotConnect created by Devart](http://www.devart.com/dotconnect/). Devart has its own set of tools for building an Entity backend, but we identified that we could also used it for Oracle connectivity with Microsoft tools using a process like the following:

* Ship Devart provider DLLs for Oracle
* Create an alternate `edmx` to tell Devart how to map object properties to database fields, among other things. In general the SQL Server `edmx` would be very similar to the Oracle `edmx`. We'd create the first one by hand, then build a generator to take care of it automatically in the future. Note that an `edmx` is often defined in its components: `csdl`, `msl`, and `ssdl`.
* Create a new Entity Framework configuration string pointing to the Devart Oracle provider and the new `edmx` files

In theory, the whole process is pretty straightforward and should be almost "drop in", involving only some modifications to database mapping types in the `ssdl` component of the `edmx`. In practice though, we found that working with the Devart provider came with a number of gotchas, and we spent the equivalent of one man-week finding suitable solutions for them. Devart is functional, but its support and implementation aren't quite there yet. Oracle is supposed to be releasing a beta for its official ADO.NET provider sometime in 2011, and we'll probably evaluate moving to it when the time comes.

Below is a rough roadmap of what needs to be done to ship Oracle support.

### Installation

The [dotConnect for Oracle](http://www.devart.com/dotconnect/oracle/download.html) package needs to be downloaded and installed in order to get the requisite provider DLLs. Only one person really needs to do this because the DLLs can be checked into the working solution from there and made available to the rest of the team.

A word of warning: after installing Devart, you may be prompted to install the Devart designer tools the next time Visual Studio starts up. I refused this install request, and Devart responded by severely crippling my Visual Studio 2010 installation; symptoms being that all toolbars and panes refused to open, some delivering an error code. It wasn't obvious that Devart was causing the problem, but I uninstalled it and the issues disappeared. I was able to reinstall Devart in a safe manner by carefully disabling any Visual Studio related components during the install process, _and_ refusing the designer tools during Visual Studio's startup.

### DLLs

For full Oracle provider support, three DLLs need to be shipped with an Entity-enabled project: `Devart.Data.dll`, `Devart.Data.Oracle.dll`, and `Devart.Data.Oracle.Entity.dll`. These should be available in the dotConnect installation directory. An appropriate license for these files will also be necessary.

### SSDL and Devart's Special Case Behavior

SSDL is short for _store schema definition language_ and is one of the three major components generated as part of an `edmx`. Its primary reponsibilities are defining associations between entities, the database table corresponding to each entity, and the database types for each entity property.

Here's a small piece of a sample SSDL:

``` xml
<Schema ...>
  <EntityType Name="User">
    <Key>
      <PropertyRef Name="USERKEY" />
    </Key>
    <Property Name="USERKEY" Type="CHAR" Nullable="false" MaxLength="36" />
    <Property Name="USERNAME" Type="VARCHAR2" Nullable="false" MaxLength="50" />
    <Property Name="DESCRIPTION" Type="VARCHAR2" MaxLength="200" />
    ...
  </EntityType>
</Schema>
```

The Oracle SSDL will look very similar to the SQL Server SSDL, with a few key differences:

* The top-level `<Schema>` tag's `Provider` attribute must be set to `Devart.Data.Oracle`
* The `<Schema>` tag's `ProviderManifestToken` attribute must be set to `ORA`
* Any `<EntitySet>` tags must have the value of their `Schema` attribute _uppercased_
* Any `<EntitySet>` tags should have a `Table` attribute added to them, where the value is the name of that entity's database table _uppercased_
* All `<Property>` and `<PropertyRef>` tags must have the value of their `Name` attribute _uppercased_
* All `<Property>` tags must have their `Type` attribute changed to an Oracle equivalent, see the table below

There's a good reason that we uppercase many of these values. As of now (December 2010), the Devart provider will wrap any table and column names that are not uppercase in double quotes when building SQL queries. This is unfortunate, because Oracle doesn't support these objects being wrapped in quotes. This problem can be worked around using uppercase for all database object names.

<table>
    <caption>SQL Server types and their Oracle equivalents (incomplete table)</caption>
    <tr>
        <th style="width: 50%;">SQL Server Type</th>
        <th style="width: 50%;">Oracle Type (case sensitive)</th>
    <tr>
    <tr>
        <td>bigint</td>
        <td>int64</td>
    <tr>
    <tr>
        <td>char</td>
        <td>CHAR</td>
    <tr>
    <tr>
        <td>datetime</td>
        <td>DATE</td>
    <tr>
    <tr>
        <td>numeric</td>
        <td>decimal</td>
    <tr>
    <tr>
        <td>tinyint</td>
        <td>byte</td>
    <tr>
    <tr>
        <td>varchar</td>
        <td>VARCHAR2</td>
    <tr>
</table>

### MSL

MSL is short for _mapping specification language_ and is the next major `edmx` component. Its job is to map the properties of an entity's class representation to the fields of the corresponding database table.

Here's a small piece of a sample MSL:

``` xml
<Mapping Space="C-S" xmlns="http://schemas.microsoft.com/ado/2008/09/mapping/cs">
  <EntitySetMapping Name="User">
    <EntityTypeMapping TypeName="Web_Model.User">
      <MappingFragment StoreEntitySet="User">
        <ScalarProperty Name="UserKey" ColumnName="USERKEY" />
        <ScalarProperty Name="UserName" ColumnName="USERNAME" />
        <ScalarProperty Name="Description" ColumnName="DESCRIPTION" />
        ...
      </MappingFragment>
    </EntityTypeMapping>
  </EntitySetMapping>
</Mapping>
```

The Oracle MSL needs fewer changes than the SSDL, and would probably require none at all if Devart's double quoting bug didn't exist:

* All `<*Property>` tags must have the value of their `ColumnName` attribute _uppercased_

### CSDL

CSDL is short for _conceptual schema definition language_ and is the last `edmx` component. It tracks relations, general property information, and navigation properties. Here's a sample section from a CSDL:

``` xml
<Schema Namespace="Web_Model" Alias="Self" xmlns:annotation="http://schemas.microsoft.com/ado/2009/02/edm/annotation" xmlns="http://schemas.microsoft.com/ado/2008/09/edm">
  <EntityType Name="User">
    <Key>
      <PropertyRef Name="UserKey" />
    </Key>
    <Property Type="String" Name="UserKey" Nullable="false" MaxLength="36" FixedLength="true" Unicode="false" />
    <Property Type="String" Name="UserName" Nullable="false" MaxLength="50" FixedLength="false" Unicode="false" />
    <Property Type="String" Name="Description" MaxLength="200" FixedLength="false" Unicode="false" />
    ...
  </Entitytype>
</Schema>
```

No changes are required to the CSDL for Oracle. We were hoping to simply reuse the CSDL from our SQL Server `edmx` file, but we ended up duplicating it because we ran into a compiler bug that would occur intermittently. We're still not exactly sure what the problem was, but the compilation error looked something like this:

```
Error   19      Value cannot be null.
Parameter name: csdlPath        Company.Namespace.Project
```

### Resources and Config Paths

After new SSDL, MSL, and CSDL files have been generated, they should be added to the entity project and their _build action_ changed to _embedded resource_. You'll want to reference these new resources from your Oracle connection string. Keep in mind that the compiler will probably prepend your full namespace path to their resource names, so reference them accordingly. If in doubt, use Reflector to check what path a resource has been embedded in a DLL under.

Here's an example Oracle Entity connection string for use in `App/Web.config`:

``` xml
<add name="Web_Entities"
     connectionString="provider=Devart.Data.Oracle;metadata=res://*/Company.Namespace.Project.WebDataModel.Oracle.csdl|res://*/Company.Namespace.Project.WebDataModel.Oracle.ssdl|res://*/Company.Namespace.Project.WebDataModel.Oracle.msl;Provider Connection String='Data Source=myoracleserver/xe;User Id=myoracleuser;Password=myoraclepassword;'"
     providerName="System.Data.EntityClient" />
```

### Conversion 

The final step is a conversion project that will allow you to modify your Entity's `edmx` and convert those changes to an equivalent Oracle version. I've included some basic ideas for writing such a program built around XPath (I'm using my own extension methods here so this code won't compile directly).

``` cs
#region Public

public static void Generate(string sourceSchemaPath, string targetPathAndRootName)
{
    XmlReader reader = XmlReader.Create(sourceSchemaPath);
    XElement root = XElement.Load(reader);

    XmlNamespaceManager ns = new XmlNamespaceManager(reader.NameTable);
    ns.AddNamespace("cs", "http://schemas.microsoft.com/ado/2008/09/mapping/cs");
    ns.AddNamespace("edmx", "http://schemas.microsoft.com/ado/2008/10/edmx");
    ns.AddNamespace("ssdl", "http://schemas.microsoft.com/ado/2009/02/edm/ssdl");
    ns.AddNamespace("edm", "http://schemas.microsoft.com/ado/2008/09/edm");

    XElement msl = GenerateMSL(root, ns);
    msl.Save(targetPathRootName + ".msl");
    Console.Out.WriteLine("Generated: " + targetPathRootName + ".msl");

    XElement ssdl = GenerateSSDL(root, ns);
    ssdl.Save(targetPathRootName + ".ssdl");
    Console.Out.WriteLine("Generated: " + targetPathRootName + ".ssdl");

    XElement csdl = GenerateCSDL(root, ns);
    csdl.Save(targetPathRootName + ".csdl");
    Console.Out.WriteLine("Generated: " + targetPathRootName + ".csdl");
}

#endregion

#region Private

private static XElement GenerateMSL(XElement root, XmlNamespaceManager ns)
{
    foreach (XElement property
        in root.XPathSelectElements(
            "//edmx:Runtime/edmx:Mappings/cs:Mapping/cs:EntityContainerMapping/" + 
            "cs:EntitySetMapping/cs:EntityTypeMapping/cs:MappingFragment/*", 
            ns
        ))
    {
        property.ModifyAttribute("ColumnName", v => v.ToUpper());
    }

    XElement msl = 
        root.XPathSelectElement("//edmx:Runtime/edmx:Mappings/cs:Mapping", ns);
    return msl;
}

private static XElement GenerateCSDL(XElement root, XmlNamespaceManager ns)
{
    // No modification -- see explanation above
    XElement csdl = 
        root.XPathSelectElement(
            "//edmx:Runtime/edmx:ConceptualModels/edm:Schema", 
            ns
        );
    return csdl;
}

private static XElement GenerateSSDL(XElement root, XmlNamespaceManager ns)
{
    foreach (XElement entitySet
        in root.XPathSelectElements(
            "//edmx:Runtime/edmx:StorageModels/ssdl:Schema/" + 
            "ssdl:EntityContainer/ssdl:EntitySet",
            ns
        ))
    {
        entitySet.ModifyAttribute("Schema", v => v.ToUpper());
        entitySet.SetAttribute("Table", entitySet.Attribute("Name").Value.ToUpper());
    }

    foreach (XElement property
        in root.XPathSelectElements("//*/ssdl:Property", ns))
    {
        property.ModifyAttribute("Name", v => v.ToUpper());

        Action<string> setAttr = v => property.SetAttribute("Type", v);

        // Not an exhaustive list
        switch (property.Attribute("Type").Value.ToLower())
        {
            case "bigint":
                setAttr("int64"); break;
            case "tinyint":
                setAttr("byte"); break;
            case "char":
                setAttr("CHAR"); break;
            case "datetime":
                setAttr("DATE"); break;
            case "numeric":
                setAttr("decimal"); break;
            case "text":
                setAttr("CLOB"); break;
            case "varchar":
                setAttr("VARCHAR2"); break;
            default:
                Console.WriteLine(
                    "Warning: unrecognized property type: " + 
                    property.Attribute("Type").Value
                );
                break;
        }
    }

    foreach (XElement propertyRef
        in root.XPathSelectElements("//*/ssdl:PropertyRef", ns))
    {
        propertyRef.ModifyAttribute("Name", v => v.ToUpper());
    }

    XElement ssdl = 
        root.XPathSelectElement("//edmx:Runtime/edmx:StorageModels/ssdl:Schema", ns);
    ssdl.SetAttribute("Provider", "Devart.Data.Oracle");
    ssdl.SetAttribute("ProviderManifestToken", "ORA");
    return ssdl;
}

#endregion
```
