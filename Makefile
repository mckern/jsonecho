BUILDDIR := build
NAME := jsonecho
VERSION := $(shell git describe --always --tags)

.DEFAULT_TARGET := build
.PHONY: build

build: export CGO_ENABLED := 0
build:
	go build \
	  -o $(BUILDDIR)/$(NAME) \
	  -ldflags "-X main.versionNumber=$(VERSION)"

clean:
	$(RM) $(BUILDDIR)/$(NAME)

cleanest: clean
	$(RM) -rv go.sum
