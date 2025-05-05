.PHONY: dep
dep:
	@ go mod tidy && go mod verify

.PHONY: lint
lint:
	@ golangci-lint run

.PHONY: build
build:
	@ go build -o ./telegram-send .
