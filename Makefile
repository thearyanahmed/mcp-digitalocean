all: lint build test
build-dist: build-bin dist

MAIN := ./cmd/mcp-digitalocean/main.go
COMMIT := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags --always --dirty)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w \
  -X 'main.Commit=$(COMMIT)' \
  -X 'main.Version=$(VERSION)' \
  -X 'main.Date=$(DATE)'

build-bin:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/mcp-digitalocean-darwin-arm64 $(MAIN)
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/mcp-digitalocean-darwin-amd64 $(MAIN)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/mcp-digitalocean-linux-arm64 $(MAIN)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/mcp-digitalocean-linux-amd64 $(MAIN)
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/mcp-digitalocean-windows-amd64.exe $(MAIN)
	GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/mcp-digitalocean-windows-arm64.exe $(MAIN)
dist:
	mkdir -p ./scripts/npm/dist
	cp ./bin/* ./scripts/npm/dist/
	cp ./internal/apps/spec/*.json ./scripts/npm/dist/
	npm install --prefix ./scripts/npm/

lint:
	revive -config revive.toml ./...

test:
	go test -v ./...

format:
	gofmt -w .
	@echo "Code formatted successfully."

format-check:
	bash -c 'diff -u <(echo -n) <(gofmt -d ./)'

gen:
	go generate ./...
