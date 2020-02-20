.PHONY: test run

TAG?=latest
NAME:=hls_go
DOCKER_REPOSITORY:=rodrigodev
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)
VERSION:=0.1

default: build

build: test
	go build -o ./hls_go ./

build-container:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

push-container:
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):latest
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)
	docker push $(DOCKER_IMAGE_NAME):latest
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) quay.io/$(DOCKER_IMAGE_NAME):$(VERSION)
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) quay.io/$(DOCKER_IMAGE_NAME):latest
	docker push quay.io/$(DOCKER_IMAGE_NAME):$(VERSION)
	docker push quay.io/$(DOCKER_IMAGE_NAME):latest

test: lint
	go fmt ./...
	go test -vet all ./...

lint: get-linter
	golangci-lint run

get-linter:
    command -v golangci-lint || curl -sfL "https://install.goreleaser.com/github.com/golangci/golangci-lint.sh" | sh -s -- -b $(go env GOPATH)/bin v1.18.0