NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

NAME=rklotz
REPO=github.com/vgarvardt/${NAME}

# Build configuration
VERSION ?= "$(shell cat ./VERSION)"
BUILD_DIR ?= $(CURDIR)/build
GO_LINKER_FLAGS=-ldflags "-s -w" -ldflags "-X ${REPO}/cmd.version=$(VERSION)"

.PHONY: all clean deps build

all: clean deps build

deps:
	@echo "$(OK_COLOR)==> Installing dev dependencies$(NO_COLOR)"
	@go get -u github.com/go-playground/overalls
	@go get -u golang.org/x/lint/golint

build:
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@CGO_ENABLED=0 go build -mod vendor $(GO_LINKER_FLAGS) -o "$(BUILD_DIR)/${NAME}"
	@GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -mod vendor $(GO_LINKER_FLAGS) -o "$(BUILD_DIR)/${NAME}.linux.amd64"
	@docker build --no-cache --pull -t vgarvardt/rklotz:`cat ./VERSION` .

push:
	@docker push vgarvardt/rklotz:`cat ./VERSION`

test: lint format vet
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@CGO_ENABLED=0 go test -cover ./... -coverprofile=coverage.txt -covermode=atomic

lint:
	@echo "$(OK_COLOR)==> Checking code style with 'golint' tool$(NO_COLOR)"
	@go list ./... | xargs -n 1 golint -set_exit_status

format:
	@echo "$(OK_COLOR)==> Checking code formating with 'gofmt' tool$(NO_COLOR)"
	@gofmt -l -s cmd pkg | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

vet:
	@echo "$(OK_COLOR)==> Checking code correctness with 'go vet' tool$(NO_COLOR)"
	@go vet ./...

clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	@if [ -d ${BUILD_DIR} ] ; then rm -rf ${BUILD_DIR}/* ; fi
