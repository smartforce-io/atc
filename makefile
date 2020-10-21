MAINPACKAGE=main
EXENAME=atcapp
BUILDPATH=$(CURDIR)
export GOPATH=$(CURDIR)

.PHONY: all clean build test

default : all

makedir :
	@if [ ! -d $(BUILDPATH)/bin ] ; then mkdir -p $(BUILDPATH)/bin ; fi
	@if [ ! -d $(BUILDPATH)/pkg ] ; then mkdir -p $(BUILDPATH)/pkg ; fi

build :
	@echo "building...."
	@go build -o $(BUILDPATH)/bin/$(EXENAME) $(MAINPACKAGE)

get :
	@echo "download 3rd party packages...."
	@go get github.com/gorilla/mux github.com/google/go-github/github

all : makedir get build

clean :
	@echo "cleaning...."
	@rm -rf $(BUILDPATH)/bin/$(EXENAME)
	@rm -rf $(BUILDPATH)/pkg
	@rm -rf $(BUILDPATH)/bin

test :
	go test -race -coverprofile=coverage.txt -covermode=atomic ./src/apiserver

cover :
	go test ./src/apiserver -v -coverprofile=apiserver.out
	go tool cover -html=apiserver.out
