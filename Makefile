.PHONY: dep
dep:
	@ go mod tidy && go mod verify

.PHONY: lint
lint:
	@ golangci-lint run

.PHONY: build
build:
	@ go build -ldflags="-X main.Version=$(shell git describe --tags --abbrev=0)" -o telegram-send .
