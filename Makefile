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

all: test clean deps build
build:
		find . -maxdepth 1 -type d | grep -v .git | grep -v "\.$$" | while read dir; do make build -C $${dir}; done
install:
		find . -maxdepth 1 -type d | grep -v .git | grep -v "\.$$" | while read dir; do make install -C $${dir}; done
test:
		find . -maxdepth 1 -type d | grep -v .git | grep -v "\.$$" | while read dir; do make test -C $${dir}; done
clean:
		find . -maxdepth 1 -type d | grep -v .git | grep -v "\.$$" | while read dir; do make clean -C $${dir}; done
#run:
#		$(GOBUILD) -o $(BINARY_NAME) -v ./...
#		./$(BINARY_NAME)
deps:
		find . -maxdepth 1 -type d | grep -v .git | grep -v "\.$$" | while read dir; do make deps -C $${dir}; done

## Cross compilation
#build-linux:
#		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
#docker-build:
#		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
