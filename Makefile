.DEFAULT_GOAL := help
.PHONY: help

BINARY=k8guard-action
VERSION=`git fetch;git describe --tags > /dev/null 2>&1`
BUILD=`date +%FT%T%z`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

all: deps build test ## install dependecies, build for linux

create-pre-commit-hooks: ## creates pre-commit hooks
	chmod +x $(CURDIR)/hooks/pre-commit
	ln -s $(CURDIR)/hooks/pre-commit .git/hooks/pre-commit || true

deps: ## install dependecies
	glide install

glide-update:
	glide cc
	glide update

build-docker: ## builds docker local/k8guard-action
	docker build -t local/k8guard-action .

build:  ## builds binary for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}

mac: ## builds binary for mac
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}

dev-setup:
	go get golang.org/x/tools/cmd/goimports

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	go clean

sclean: clean ## super clean, cleans binary and glide lock
	rm glide.lock

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
