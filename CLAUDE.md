# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the **MCP DigitalOcean Integration** - a Go-based MCP (Model Context Protocol) server that provides comprehensive tools for managing DigitalOcean resources through the DigitalOcean API. The project uses the `godo` library for API interactions and the `mcp-go` framework for MCP server functionality.

## Development Commands

### Testing
```bash
make test           # Run all tests
go test -v ./...    # Run tests with verbose output
```

### Linting
```bash
make lint          # Run revive linter with revive.toml config
revive -config revive.toml ./...
```

### Formatting
```bash
make format        # Format code with gofmt
make format-check  # Check if code is properly formatted
```

### Building
```bash
make build-bin     # Build binaries using goreleaser
make build-dist    # Build binaries and create npm distribution
make dist          # Create npm package distribution
make all           # Run lint, test, and build-dist
```

### Code Generation
```bash
make gen           # Run go generate for all packages
go generate ./...  # Generate mocks and other generated code
```

### MCP Inspector (for App Platform deployment)
```bash
npm start          # Start MCP inspector with apps service only
npm run dev        # Start MCP inspector with apps,droplets,databases services
npm run inspector  # Start MCP inspector with custom services via SERVICES env var
```

## Architecture

### Core Components

- **Main Entry Point**: `cmd/mcp-digitalocean/main.go` - Sets up the MCP server with DigitalOcean client authentication
- **Service Registry**: `internal/registry.go` - Central registration system for all service tools
- **Service Modules**: `internal/{service}/` - Each DigitalOcean service (apps, droplets, networking, etc.) has its own module

### Service Architecture

The codebase is organized into service-specific modules under `internal/`:

- `apps/` - App Platform management (deployments, configurations)
- `droplet/` - Droplet (VM) management and actions
- `account/` - Account information, billing, SSH keys
- `networking/` - Domains, certificates, firewalls, VPCs
- `dbaas/` - Database clusters (Postgres, MySQL, Redis, etc.)
- `insights/` - Monitoring and alerting tools
- `spaces/` - Object storage and CDN management
- `marketplace/` - DigitalOcean Marketplace applications
- `doks/` - Kubernetes cluster management
- `common/` - Shared utilities (regions, etc.)

### Tool Registration Pattern

Each service follows a consistent pattern:
1. Implements a tool struct with `Tools()` method returning MCP tools
2. Tools are registered in `registry.go` via service-specific register functions
3. Services are activated via `--services` flag (e.g., `--services apps,droplets`)

### Testing Strategy

- Each service has corresponding `*_test.go` files
- Uses testify for assertions and go.uber.org/mock for mocking
- Mock interfaces are generated using `go:generate` directives

### Key Dependencies

- `github.com/digitalocean/godo` - DigitalOcean API client
- `github.com/mark3labs/mcp-go` - MCP server framework
- `github.com/invopop/jsonschema` - JSON schema generation for app specs

## Service Configuration

Services are selectively enabled using the `--services` flag:
```bash
# Enable specific services
npx @digitalocean/mcp --services apps,droplets,databases

# All services are enabled if none specified
```

Supported services: `apps`, `droplets`, `accounts`, `networking`, `insights`, `spaces`, `databases`, `marketplace`, `doks`

## Authentication

Requires `DIGITALOCEAN_API_TOKEN` environment variable or `--digitalocean-api-token` flag.

## Distribution

The project distributes as an npm package (`@digitalocean/mcp`) with platform-specific Go binaries, allowing Node.js-based MCP clients to use the Go-based server.