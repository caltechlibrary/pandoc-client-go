package pandoc_client

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	mdText := []byte(`---
title: "Hello World"
author: "jane.doe@example.org (Jane Doe)"
pubDate: 2022-11-04
---

Hello World
===========

By Jane Doe

Hi there Universe!

`)
	cfgText := []byte(`{
	"from": "markdown",
	"to": "html5",
	"standalone": true
}`)
	if _, err := os.Stat("testout"); os.IsNotExist(err) {
		os.MkdirAll("testout", 0775)
	}
	if err := os.WriteFile("testout/hello-world.json", cfgText, 0664); err != nil {
		t.Error(err)
		t.FailNow()
	}
	cfg, err := Load("testout/hello-world.json")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// Use verbose logging for tests.
	cfg.Verbose = true
	src, err := cfg.Convert(bytes.NewReader(mdText), "text/plain")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Fprintf(os.Stdout, "DEBUG src: %s\n", src)
}
