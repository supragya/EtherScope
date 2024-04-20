GO=go
GOVER=$(shell go version)
GOBUILD=$(GO) build
BINDIR=build
BINCLI=escope
INSTALLLOC=/usr/local/bin/$(BINCLI)
RELEASE=$(shell git describe --tags --abbrev=0)
BUILDCOMMIT=$(shell git rev-parse HEAD | cut -c 1-7)
BUILDLINE=$(shell git rev-parse --abbrev-ref HEAD)
CURRENTTIME=$(shell date -u '+%d-%m-%Y %H:%M:%S')
CURRENTTAG=$(shell git tag -l --contains HEAD)
BUILDER=$(shell uname -n)

build:
	$(GOBUILD) -ldflags="\
	-X 'github.com/supragya/EtherScope/version.buildCommit=$(BUILDLINE)@$(BUILDCOMMIT)' \
	-X 'github.com/supragya/EtherScope/version.gittag=$(CURRENTTAG)' \
	-X 'github.com/supragya/EtherScope/version.buildTime=$(CURRENTTIME)' \
	-X 'github.com/supragya/EtherScope/version.builder=$(BUILDER)' \
	-X 'github.com/supragya/EtherScope/version.gover=$(GOVER)' \
	-linkmode=external" \
	-o $(BINDIR)/$(BINCLI)
clean:
	rm -rf $(BINDIR)/*

install:
	cp $(BINDIR)/$(BINCLI) $(INSTALLLOC)

uninstall:
	rm $(INSTALLLOC)

docker:
	docker build \
	-t geth-indexer \
	--build-arg buildCommit=$(BUILDLINE)@$(BUILDCOMMIT) \
	--build-arg gittag=$(CURRENTTAG) \
	--build-arg 'buildTime=$(CURRENTTIME)' \
	--build-arg 'builder=$(BUILDER)' \
	--build-arg 'gover=$(GOVER)' .

dockerbuildx:
	docker buildx create --use && \
	docker buildx build --platform=linux/amd64,linux/arm64 \
	--push \
	-t geth-indexer \
	--build-arg buildCommit=$(BUILDLINE)@$(BUILDCOMMIT) \
	--build-arg gittag=$(CURRENTTAG) \
	--build-arg 'buildTime=$(CURRENTTIME)' \
	--build-arg 'builder=$(BUILDER)' \
	--build-arg 'gover=$(GOVER)' \
	-f ./Dockerfile.multiarch .
