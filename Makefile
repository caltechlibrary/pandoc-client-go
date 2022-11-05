#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = pandoc_client

VERSION = $(shell grep '"version":' codemeta.json | cut -d\"  -f 4)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PANDOC = $(shell which pandoc)

PROGRAMS = $(shell ls -1 cmd)

MANPAGES = md2html.1

PACKAGE = $(shell ls -1 *.go | grep -v 'version.go')

SUBPACKAGES = $(shell ls -1 */*.go)

#PREFIX = /usr/local/bin
PREFIX = $(HOME)

ifneq ($(prefix),)
        PREFIX = $(prefix)
endif

OS = $(shell uname)

EXT = 
ifeq ($(OS), Windows)
	EXT = .exe
endif

DIST_FOLDERS = bin/*

build: version.go $(PROGRAMS) CITATION.cff about.md

version.go: .FORCE
	@echo "package $(PROJECT)" >version.go
	@echo '' >>version.go
	@echo 'const (' >>version.go
	@echo '// Version of package' >>version.go
	@echo '    Version = "$(VERSION)"' >>version.go
	@echo '' >>version.go
	@echo '    License = `' >>version.go
	@cat LICENSE >>version.go
	@echo '`' >>version.go
	@echo ')' >>version.go
	@echo '' >>version.go
	@git add version.go

CITATION.cff: codemeta.json
	@cat codemeta.json | sed -E   's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	@if [ -f $(PANDOC) ]; then echo "" | $(PANDOC) --metadata title="Cite $(PROJECT)" --metadata-file=_codemeta.json --template=codemeta-cff.tmpl >CITATION.cff; fi

about.md: codemeta.json $(PROGRAMS)
	@cat codemeta.json | sed -E 's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	@if [ -f $(PANDOC) ]; then echo "" | pandoc --metadata-file=_codemeta.json --template codemeta-md.tmpl >about.md 2>/dev/null; fi
	@if [ -f _codemeta.json ]; then rm _codemeta.json; fi

$(PROGRAMS): cmd/*/*.go $(PACKAGE)
	@mkdir -p bin
	go build -o bin/$@$(EXT) cmd/$@/*.go

$(MANPAGES): .FORCE
	mkdir -p man/man1
	$(PANDOC) $@.md --from markdown --to man -s >man/man1/$@

man: $(MANPAGES)

# NOTE: on macOS you must use "mv" instead of "cp" to avoid problems
install: build man
	@if [ ! -d $(PREFIX)/bin ]; then mkdir -p $(PREFIX)/bin; fi
	@echo "Installing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f ./bin/$$FNAME ]; then mv -v ./bin/$$FNAME $(PREFIX)/bin/$$FNAME; fi; done
	@echo ""
	@echo "Make sure $(PREFIX)/bin is in your PATH"
	@echo ""
	@echo "Installing manpages in $(PREFIX)/bin"
	@if [ ! -d $(PREFIX)/man/man1 ]; then mkdir -p $(PREFIX)/man/man1; fi
	@for FNAME in $(MANPAGES); do if [ -f ./man/man1/$$FNAME ]; then cp -v ./man/man1/$$FNAME $(PREFIX)/man/man1/$$FNAME; fi; done
	@echo ""
	@echo "Make sure $(PREFIX)/man is in your MANPATH"
	@echo ""

uninstall: .FORCE
	@echo "Removing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f $(PREFIX)/bin/$$FNAME ]; then rm -v $(PREFIX)/bin/$$FNAME; fi; done
	@echo "Removing manpages in $(PREFIX)/man"
	@for FNAME in $(MANPAGES); do if [ -f $(PREFIX)/man/man1/$$FNAME ]; then rm -v $(PREFIX)/man/man1/$$FNAME; fi; done


website: page.tmpl README.md INSTALL.md LICENSE
	make -f website.mak

check: .FORCE
	go vet *.go

test: clean build
	go test

cleanweb:
	make -f website.mak clean

clean: 
	@if [ -d bin ]; then rm -fR bin; fi
	@if [ -d dist ]; then rm -fR dist; fi
	@if [ -d testout ]; then rm -fR testout; fi

dist/linux-amd64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env  GOOS=linux GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-linux-amd64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/macos-amd64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macos-amd64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/macos-arm64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=arm64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macos-arm64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/windows-amd64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=amd64 go build -o dist/bin/$$FNAME.exe cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-windows-amd64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/windows-arm64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=arm64 go build -o dist/bin/$$FNAME.exe cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-windows-arm64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/raspbian-arm7:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-rasperry-pi-os-arm7.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

distribute_docs:
	if [ -d dist ]; then rm -fR dist; fi
	mkdir -p dist
	cp -v codemeta.json dist/
	cp -v CITATION.cff dist/
	cp -v README.md dist/
	cp -v about.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -vfR man dist/

update_version:
	$(EDITOR) codemeta.json
	codemeta2cff

release: clean build CITATION.cff about.md man distribute_docs dist/linux-amd64 dist/windows-amd64 dist/windows-arm64 dist/macos-amd64 dist/macos-arm64 dist/raspbian-arm7

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

#publish:
#	bash mk-website.bash
#	bash publish.bash

.FORCE:
