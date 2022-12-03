# Copyright 1999-2022. Plesk International GmbH.

OUTFILE=plesk
COMMIT:=$(shell git rev-parse --short HEAD)
TAG:=$(shell git describe --abbrev=0 --tags)
VERSION:=$(TAG:v%=%)
DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-X main.commit=$(COMMIT) -X main.date=$(DATE) -X main.version=$(VERSION)

.PHONY: all build clean test

build: test
	go build -ldflags "$(LDFLAGS)"

release: test
	goreleaser release

run:
	go run main.go

clean:
	$(RM) $(OUTFILE)

test:
	go test -v -cover ./...

all: clean release build
