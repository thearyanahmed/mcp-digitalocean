# Common MCP Tools and Resources

This directory contains shared tools, resources, and utility functions that are always included as part of the DigitalOcean MCP integration. The contents of this directory provide Utilities and abstractions that are leveraged across all DigitalOcean MCP resource and tools.

## Purpose

The `common` package is designed to:

- Provide utility functions for extracting and parsing resource identifiers from URIs.
- Offer core resources that are universally useful for DigitalOcean interactions (such as region listings).
- Serve as a central place for logic and helpers that are not specific to a single DigitalOcean product or service, but are required by many.

## Included Resources

- **Regions Resource**
  - **`regions://all`**
    Returns a list of all available DigitalOcean regions, including their features and droplet size availability.
    Useful for workflows that need to present or validate region options for resource creation or migration.
