% md2html(1) user manual
% R. S. Doiel
% 2022-11-04

# NAME

md2html

# SYNOPSIS

```
md2html [OPTIONS] CONFIG_JSON HTDOCS
```


# DESCRIPTION

md2html convert files Markdown files found in `HTDOC_DIR` to
HTML files. Two required parameters are `CONFIG_JSON` file
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

# EXAMPLE

In this example we have markdown files in a directory structure
called `/var/www/htdocs`. We're using a config.json file contains

```json

```

The command to convert the Markdown files to HTML is 

```shel
md2html config.json /var/www/htdocs
```

The `/var/www/htdocs` directory needs to have write permission
by the user running `md2html`. As the Markdown files are encountered
a log message will be written indicating any errors or that the file
was successful converted.



