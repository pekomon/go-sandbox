APP := snake

.PHONY: deps build test cover clean ensure-tidy check-go

# Minimal Go >= 1.25 guard
check-go:
	@v=$$(go version | awk '{print $$3}' | sed 's/^go//' | cut -d. -f1,2); \
	req="1.25"; \
	awk 'BEGIN { split("'"$$v"'", a, "."); split("'"$$req"'", b, "."); cur=a[1]*100+a[2]; need=b[1]*100+b[2]; if (cur<need) exit 1 }' \
	&& exit 0 || (echo "Go $$v detected; Go >= $$req required. Please install Go 1.25+ or run in CI."; exit 2)

deps: check-go
	@go mod tidy

# Build only the main package (to be added in a later PR)
build: check-go deps
	@echo "Building $(APP)"
	@mkdir -p bin
	@go build -o bin/$(APP) ./cmd/$(APP)

test: check-go deps
	@go test ./...

cover: check-go deps
	@go test ./... -coverprofile=cover.out -covermode=atomic
	@go tool cover -func=cover.out | tail -n +1

clean:
	@rm -rf bin cover.out

ensure-tidy: check-go
	@go mod tidy
	@git diff --quiet -- go.mod go.sum || (echo "go.mod/go.sum not tidy"; exit 1)
