SHELL:=/bin/bash
LINT_FILE_TAG=v0.2.0
LINT_FILE_URL=https://raw.githubusercontent.com/tuihub/librarian/$(LINT_FILE_TAG)/.golangci.yml

GOLANG_CROSS_VERSION ?= v1.20.7
PACKAGE_NAME := github.com/tuihub/tuihub-steam
VERSION=$(shell git describe --tags --always)
PROTO_VERSION=$(shell go list -m -f '{{.Version}}' github.com/tuihub/protos)

.PHONY: init
# init env
init:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

.PHONY: lint
# lint files
lint:
	golangci-lint run --fix -c <(curl -sSL $(LINT_FILE_URL))
	golangci-lint run -c <(curl -sSL $(LINT_FILE_URL)) # re-run to make sure fixes are valid, useful in some condition

.PHONY: release-dry-run
# build server in release mode, for manual test
release-dry-run:
	@docker run \
		--rm \
		-e PROTO_VERSION=$(PROTO_VERSION) \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip-validate --skip-publish

.PHONY: release
# build server in release mode, for CI, do not run manually
release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		-e PROTO_VERSION=$(PROTO_VERSION) \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
