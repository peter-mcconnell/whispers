DOCKER_IMG?=pemcconnell/whispers
BIN_DIR?=./bin
TAG?=$(shell git rev-parse --short HEAD)
BUILD_VCS?=true
KERNEL?=$(shell uname -r)
GOOS?=linux
GOARCH?=amd64

.PHONY: test
test:
	go test -v -race ./...

.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: whispers
whispers:
	mkdir -p $(BIN_DIR)
	go mod tidy
	go generate ./...
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -buildvcs=$(BUILD_VCS) -o $(BIN_DIR)/whispers ./cmd/whispers

.PHONY: docker-build
docker-build:
	docker build --build-arg GOARCH=$(GOARCH) --platform=$(GOOS)/$(GOARCH) -t=$(DOCKER_IMG):$(GOOS)-$(GOARCH)-$(TAG) -t=$(DOCKER_IMG):$(GOOS)-$(GOARCH)-latest -f Dockerfile .

# we use the "base" target in CI
.PHONY: docker-base
docker-base:
	docker build --build-arg GOARCH=$(GOARCH) --platform=$(GOOS)/$(GOARCH) -t=$(DOCKER_IMG)-base:$(TAG) -t=$(DOCKER_IMG)-base:latest -f Dockerfile --target base .

.PHONY: docker-run
docker-run: docker-build
	@-docker rm -f whispers > /dev/null 2>&1
	docker run --platform=$(GOOS)/$(GOARCH) --privileged --name whispers --rm -p 2222:22 -d $(DOCKER_IMG):$(TAG)

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMG):$(GOOS)-$(GOARCH)-$(TAG)
	docker push $(DOCKER_IMG):$(GOOS)-$(GOARCH)-latest

.PHONY: docker-exec
docker-exec:
	docker exec -ti whispers bash

.PHONY: vmlinux
vmlinux:
	if [ ! -f bpf/headers/vmlinux.h ]; then bpftool btf dump file /sys/kernel/btf/vmlinux format c > bpf/headers/vmlinux.h; fi
