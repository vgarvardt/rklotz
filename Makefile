NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

NAME=rklotz

# Build configuration
VERSION ?= "0.0.0-dev-$(shell git rev-parse --short HEAD)"
BUILD_DIR ?= $(CURDIR)
GO_LINKER_FLAGS=-ldflags "-s -w" -ldflags "-X main.version=$(VERSION)"

.PHONY: all
all: test build

.PHONY: build
build:
	@echo "$(OK_COLOR)==> Building (v${VERSION}) ... $(NO_COLOR)"
	@CGO_ENABLED=0 go build $(GO_LINKER_FLAGS) -o "$(BUILD_DIR)/${NAME}"

.PHONY: test
test:
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@CGO_ENABLED=0 go test -cover -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: lint
lint:
	golangci-lint run --config=./.github/linters/.golangci.yml --fix
