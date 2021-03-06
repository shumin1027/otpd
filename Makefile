PWD=$(shell pwd)
DIST=$(shell pwd)/bin
DATE=$(shell date --iso-8601=seconds)

GIT_SHA=$(shell git rev-parse HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GIT_CLOSEST_TAG=$(shell git describe --abbrev=0 --tags)

PKGS=$(shell go list ./... | grep -v /vendor/)

PROJECT="otpd"

BUILD_INFO_IMPORT_PATH=github.com/shumin1027/otpd/app
BUILD_INFO='-X $(BUILD_INFO_IMPORT_PATH).BuildTime=$(DATE) -X $(BUILD_INFO_IMPORT_PATH).GitCommit=$(GIT_SHA) -X $(BUILD_INFO_IMPORT_PATH).GitBranch=$(GIT_BRANCH) -X $(BUILD_INFO_IMPORT_PATH).GitTag=$(GIT_CLOSEST_TAG)'

GOPATH := $(HOME)/go
PATH := $(GOPATH)/bin/:$(PATH)

.PHONY: dev-tools
dev-tools:
ifneq (,$(wildcard $(GOPATH)/bin/swag))
	@swag --version
else
	@go install github.com/swaggo/swag/cmd/swag@latest
endif

.PHONY: build
build:clean fmt vet doc
	@echo ">> building code"
	go build -mod=vendor -tags=jsoniter -ldflags='-w -s -linkmode=external' -ldflags=$(BUILD_INFO) -o $(DIST)/otpd $(PWD)/main.go
	strip $(DIST)/otpd

.PHONY: docker
docker:
	@echo ">> building docker images"
	 docker build -t otpd .

.PHONY: fmt
fmt:
	@echo ">> fmt code"
	go fmt $(PKGS)

.PHONY: vet
vet:
	@echo ">> vetting code"
	go vet $(PKGS)

.PHONY: clean
clean:
	@echo ">> clean build"
	go clean -i -x 
	rm -rf $(BUILDDIR)

.PHONY: clean-cache
clean-cache:
	@echo ">> clean build cache"
	go clean -cache -testcache

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: update
update:
	go get -u
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: doc
doc:
	@echo ">>Generating API DOC"
	swag init http/web.go
	@echo "Done."