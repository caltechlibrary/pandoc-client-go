/*
Copyright (c) 2022, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/pandoc_client"
)

var (
	helpText = `% {app_name}(1) user manual
% R. S. Doiel
% 2022-11-04

# NAME

{app_name}

# SYNOPSIS

~~~
{app_name} [OPTIONS] CONFIG_JSON HTDOCS
~~~

# DESCRIPTION

{app_name} convert files Markdown files found in ` + "`" + `HTDOC_DIR` + "`" + ` to
HTML files. Two required parameters are ` + "`" + `CONFIG_JSON` + "`" + ` file
and HTDOCS directory containing Markdown documents. The
configuration JSON file should include any parameters need to 
format the "POST" sent to the [Pandoc Server](https://pandoc.org/pandoc-server.html)
(see the API documentaiton). 

The HTDOCS directory path will be recusively walked to find
files ending in ".md" and write successfull conversions to 
the same file path using a ".html" extension instead of ".md".

# OPTIONS

-help
: display help

-version
: display version

-license
: display license

-verbose
: use verbose log output

# EXAMPLE

In this example we have markdown files in a directory structure
called ` + "`" + `/var/www/htdocs` + "`" + `. We're using a config.json file contains

~~~
{
	"from": "markdown",
	"to": "html5",
	"standalone": true
}
~~~


The command to convert the Markdown files to HTML is 

~~~
{app_name} config.json /var/www/htdocs
~~~

The ` + "`" + `/var/www/htdocs` + "`" + ` directory needs to have write permission
by the user running ` + "`" + `{app_name}` + "`" + `. As the Markdown files are encountered
a log message will be written indicating any errors or that the file
was successful converted.

`
)

func usage(appName string) string {
	return strings.ReplaceAll(helpText, "{app_name}", appName)
}

func main() {
	appName := path.Base(os.Args[0])
	showHelp, showVersion, showLicense := false, false, false
	verbose := false
	flag.BoolVar(&showHelp, "help", showHelp, "display help")
	flag.BoolVar(&showVersion, "version", showVersion, "display version")
	flag.BoolVar(&showLicense, "license", showLicense, "display license")
	flag.BoolVar(&verbose, "verbose", verbose, "verbose log output")
	flag.Parse()

	if showHelp {
		fmt.Fprintf(os.Stdout, "%s\n", usage(appName))
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(os.Stdout, "%s %s\n", appName, pandoc_client.Version)
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(os.Stdout, "%s\n", pandoc_client.License)
		fmt.Fprintf(os.Stdout, "%s %s\n", appName, pandoc_client.Version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "ERROR: expected a json configuration filename and htdocs path\n")
		os.Exit(1)
	}
	cfg, err := pandoc_client.Load(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	cfg.Verbose = verbose
	cfg.From = "markdown"
	cfg.To = "html5"
	if err := cfg.Walk(args[1], ".md", ".html"); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
