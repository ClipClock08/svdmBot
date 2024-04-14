 IMAGE_REGISTRY 		?= unknown
IMAGE_NAME 			?= s1-secret-rotation
FULL_IMAGE_NAME		:= $(IMAGE_REGISTRY)/$(IMAGE_NAME)
VERSION 			?= dev
GO_PACKAGES			:= $(shell go list ./...)
GO_FILES        	:= $(shell find . -type f -name '*.go')
DOCKER_RUNNER 		:= docker run $(FULL_IMAGE_NAME):$(VERSION)
BASE_IMAGE          ?= public.ecr.aws/lambda/provided
BASE_IMAGE_VERSION  ?= al2023
REPOSITORY_PATH 	:= lambda-s1-secret-rotation
GO_VERSION 			:= golang:1.21.8-bullseye
DOCKER_RUNNER 		:= docker run -v $(CURDIR):/go/src/$(REPOSITORY_PATH) -w /go/src/$(REPOSITORY_PATH) ${GO_VERSION}

.PHONY: build
build:
	@docker build -f ./build/Dockerfile -t gcr.io/svdm-test/svdmtgbot .

.PHONY: test
test:
	${DOCKER_RUNNER} sh -c "go test -v --race -cover $(GO_PACKAGES)"

.PHONY: push
push:
	docker push $(FULL_IMAGE_NAME):$(VERSION)

.PHONY: clean
clean:
	@docker rmi -f $(shell docker images -q $(FULL_IMAGE_NAME)) || true

.PHONY: fmt
fmt:
	@go fmt $(GO_PACKAGES)

.PHONY: lint
lint:
	@! gofmt -l . | grep -v vendor/
	@golangci-lint run
	@gocritic check -enableAll ./...
	@gosec ./...

.PHONY: goimports
goimports:
	@docker run --rm -v $(shell pwd):/data cytopia/goimports -w "$(GO_FILES)"

.PHONY: all
all: fmt lint test build push
