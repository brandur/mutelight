+++
location = "Calgary"
published_at = 2010-07-15T00:00:00-06:00
slug = "building-a-command-line-environment-for-net-development-with-rake"
tiny_slug = "25"
title = "Building a Command Line Environment for .NET Development with Rake"
+++

While developing .NET projects, it's useful to have a set of simple tasks that are are runnable on the command line (rather than through Visual Studio) for a few reasons:

* You can build your solution and run its tests without opening <acronym title="Visual Studio">VS</acronym>
* During development, you can run your test suite without going through VS, which is tempermental and can be slow to respond to input
* In the case of ASP.NET, you can start a development server to view your application
* For free, you get easy commands for use with your <acronym title="Continuous Integration">CI</acronym> server

History
-------

The defacto standard for building .NET from the command line used to be a project called [NAnt](http://en.wikipedia.org/wiki/Nant), a port from Java's Ant that built projects and ran tasks based on descriptions in build XML files. A huge drawback to NAnt was that even though VS already knew how to build your solution, that information was locked up in an `.sln` file, and you'd have to define a separate set of instructions for NAnt, thus forcing you to maintain the build in two places. NAnt build files were also extremely verbose, and would become difficult to maintain for large projects.

Microsoft later introduced their own build tool called [MSBuild](http://en.wikipedia.org/wiki/Msbuild). In its most basic form, MSBuild is an executable that reads a VS solution file, and builds it for you. It also supports more sophisticated tasks that can be defined in MSBuild project files, these files have XML syntax similar to NAnt's. 

They key to MSBuild's usefulness was that you'd give it an output target, and it would resolve all the necessary dependencies for you by reading your solution; no more manually building intermediary DLLs with `csc`. It even became fairly common to use a NAnt build file, but to build using an `<exec>` call to MSBuild instead of using built-in NAnt compilation tasks. In fact, I wrote and [building with MSBuild and NAnt](http://mutelight.org/articles/the-anatomy-of-a-nant-build-file.html) last year.

The problem with both NAnt and MSBuild projects is that they're XML, a markup language that is fundamentally hard to read and write for humans. A new movement appeared that started to use Ruby Rakefiles to to build .NET solutions.

Rake
----

[Rake](http://en.wikipedia.org/wiki/Rakefile) is a build tool commonly used for building projects in the Ruby world. One nice thing about it is that it is both written in Ruby, *and* allows the user to write build tasks in Ruby, keeping Rakefiles easy to read, and massively powerful with access to the complete Ruby language and all its Gems.

I'm going to provide a small walkthrough on how to get a Ruby environment installed, and getting a Rakefile up and running for a project:

1. Install [Cygwin](http://www.cygwin.com/), a Linux shell for Windows. When walking through the installer, select the *Devel* tree which includes Ruby.

2. Open a Cygwin shell

3. Download the Rubygems tarball from [Rubyforge](http://rubyforge.org/projects/rubygems/), unpack it, and from Cygwin run `ruby setup.rb install` (your home directory in Cygwin is at `C:\cygwin\home\fyrerise` by default, a good choice of download location for the tarball)

4. Install Rake with `gem install rake` and Bundler with `gem install bundler`

5. In Cygwin, navigate to your .NET project's path (get to a Windows path using `cd /cygdrive/c/Users/...`)

Bundler
-------

Ruby Gems are small software libraries that extend the core language. In most Ruby projects, a large number of Gems are usually needed, and to help manage these dependencies, a package called [Bundler](http://gembundler.com/) was written.

We'll only need a few Gems for our simple build file, but it's best to use Bundler anyway because it's easy and our complexity may increase later on. Project Gem dependencies are tracked in a file called `Gemfile`, create one in your solution's root with the following contents:

``` ruby
source 'http://rubygems.org'

gem 'albacore'
gem 'haml' # Includes Sass
```

Run `bundle install .` (from Cygwin). Two dependencies are now installed:

* **Albacore** &mdash; provides a set of .NET build tasks for Rake
* **Haml** &mdash; allows us to access the [Sass](http://sass-lang.com/) compiler, which I use for my ASP.NET development for a more literate CSS

Our Rakefile
------------

As promised, now it's time for our build file. Create `Rakefile` with these contents (replacing solution, and test project paths appropriately):

``` ruby
# Initialize the Bundler environment, it will handle all other dependencies
require 'rubygems'
require 'bundler'
Bundler.setup

require 'albacore'
require 'mstest_task'
require 'sass/plugin'

# Albacore still defaults to MSBuild 3.5, so specify the exe location manually
MsBuild = 'C:/Windows/Microsoft.NET/Framework/v4.0.30319/MSBuild.exe'
WebDev  = 'C:/Program Files/Common Files/microsoft shared/DevServer/10.0/WebDev.WebServer40.EXE'

task :default => :build

desc 'Alias for build:debug'
task :build => 'build:debug'

namespace :build do
  [ :debug, :release ].each do |t|
    desc "Build the Project solution with #{t} configuration"
    msbuild t do |b|
      b.path_to_command = MsBuild
      b.properties :configuration => t
      b.solution = 'Project/Project.sln'
      b.targets :Build
    end
  end
end

desc 'Clean build files from the directory structure'
msbuild :clean do |b|
  b.path_to_command = MsBuild
  b.solution = 'Project/Project.sln'
  b.targets :Clean
end

desc 'Alias for test:all'
task :test => 'test:all'

namespace :test do
  desc 'Run all tests'
  # Run the category task with no parameters, and therefore no category
  task :all => 'test:category'

  # Usage -- rake test:category[<category name>]
  desc 'Run all tests in a category'
  mstest :category, :cat, :needs => 'build:debug' do |t, args|
    t.category = args[:cat] if args[:cat]
    t.container = 'Project/Project.Tests/bin/Debug/Project.Tests.dll'
  end
end

desc 'Start a development server'
exec :server => 'build:debug' do |cmd|
  cmd.path_to_command = WebDev
  path = 'Project/Project/'
  # WebDev.WebServer is *extremely* finicky and is a typical example of 
  # fragile Microsoft coding. For it to work, its path MUST (a) use 
  # backslash directory separators, and (b) be absolute. Here we use Cygwin 
  # to convert a relative Unix path to an absolute Windows path.
  path = `cygpath -a -w #{path}`.strip
  cmd.parameters << %-"/path:#{path}"-
  cmd.parameters << %-"/port:3001"-
  puts ''
  puts 'Starting development server on http://localhost:3001'
  puts 'Ctrl-C to shutdown server'
end

desc 'Updates stylesheets if necessary from their Sass templates'
task :sass do
  # @todo: change this once we know our stylesheets location
  Sass::Plugin.add_template_location '.'
  Sass::Plugin.on_updating_stylesheet do |template, css|
    puts "Compiling #{template} to #{css}"
  end
  Sass::Plugin.update_stylesheets
end
```

Run `rake -T` for a list of available tasks. Here are the important ones:

* `rake build` &mdash; build the project
* `rake build:release` &mdash; build the project with release configuration
* `rake test` &mdash; run our test suite (I'm using the MSTest framework here)
* `rake server` &mdash; start a development server pointing to our project for ASP.NET
