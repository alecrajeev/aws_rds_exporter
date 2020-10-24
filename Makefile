IMAGE 		?= alecrajeev/aws_rds_exporter
VERSION 	= $(shell cat VERSION)
GO_PACKAGES = $(shell go list ./... | grep -v /vendor/)

all: format build test

test:
	@echo ">> running tests"
	@go test $(GO_PACKAGES)

format:
	@echo ">> formatting code"
	@go fmt $(GO_PACKAGES)

build:
	@echo ">> building binaries"
	@go build
	@echo ">> building docker image"
	@docker build \
		-t $(IMAGE):$(VERSION) \
		.
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest	

version:
	echo $(DOCKER_IMAGE_TAG)

.PHONY: all format build test docker
