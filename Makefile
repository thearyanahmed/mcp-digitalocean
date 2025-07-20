all: lint test build-dist
build-dist: build-bin dist

build-bin:
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo "goreleaser not found, installing..."; \
		go install github.com/goreleaser/goreleaser@latest; \
	fi
	goreleaser build --auto-snapshot --clean

.PHONY: dist
dist:
	mkdir -p ./scripts/npm/dist
	cp ./dist/*/mcp-digitalocean* ./scripts/npm/dist/
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
