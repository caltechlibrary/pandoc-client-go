/*
Copyright (c) 2022, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package pandoc_client

import (
	"bytes"
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
	if len(src) == 0 {
		t.Errorf("Expected content returned from cfg.Convert(), got none")
		t.FailNow()
	}
	t.Errorf("FIXME: Need to make sure I am getting valid HTML ... ->\n%s\n", src)
}
