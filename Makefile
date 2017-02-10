TEST = $$(go list ./... | grep -v '/vendor/')
NAME = $(shell awk -F\" '/^const Name/ { print $$2 }' main.go)
VERSION = $(shell awk -F\" '/^const Version/ { print $$2 }' main.go)

all: clean deps build release

clean:
	rm -rf dist/

deps:
	go get -u github.com/Masterminds/glide
	-glide create
	-glide install

deps-update:
	glide update

build: data
	@mkdir -p dist
	go build -o dist/$(NAME) .

release:
	@mkdir -p releases/v$(VERSION)/linux_amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o $(NAME)
	@mv $(NAME) releases/v$(VERSION)/linux_amd64/

test:
	go test ./
	go test ./cmd

.PHONY: all clean build tag proto test
.PHONY: data deps deps-lock deps-update
.PHONY: docker-deps docker docker-push
