GOLINTFLAGS ?=
GOBIN        = $(shell go env GOPATH)/bin
GOLINT       = PATH=$(GOBIN):$(PATH) golangci-lint --color=always $(GOLINTFLAGS)

VERSION  = $(shell date '+%Y%m%d%H%M%S')
BRANCH   = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse --short HEAD)

GOPREFIX  = github.com/go-park-mail-ru/2020_2_softree/build
LDFLAGS   = \
			-X $(GOPREFIX).Version=$(VERSION) \
			-X $(GOPREFIX).Branch=$(BRANCH) \
			-X $(GOPREFIX).Revision=$(REVISION)

GOTAGS ?=
ifeq ($(DEBUG), 1)
	GOTAGS = debug
endif

ifneq ($(GOTAGS),)
	GOTAGS := -tags $(GOTAGS)
endif

DOCKERFILE = Dockerfile
ifeq ($(DEBUG), 1)
	DOCKERFILE = Dockerfile.debug
endif


.PHONY: all
all: build

.PHONY: build
build:
	go build $(GOTAGS) \
		-ldflags '$(LDFLAGS)' \
		-o bin/mc \
		./src

.PHONY: test
test:
	sh test.sh

.PHONY: fmt
fmt:
	gofmt -s -w src

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor: tidy
	go mod vendor

.PHONY: docker-build
docker-build: vendor
	docker build -f $(DOCKERFILE) -t mc_api .

.PHONY: run
run: deps
	go run ./src/main.go

.PHONY: deps
deps:
	go mod download

.PHONY: lint
lint:
	@$(GOLINT) run ./...
