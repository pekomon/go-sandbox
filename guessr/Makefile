APP := guessr

.PHONY: build test cover
build:
@echo "Building $(APP)"
@mkdir -p bin
@go build -o bin/$(APP) ./cmd/guessr

test:
@go test ./...

cover:
@go test ./... -coverprofile=cover.out -covermode=atomic
@go tool cover -func=cover.out | tail -n +1
