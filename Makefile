GO_CMD=go
GO_BUILD=$(GO_CMD) build
CMD_PATH=./cmd/spc/main.go
DIST=dist
DIST_LINUX=$(DIST)/linux
DIST_MACOS=$(DIST)/macos
BINARY_NAME=spc

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
	@go install golang.org/x/tools/cmd/goimports@latest
	goimports -local github.com/SpecularL2/specular-cli -w .
	go fmt ./...
	@golangci-lint -v run ./...
	@test -z "$$(golangci-lint run ./...)"

.PHONY: lint # Run lint-tests (alias)
lint: lint-test

.PHONY: build-linux # Build Linux binary
build-linux:
	mkdir -p $(DIST_LINUX)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o ./$(DIST_LINUX)/$(BINARY_NAME)-linux-amd64 -v $(CMD_PATH)

.PHONY: build-macos # Build macOS binary
build-macos:
	mkdir -p $(DIST_MACOS)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO_BUILD) -o ./$(DIST_MACOS)/$(BINARY_NAME)-macos-arm64 -v $(CMD_PATH)

.PHONY: build-docker # Build Docker image
build-docker:
	docker build . -t ghcr.io/specularl2/specular-cli:latest -t spc

.PHONY: docker-push # Push Docker image to registry
docker-push:
	docker push ghcr.io/specularl2/specular-cli:latest

.PHONY: build # Build Linux binary (alias)
build:
	make build-linux
	make build-macos
	make build-docker

.PHONY: wire-generate # Generate Wire bindings
wire-generate:
	cd internal/service/di ;\
	wire

.PHONY: test # Run tests
test:
	go test ./...
