.PHONY: build test run-local run-docker

RFC_3339 := "+%Y-%m-%dT%H:%M:%SZ"
DATE := $(shell date -u $(RFC_3339))
COMMIT := $(shell git rev-list -1 HEAD)
VERSION := $(shell git describe)

BUILDINFO_PATH := main

BUILDINFO_VERSION := -X $(BUILDINFO_PATH).Version=$(VERSION)
BUILDINFO_DATE := -X $(BUILDINFO_PATH).Date=$(DATE)
BUILDINFO_COMMIT := -X $(BUILDINFO_PATH).Commit=$(COMMIT)

BUILDINFO?=$(BUILDINFO_VERSION) $(BUILDINFO_DATE) $(BUILDINFO_COMMIT)

BUILD_OPTS?="-ldflags=$(BUILDINFO)" $(RACE_FLAG)

test:
	go test ./pkg/...

build:
	go build ${BUILD_OPTS} -o booker ./cmd/booker/booker.go

run-local:
	DB=memory APP_PORT=8080 LOG_LEVEL=debug LOG_PRETTY=1 ./booker

run-docker:
	docker-compose down && docker-compose up -d --build