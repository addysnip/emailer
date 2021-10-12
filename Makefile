VERSION="dev"
ifdef DRONE_TAG
  VERSION=$(DRONE_TAG)
endif
LDFLAGS=-w -X 'addysnip.dev/emailer/pkg/version.BuildTime=$(shell TZ='UTC' date)' -X 'addysnip.dev/emailer/pkg/version.GitCommit=$(shell git rev-parse --short HEAD)' -X 'addysnip.dev/emailer/pkg/version.Version=$(VERSION)' -X 'addysnip.dev/emailer/pkg/version.GoVersion=$(shell go version | awk '{print $$3}')'
COMPILER=go build -ldflags="$(LDFLAGS)"

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(COMPILER) -o build/app

clean:
	rm -rf build