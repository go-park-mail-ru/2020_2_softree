GOLINTFLAGS ?=
GOBIN        = $(shell go env GOPATH)/bin

VERSION  = $(shell date '+%Y%m%d%H%M%S')
BRANCH   = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse --short HEAD)

GOPREFIX  = github.com/go-park-mail-ru/2020_2_softree
LDFLAGS   = \
			-X $(GOPREFIX).Version=$(VERSION) \
			-X $(GOPREFIX).Branch=$(BRANCH) \
			-X $(GOPREFIX).Revision=$(REVISION)

GOTAGS ?=
ifneq ($(GOTAGS),)
	GOTAGS := -tags $(GOTAGS)
endif

.PHONY: all
all: build

.PHONY: build
build:
ifndef TARGET
	@echo 'build target is not defined'
else
	go build $(GOTAGS) \
		-ldflags '$(LDFLAGS)' \
		-o bin/$(TARGET) \
		./cmd/$(TARGET)
endif

.PHONY: test
test:
	sh test.sh

.PHONY: fmt
fmt:
	gofmt -s -w ./cmd

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor: tidy
	go mod vendor

.PHONY: docker-build
docker-build: vendor
ifndef TARGET
	@echo 'build target is not defined'
else
		docker build -f docker-images/Dockerfile.$(TARGET) -t $(TARGET) .
endif

.PHONY: run
run:
	sudo docker-compose up -d --build

.PHONY: deps
deps:
	go mod download

.PHONY: lint
lint:
	golangci-lint --color=always run ./... $(GOLINTFLAGS)
