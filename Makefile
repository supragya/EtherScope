GO=go
GOBUILD=$(GO) build
BINDIR=build
BINCLI=bgidx
INSTALLLOC=/usr/local/bin/$(BINCLI)
RELEASE=$(shell git describe --tags --abbrev=0)
BUILDCOMMIT=$(shell git rev-parse HEAD)
BUILDLINE=$(shell git rev-parse --abbrev-ref HEAD)
CURRENTTIME=$(shell date -u '+%d-%m-%Y_%H-%M-%S')@UTC
BUILDER=$(shell uname -n)
VER='unversioned'

build:
	$(GOBUILD) -ldflags="\
	-X github.com/Blockpour/Blockpour-Geth-Indexer/version.ApplicationVersion=$(VER) \
	-X github.com/Blockpour/Blockpour-Geth-Indexer/version.buildCommit=$(BUILDLINE)@$(BUILDCOMMIT) \
	-X github.com/Blockpour/Blockpour-Geth-Indexer/version.buildTime=$(CURRENTTIME) \
	-X github.com/Blockpour/Blockpour-Geth-Indexer/version.builder=$(BUILDER) \
	-linkmode=external" \
	-o $(BINDIR)/$(BINCLI)
clean:
	rm -rf $(BINDIR)/*

install:
	cp $(BINDIR)/$(BINCLI) $(INSTALLLOC)

uninstall:
	rm $(INSTALLLOC)