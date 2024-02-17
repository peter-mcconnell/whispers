DOCKER_IMG?=pemcconnell/whispers
BIN_DIR?=./bin
TAG?=$(shell git rev-parse --short HEAD)

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
	CGO_ENABLED=1 go build -o $(BIN_DIR)/whispers ./cmd/whispers

.PHONY: docker-build
docker-build:
	docker build -t=$(DOCKER_IMG):$(TAG) -f Dockerfile .

.PHONY: docker-run
docker-run: docker-build
	@-docker rm -f whispers > /dev/null 2>&1
	docker run --privileged --pid=host -v /sys/kernel/:/sys/kernel/ --name whispers --rm -p 2222:22 -d $(DOCKER_IMG):$(TAG)

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMG):$(TAG)

.PHONY: docker-exec
docker-exec:
	docker exec -ti whispers bash
