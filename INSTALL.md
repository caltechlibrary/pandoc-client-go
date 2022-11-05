INSTALL
=======

This is an experimental Go package and cli. It is only distributed in source code. If you wish compile and test it then you need the required development environment.  Follow the steps listed below in "Compiling from Source".

Requirements
------------

pandoc_client is currently implemented as a Go application. This may change in the future.

- Git to clone the repository
- Compiling the cli
    - [Golang](https://golang.org) 1.19.2 or better
    - GNU Make
    - Pandoc 2.19 or better (you need to run pandoc-server for the client to work)

Compiling from Source
---------------------

1. clone the repository
2. change into the cloned directory
3. run "make", "make test", "make install" to install the cli, man pages

Here's the steps I take on my macOS box or Linux box.

~~~
git clone git@github.com:caltechlibrary/pandoc_client.git
cd pandoc_client
make
make test
make install
~~~

