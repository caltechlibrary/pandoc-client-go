package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library packages
	"github.com/caltechlibrary/pandoc_client"
)

func main() {
	appName := path.Base(os.Args[0])
	showHelp, showVersion, showLicense := false, false, false
	flag.BoolVar(&showHelp, "help", showHelp, "display help")
	flag.BoolVar(&showVersion, "version", showVersion, "display version")
	flag.BoolVar(&showLicense, "license", showLicense, "display license")
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
	cfg.From = "markdown"
	cfg.To = "html5"
	if err := cfg.Walk(args[1], ".md", ".html"); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
