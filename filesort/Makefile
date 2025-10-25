APP := filesort

.PHONY: deps build test cover clean ensure-tidy

deps:
	@go mod tidy

# Build only the main package (to be added in a later PR)
build: deps
	@echo "Building $(APP)"
	@mkdir -p bin
	@go build -o bin/$(APP) ./cmd/filesort

test: deps
	@go test ./...

cover: deps
	@go test ./... -coverprofile=cover.out -covermode=atomic
	@go tool cover -func=cover.out | tail -n +1

clean:
	@rm -rf bin cover.out

ensure-tidy:
	@go mod tidy
	@git diff --quiet -- go.mod go.sum || (echo "go.mod/go.sum not tidy"; exit 1)
