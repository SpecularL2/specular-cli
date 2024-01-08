GO_CMD=go
GO_BUILD=$(GO_CMD) build
CMD_PATH=./cmd/spc/main.go
DIST=dist
DIST_LINUX=$(DIST)/linux
DIST_MACOS=$(DIST)/macos
BINARY_NAME=spc
REPO=$(shell .github/scripts/reponame.sh)
DOCKER_BUILD_ARGS=($BUILD_ARGS)
GIT_TAG=$(TAG)

.PHONY: help # Help - list of targets with descriptions
help:
	@echo ''
	@echo 'Usage:'
	@echo ' ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/ \1\t\2/' | expand -t20

.PHONY: lint-test # Run lint-tests
lint-test:
	go mod tidy
	goimports -local github.com/SpecularL2/specular-cli -w .
	go fmt ./...
	@golangci-lint -v run ./...
	@test -z "$$(golangci-lint run ./...)"

.PHONY: lint # Run lint-tests (alias)
lint: lint-test

.PHONY: build-linux # Build Linux binary
build-linux:
	mkdir -p $(DIST_LINUX)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o ./$(DIST_LINUX)/$(BINARY_NAME) -v $(CMD_PATH)

.PHONY: build-macos # Build macOS binary
build-macos:
	mkdir -p $(DIST_MACOS)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o ./$(DIST_MACOS)/$(BINARY_NAME) -v $(CMD_PATH)

.PHONY: build # Build Linux binary (alias)
build: build-linux

.PHONY: wire-generate # Generate Wire bindings
wire-generate:
	cd internal/service/di ;\
	wire

.PHONY: test # Run tests
test:
	go test --short ./...
