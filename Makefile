BINARY=k8guard-action

VERSION=`git fetch;git describe --tags > /dev/null 2>&1`
BUILD=`date +%FT%T%z`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

all: deps build test

deps:
	glide install

glide-update:
	glide cc
	glide update

build-docker:
	docker build -t local/k8guard-action .

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}

mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY}

dev-setup:
	go get golang.org/x/tools/cmd/goimports

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	go clean

sclean: clean
	rm glide.lock

.PHONY: build
