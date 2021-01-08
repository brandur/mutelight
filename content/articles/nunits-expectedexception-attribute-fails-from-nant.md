+++
location = "Calgary"
published_at = 2009-06-22T00:00:00-06:00
slug = "nunits-expectedexception-attribute-fails-from-nant"
tiny_slug = "12"
title = "NUnit's ExpectedException Attribute Fails from NAnt"
+++

It seems like there is a bug in NAnt's `<nunit2>` task that occurs when using NUnit's `ExpectedException` attribute. There seems to be [an open bug](http://sourceforge.net/tracker/?func=detail&atid=402868&aid=1942863&group_id=31650) related to this problem, but it's status has been unchanged for some time now (since April 2008). I'm using the NAnt 0.86 beta and NUnit 2.5. Here's a simple test for this problem:

``` csharp
[Test]
[ExpectedException(typeof(Exception))]
public void TestThatExpectedExceptionFails()
{
    throw new Exception("wef");
}
```

And a NAnt file to build the test's project it and try to run it:

``` xml
<?xml version="1.0"?>

<project default="default">
    <property name="msbuild" value="${framework::get-framework-directory(framework::get-target-framework())}\MSBuild.exe" />
    <target name="default">
        <exec program="${msbuild}">
            <arg value="Proj.csproj" />
            <arg value="/target:Rebuild" />
            <arg value="/property:Configuration=Debug" />
            <arg value="/verbosity:quiet" />
        </exec>
        <nunit2>
            <formatter type="Plain" />
            <test assemblyname="bin\Debug\Proj.dll" />
        </nunit2>
    </target>
</project> 
```

Here's what NAnt has to say about that:

```
NAnt 0.86 (Build 0.86.2898.0; beta1; 12/8/2007)
Copyright (C) 2001-2007 Gerry Shaw
http://nant.sourceforge.net

Buildfile: file:///C:/Temp/NAntTest/NAntTest/nant.build
Target framework: Microsoft .NET Framework 3.5
Target(s) specified: default


default:

    [exec] Microsoft (R) Build Engine Version 3.5.30729.1
    [exec] [Microsoft .NET Framework, Version 
2.0.50727.3053]
    [exec] Copyright (C) Microsoft Corporation 2007. All 
rights reserved.
    [exec]
  [nunit2] Tests run: 1, Failures: 1, Not run: 0, Time: 
0.031 seconds
  [nunit2]
  [nunit2] Failures:
  [nunit2] 1) Proj.Tests.TestThatExpectedExceptionFails : 
System.Exception : wef
  [nunit2]    at Proj.Tests.TestThatExpectedExceptionFails()
in c:\Temp\NAntTest\NAntTest\Tests.cs:line 17
  [nunit2]
  [nunit2]
  [nunit2]

BUILD FAILED

C:\Temp\NAntTest\NAntTest\nant.build(11,6):
Tests Failed.

Total time: 1 seconds.
```

Use Assert.Throws<> Instead
---------------------------

Prior to finding this bug, I used to use `ExpectedException` fairly regularly. It wasn't until I started looking for alternatives that work in NAnt that I realized a fairly persuasive argument could be built that `ExpectedException` isn't the best tool for testing exceptions anyway.

Here are my talking points regarding what `Assert.Throws<>` has over `ExpectedException`:

1. **Readability:** when reading tests that rely on `ExpectedException`, its easy to miss the attribute declarations, which might leave you confused as to what the test was supposed to accomplish upon reading it through. Tests using asserts read cleanly top to bottom.
2. **Test multiple conditions:** a test using asserts can test for multiple thrown exceptions, while a test using `ExpectedException` can only test for one, forcing you to write multiple tests to cover all conditions.
3. **Type parameters look nice:** this one is pretty minor, but using the exception assert allows you to leave out the `typeof()` operator.

An example that demonstrates the advantages of `Assert.Throws<>` follows:

``` csharp
private void MethodWithArgChecks(string arg1, string arg2)
{
    if (arg1 == null)
        throw new ArgumentException("arg1 = null", "arg1");

    if (arg2 == null)
        throw new ArgumentException("arg2 = null", "arg2");
}

[Test]
public void TestForThrownExceptionOnBadArgs()
{
    Assert.Throws<ArgumentException>(
        () => MethodWithArgChecks(null, "ok")
    );
    Assert.Throws<ArgumentException>(
        () => MethodWithArgChecks("ok", null)
    );
}
```
