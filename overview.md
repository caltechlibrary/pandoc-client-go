

Overview
========

This Go package demonstrates who you can interact with [Pandoc Server](https://pandoc.org/pandoc-server.html)
from a Go program.

The provided proof of concept is called [md2html](md2html.1.html). It will
walk a directory finding the ".md" files and convert them to ".html" files
based on the defaults specified in a JSON configuration file.

Concept
-------

The pandoc server is launched via systemd when the machine starts up. It listens on a
localhost port ONLY.  When the site building processes startup they can read the 
documents that need to be converted from either a database (e.g. SQLite3, MySQL,
Postgres) or the file system. That content is then turned into a structure that
the Pandoc server understands and is sent to it as a POST per [Pandoc Server](https://pandoc.org/pandoc-server.html)
documentation.  The response then is written to disk (or S3 bucket) as appropriate.

This should result in a relatively simple Go package and can work with io.Reader,
io.Writer types for maximum flexibility.  Combined with other services we hypothisis
is that we would see improved performance in rendering the website. The expectation
is that pandoc-server launches once, it is a single process (so no overhead on startup).
We have the existing overhead of the data source so that doesn't change. The documents 
are small for the most part so the network overhead between the client and pandoc-server
should be minimal (they are running on the same machine after all). The write of the
rendered document should be the same as our previous approach. The wind down of the process
from the exec is avoided.  We should be able to run conversions in parallel without
worrying about running out of process handles. More parallel writes should imply that
the overall time of the updates can be lowered.

Requirements
------------

- Go 1.19.2 or better
- Pandoc 2.19 or better
- A data source (e.g. file system with markdown documents)
- A place to write the output (e.g. a file system with render documents)




