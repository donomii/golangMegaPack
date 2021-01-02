[![Build Status](https://travis-ci.org/donomii/portprog.svg?branch=master)](https://travis-ci.org/donomii/portprog)

# [The Portable Programmer](https://donomii.github.io/portprog/)

Programming languages and libraries for the programmer on the move.  [Read the website](https://donomii.github.io/portprog/)

Double click to install GCC, Java, Perl, TclTk, Nim, Lazarus, Squeak Smalltalk, make, maven, Eclipse, IntelliJ, Codeblocks, Notepad++, and many supporting tools and libraries.

After installation, run environment.bat to open a command shell with all your paths correctly set up.

## Easy setup

I frequently have to install my programming tools on a fresh computer.
  Sometimes it's a new laptop, but usually it's a fresh VM or a computer at a client site.
  And it is always frustrating, because it takes me hours to track down every minor library and
  patch that I need to get something compiled.
  
Of course, there are no _minor_ libraries, because missing any of them halts my compile.  So
I put together this installer to get my build environment set up as quickly as possible.

# Installation

## Windows

Download a new release from the [Releases page](https://github.com/donomii/portprog/releases).  Unpack it and double click the exe.

## Linux and Mac

	go get -u github.com/donomii/portprog
	go build
	./portprog

## Options

There aren't a lot.  This isn't another distribution, it's just a fancy downloader and unpacker.  There's no dependency management or build flags 
or whatever.

	--no-gcc	Don't download or install gcc
	--no-golang 	Don't download or install golang
	--no-git 	Don't attempt to clone or update any repositories via git
	--devel		Only install packages from the packages-devel directory
        --file          Install using settings from file   
	--no-install    Don't install anything, but do regenerate the environment scripts

## Operation

Portprog starts and checks the packages (or packages-windows) directory, then attempts to download all the files to the zips directory, unpacking them in the background.  It will not download any zips already present in the zips directory.

## Adding your own

The goal of PortProg is to manage your own downloads.  You can easily add any download that you want:

*  Follow the [simple instructions](MakeYourOwnPackage.md)
*  Send me a pull request so I can add your package to the distribution!


## Uninstall

Delete the directory.  Portprog does not modify any part of your system outside of its own directory.

** Warning **

I can't control other programs and libraries, so when you use a program or library that portprog downloads for you, this third party might change your system.  

I can't stop that, but I try not to use any programs that would do that.  You might end up with some files left in Application Data and the registry

