OUTFILE=plesk
REVISON:=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date +'%Y-%m-%d_%T')
LDFLAGS=-X main.revision=$(REVISON) -X main.buildTime=$(BUILD_TIME)
RELEASE_LDFLAGS=$(LDFLAGS) -s -w

.PHONY: all build clean test

build: test
	go build -ldflags "$(LDFLAGS)"

release: test
	GOOS=linux go build -ldflags "$(RELEASE_LDFLAGS)" -o ./build/linux/$(OUTFILE)
	GOOS=darwin go build -ldflags "$(RELEASE_LDFLAGS)" -o ./build/mac/$(OUTFILE)
	GOOS=windows go build -ldflags "$(RELEASE_LDFLAGS)" -o ./build/win/$(OUTFILE)

run:
	go run main.go

clean:
	$(RM) $(OUTFILE) ./build/*/$(OUTFILE)

test:
	go test -v -cover ./...

all: clean release build