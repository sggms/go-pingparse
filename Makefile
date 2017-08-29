all: build test

build:
	go install -v github.com/sggms/go-pingparse/pinger github.com/sggms/go-pingparse/pinger/parser

test:
	go test -v github.com/sggms/go-pingparse/pinger github.com/sggms/go-pingparse/pinger/parser

.PHONY: all build test
