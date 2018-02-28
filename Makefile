# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOLINT=golint
GOX=gox
DEPCMD=dep
DEPENSURE=$(DEPCMD) ensure

all: test clean build
build:
		for dir in $(find . -type d -maxdepth 1 | grep -v .git); do pushd ${dir}; make build; popd; done
install:
		for dir in $(find . -type d -maxdepth 1 | grep -v .git); do pushd ${dir}; make install; popd; done
test:
		for dir in $(find . -type d -maxdepth 1 | grep -v .git); do pushd ${dir}; make test; popd; done
clean:
		for dir in $(find . -type d -maxdepth 1 | grep -v .git); do pushd ${dir}; make clean; popd; done
#run:
#		$(GOBUILD) -o $(BINARY_NAME) -v ./...
#		./$(BINARY_NAME)
deps:
		for dir in $(find . -type d -maxdepth 1 | grep -v .git); do pushd ${dir}; make deps; popd; done

## Cross compilation
#build-linux:
#		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
#docker-build:
#		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
