GO = go
GIT = git
GOLANGCI-LINT = golangci-lint
GORELEASER = goreleaser
YARN = yarn
INSTALL = sudo install
DOCKER = docker
DOCKER-COMPOSE = $(DOCKER) compose

BIN = /usr/local/bin

GOOS = $(shell $(GO) env GOOS)
GOARCH = $(shell $(GO) env GOARCH)

SEMVER ?= 0.1.0

up:
	@$(DOCKER-COMPOSE) $@ --build

generate test:
	@$(GO) $@ ./...

fmt:
	@$(GO) $@ ./...
	@$(YARN) $@

download tidy vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

release:
	@$(YARN) version --new-version=$(SEMVER)
	@$(GIT) tag v$(SEMVER)
	@$(GIT) push origin --tags

gen: generate
dl: download
ven: vendor
ver: verify
format: fmt

.PHONY: build dl download fmt format gen generate lint release test ven vendor ver verify
