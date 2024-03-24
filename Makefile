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
	@goreleaser --skip=publish --snapshot --clean

.PHONY: test
test:
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@CGO_ENABLED=0 go test -cover -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: lint
lint:
	@echo "$(OK_COLOR)==> Running code linters$(NO_COLOR)"
	@golangci-lint run --config=./.github/linters/.golangci.yml --fix

.PHONY: spell
spell:
	@echo "$(OK_COLOR)==> Running spell check linters$(NO_COLOR)"
	@docker run \
		--interactive --tty --rm \
		--volume "$(CURDIR):/workdir" \
		--workdir "/workdir" \
		python:3.12-slim bash -c "python -m pip install --upgrade pip && pip install 'codespell>=2.2.4' && codespell"
