BUILDDIR := build
CGO_ENABLED := 0
NAME := jsonecho
VERSION := $(shell git describe --always --tags)

.DEFAULT_TARGET := build
.PHONY: build

build:
	go build \
	  -o $(BUILDDIR)/$(NAME) \
	  -ldflags "-X main.versionNumber=$(VERSION)"
