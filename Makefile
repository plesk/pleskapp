# Copyright 1999-2021. Plesk International GmbH.

OUTFILE=plesk
REVISON:=$(shell git rev-parse --short HEAD)
VERSION:=$(shell cat VERSION)
BUILD_TIME=$(shell date +'%Y-%m-%d_%T')
LDFLAGS=-X main.revision=$(REVISON) -X main.buildTime=$(BUILD_TIME) -X main.version=$(VERSION)
RELEASE_LDFLAGS=$(LDFLAGS) -s -w

.PHONY: all build clean test

build: test
	go build -ldflags "$(LDFLAGS)"

release: test
	GOOS=linux go build -ldflags "$(RELEASE_LDFLAGS)" -o ./build/linux/$(OUTFILE)
	tar czf ./build/$(OUTFILE)-v$(VERSION)-linux.tgz build/linux/$(OUTFILE)
	GOOS=darwin go build -ldflags "$(RELEASE_LDFLAGS)" -o ./build/mac/$(OUTFILE)
	tar czf ./build/$(OUTFILE)-v$(VERSION)-mac.tgz build/mac/$(OUTFILE)
	GOOS=windows go build -ldflags "$(RELEASE_LDFLAGS)" -o ./build/win/$(OUTFILE).exe
	tar czf ./build/$(OUTFILE)-v$(VERSION)-win.tgz build/win/$(OUTFILE).exe

run:
	go run main.go

clean:
	$(RM) $(OUTFILE) ./build/*/$(OUTFILE)

test:
	go test -v -cover ./...

all: clean release build