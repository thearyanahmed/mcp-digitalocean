all: lint test build-dist
build-dist: build-bin dist

build-bin:
	goreleaser build --auto-snapshot --clean --skip validate

.PHONY: dist inspector
dist:
	mkdir -p ./scripts/npm/dist
	cp ./README.md ./scripts/npm/README.md
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

inspector: build-dist
	SERVICES=apps,droplets,accounts,networking,insights,spaces,databases,marketplace,doks npm run inspector
