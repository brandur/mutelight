+++
location = "Calgary"
published_at = 2009-06-02T00:00:00-06:00
slug = "the-anatomy-of-a-nant-build-file"
tiny_slug = "10"
title = "The Anatomy of a NAnt Build File"
+++

Over the weekend I was looking at building a simple NAnt file layout that I could use for building one of my projects. Previously, I'd considered NAnt files to be problematical because they have a tendency of evolving into unmanageable behemoths that are thousands of lines long. For example, the build file we use at my day job is a little over 2600 lines and has a layout like this:

``` xml
<project name="Unmanageable" default="build">

    ...

    <target name="build">
        <call target="buildProject1" />
        <call target="buildProject2" />
        <call target="buildProject3" />
    </target>

    ...

    <target name="buildProject1">
        <csc target="library" output="buildProject1.dll">
            <sources>
                <include name="**/*.cs" />
            </sources>
            <references>
                <include name="System.dll" />
                ...
            </references>
        </csc>
    </target>
    <target name="buildProject2" depends="buildProject1">
        ...
    </target>
    <target name="buildProject3" 
            depends="buildProject1, buildProject2">
        ...
    </target>
</project>
```

Add in 30 or 40 other projects and about 100 lines worth of property definitions and before you know it you've got a file that will turn developers pale at its very mention. That might not be so bad, but because of the way references have to be entered manually under the `<csc>` tag, everytime you update a reference under one of your projects, it has to be added to the build file as well; so the build file ends up being updated constantly.

Not wanting to run into this problem in my own projects, I started looking for a better way to build projects, and quickly discovered the NAnt [`<solution>` task (NAnt solution task documentation)](http://nant.sourceforge.net/nightly/latest/help/tasks/solution.html). This task can theoretically read a solution or project file, build dependencies, then build the file's target; thus working around the normal NAnt "creep". Sounds great, but it only supports VS 2002 and 2003 solution files, so in practice it's not much use.

MSBuild
-------

With the release of Visual Studio 2005, Microsoft started to include a new build tool called [MSBuild (MSBuild documentation)](http://msdn.microsoft.com/en-us/library/0k6kkbsd.aspx). In its most basic form, the MSBuild command can be used with a few command line switches to build a VS solution or project: `MSBuild "path/Project.csproj" /v:n /t:Build /p:Configuration=release;OutDir="path/bin/"`

MSBuild is more than just a command though; it can also be used to read an MSBuild file written in XML to follow much more complex sets of instructions in a similar fashion to NAnt. I haven't had the chance to explore it fully, but I hope to do a future post comparing its features and ease-of-use to NAnt's.

For now, I found that calling the MSBuild executable from my NAnt scripts acts like a drop-in replacement for NAnt's `<solution>` task with one important difference &mdash; MSBuild works.

A NAnt Solution
---------------

Without further ado, I'd like to present a template that I've been using for my projects that can build multiple application targets using MSBuild, and as you'd expect, has testing targets.

``` xml
<project name="BrandursExcellentTemplate" default="rebuild">

    <!-- ============= -->
    <!-- Configuration -->
    <!-- ============= -->

    <!-- Build as: release, debug, etc. -->
    <property name="configuration" value="release" />

    <!-- Output directory where our executables should be written to -->
    <property name="bin-directory" value="${directory::get-current-directory()}/bin/" />

    <!-- Location of the MSBuild executable, we use this to build projects -->
    <property name="msbuild" value="${framework::get-framework-directory(framework::get-target-framework())}\MSBuild.exe" />

    <!-- ============ -->
    <!-- Main Targets -->
    <!-- ============ -->

    <target name="clean" description="Delete all previously compiled binaries.">
        <delete>
            <fileset>
                <include name="**/bin/**" />
                <include name="**/obj/**" />
                <include name="**/*.suo" />
                <include name="**/*.user" />
            </fileset>
        </delete>
    </target>

    <target name="build" description="Build all application targets.">
        <mkdir dir="${bin-directory}" />
        <!-- Neither of these are secondary projects like   -->
        <!-- libraries, they are executable projects. Their -->
        <!-- dependency projects will be built for us by    -->
        <!-- MSBuild automatically.                         -->
        <call target="build.app1" />
        <call target="build.app2" />
    </target>

    <target name="rebuild" depends="clean, build" />

    <target name="test" description="Build test project and run all tests.">
        <mkdir dir="${bin-directory}" />
        <call target="build.tests" />
        <nunit2>
            <formatter type="Plain" />
            <test assemblyname="${bin-directory}/Test.dll" />
        </nunit2>
    </target>

    <target name="testimmediate" description="Build test project and run all tests.">
        <mkdir dir="${bin-directory}" />
        <call target="build.tests" />
        <nunit2>
            <formatter type="Plain" />
            <test>
                <assemblies>
                    <include name="${bin-directory}/Test.dll" />
                </assemblies>
                <categories>
                    <include name="Immediate" />
                </categories>
            </test>
        </nunit2>
    </target>

    <!-- ================= -->
    <!-- Secondary Targets -->
    <!-- ================= -->

    <target name="build.app1">
        <exec program="${msbuild}" commandline='"src/App1/App1.csproj" /v:n /nologo /t:Build /p:Configuration=${configuration};OutDir="${directory::get-current-directory()}/bin/"' />
    </target>

    <target name="build.app2">
        <exec program="${msbuild}" commandline='"src/App2/App2.csproj" /v:n /nologo /t:Build /p:Configuration=${configuration};OutDir="${directory::get-current-directory()}/bin/"' />
    </target>

    <target name="build.tests">
        <!-- Do not build verbosely (/v:q), user wants to see test results, not build output -->
        <exec program="${msbuild}" commandline='"src/Test/Test.csproj" /v:q /nologo /t:Build /p:Configuration=Debug;OutDir="${directory::get-current-directory()}/bin/"' />
    </target>

</project>
```

One part of the template I'll comment on is the `testimmediate` target. At work, we used to categorize all our tests using NUnit's category feature, but found over time that we weren't gaining a whole lot of value from the process. Instead of letting categories go to waste, my boss came up with the idea of tagging test fixtures relevant to current features with `[Category("Immediate")]`. That way, developers can run only the immediate category to quickly run the most immediately important tests. I've frequently started using the same concept outside of work as well.
