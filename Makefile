GO = go
GIT = git
GOLANGCI-LINT = golangci-lint
GORELEASER = goreleaser
BUF = buf
YARN = yarn

BIN = /usr/local/bin

SEMVER ?= 0.1.0

build:
	@$(GORELEASER) release --snapshot --rm-dist

generate test:
	@$(GO) $@ ./...

fmt:
	@$(GO) $@ ./...
	@$(BUF) format -w
	@$(YARN) $@

download vendor verify:
	@$(GO) mod $@

lint:
	@$(GOLANGCI-LINT) run --fix

release:
	@$(GIT) tag v$(SEMVER)
	@$(GIT) push origin --tags

pb:
	@$(BUF) generate

gen: generate
dl: download
ven: vendor
ver: verify
format: fmt
proto protos buf: pb

.PHONY: buf build dl download fmt format gen generate lint pb proto protos release test ven vendor ver verify
