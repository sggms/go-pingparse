all: build test

build:
	go install -v github.com/sggms/go-pingparse/pinger github.com/sggms/go-pingparse/pinger/parser

test:
	go test -v github.com/sggms/go-pingparse/pinger/parser
	go test -v github.com/sggms/go-pingparse/pinger

.PHONY: all build test
