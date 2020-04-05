TARGET   = claim-parser
PACKAGE  = github.com/alexey-zayats/claim-parser/cmd/claim-parser
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/VERSION 2> /dev/null || echo v0.0.0)
BASE     = $(CURDIR)
BIN      = $(CURDIR)/bin
PKGS     = $(or $(PKG),$(shell cd $(BASE) && $(GO) list ./...))

PROTOS  = $(wildcard internal/proto/*.proto)
PROTOC  = protoc -I$(GOROOT)/src -I. -I/usr/local/include -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

GO      = go
GODOC   = godoc
GOFMT   = gofmt
GLIDE   = glide
TIMEOUT = 30
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: fmt proto lint binary

binary: | ; $(info $(M) building executable…) @ ## Build program binary
		$Q cd $(BASE) && $(GO) build \
		-tags release \
		-ldflags "-X $(PACKAGE)/cmd.Version='$(VERSION)' -X $(PACKAGE)/cmd.BuildDate='$(DATE)'" \
		-o bin/$(TARGET) ./cmd/$(TARGET)/main.go

MIGRATE = $(GOROOT)/bin/migrate
$(GOROOT)/bin/migrate: | ; $(info $(M) building migrate…)
	$Q go get -u -d github.com/golang-migrate/migrate/cli github.com/lib/pq
	$Q go build -tags 'postgres' -o $(GOROOT)/bin/migrate github.com/golang-migrate/migrate/cli

GOLINT = $(GOROOT)/bin/golint
$(GOROOT)/bin/golint: | ; $(info $(M) building golint…)
	$Q go get -u -d golang.org/x/lint/golint
	$Q go build -o  $(GOROOT)/bin/golint golang.org/x/lint/golint

.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	done ; exit $$ret

.PHONY: proto
proto: ; $(info $(M) running protoc…) @ ## Run protoc for all grpc proto
	@ret=0 && for d in $(PROTOS); do \
		$(PROTOC) --go_out=plugins=grpc:. ./$$d  || ret=$$? ; \
	done ; exit $$ret

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./...); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
		done ; exit $$ret

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf bin

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)

