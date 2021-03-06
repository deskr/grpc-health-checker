TEST = $$(go list ./... | grep -v '/vendor/')
NAME = $(shell awk -F\" '/^const Name/ { print $$2 }' main.go)
VERSION = $(shell awk -F\" '/^const Version/ { print $$2 }' main.go)

all: clean deps build release proto

clean:
	rm -rf dist/

deps:
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go
	-glide create
	-glide install

proto:
	protoc --go_out=plugins=grpc:. grpc_health_v1/health.proto

deps-update:
	glide update

build: data
	@mkdir -p dist
	go build -o dist/$(NAME) .

release:
	@mkdir -p releases
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $(NAME)
	tar -cvzf releases/$(NAME)_v$(VERSION)_linux_amd64.tar.gz $(NAME)
	rm $(NAME)

test:
	go test ./
	go test ./cmd

.PHONY: all clean build tag proto test
.PHONY: data deps deps-lock deps-update
.PHONY: docker-deps docker docker-push
