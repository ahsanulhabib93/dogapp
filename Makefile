.PHONY: default setup lint build test test.cover test.report

SHELL := /bin/bash # Use bash syntax

APP_EXECUTABLE="out/ss2"

GOPATH=$(shell go env GOPATH)
GONOSUMDB="github.com/voonik/*,github.com/shopuptech/*"
ENV=test
FILES_TO_EXCLUDE="'|$(shell yq e '.files' coverignore.yaml | tr '\n' '|' | tr -d '-' | tr -d [:blank:])'"
DIRS_TO_EXCLUDE="'|$(shell yq e '.dirs' coverignore.yaml | tr '\n' '|' | tr -d '-' | tr -d [:blank:])'"
PKGLIST="$(shell go list ./... | grep -v -E $(DIRS_TO_EXCLUDE) | tr '\n' ',')"


export GOPATH
export GONOSUMDB
export ENV
export FILES_TO_EXCLUDE
export DIRS_TO_EXCLUDE
export PKGLIST

default: setup test build

setup:
	go install github.com/axw/gocov/gocov@latest
	go install github.com/t-yuki/gocover-cobertura@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin

lint:
	golangci-lint run

build:
	mkdir -p out/
	GO111MODULE=on go build -o $(APP_EXECUTABLE) ./cmd/server/

# Ensure the test environment is correctly set up and dependencies are installed
test:
# Display environment variables to help debug potential issues
	echo "GOPATH: $(GOPATH)"
	echo "PKGLIST: $(PKGLIST)"
	echo "FILES_TO_EXCLUDE: $(FILES_TO_EXCLUDE)"
	echo "DIRS_TO_EXCLUDE: $(DIRS_TO_EXCLUDE)"

# Run tests with detailed output to help identify any failures
	mkdir -p ./coverage/
	GO111MODULE=on go clean -testcache ./internal/app/test/... && go test -race ./internal/app/test/... -v -count=1 -p 1 -covermode=atomic -coverprofile=./coverage/coverage.out.temp -coverpkg=$(PKGLIST)

test.cover: test
	cat ./coverage/coverage.out.temp | grep -v -E $(FILES_TO_EXCLUDE) > ./coverage/coverage.out
	GO111MODULE=on gocov convert ./coverage/coverage.out | gocov report 2>&1 | tee ./coverage/coverage.txt

test.report: test.cover
	GO111MODULE=on go tool cover -html ./coverage/coverage.out -o ./coverage/coverage.html
	GO111MODULE=on gocover-cobertura < ./coverage/coverage.out > ./coverage/coverage.xml
