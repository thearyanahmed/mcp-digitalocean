# MCP DigitalOcean Integration

This project provides an integration with DigitalOcean's API, exposing resources and tools for managing applications, droplets, and other DigitalOcean services. It is built using the [godo](https://github.com/digitalocean/godo) library and the [MCP framework](https://github.com/mark3labs/mcp-go).

## Features

### Resources

Resources are read-only entities that provide information about DigitalOcean resources.

#### Apps Resource
- **URI Template**: `apps://{id}`
- **Description**: Returns information about an app.
- **Supported Operations**:
  - Get app details by ID.

#### App Deployment Resource
- **URI Template**: `apps://{id}/deployments/{deployment_id}`
- **Description**: Returns deployment information for an app.
- **Supported Operations**:
  - Get deployment details by app ID and deployment ID.

#### App Tier Resource
- **URI Template**: `apps://{id}/tier`
- **Description**: Returns tier information for an app.
- **Supported Operations**:
  - Get tier details by app ID.

#### Droplets Resource
- **URI Template**: `droplets://{id}`
- **Description**: Returns information about a droplet.
- **Supported Operations**:
  - Get droplet details by ID.

#### Droplet Actions Resource
- **URI Template**: `droplets://{id}/actions/{action_id}`
- **Description**: Returns information about a specific action performed on a droplet.
- **Supported Operations**:
  - Get action details by droplet ID and action ID.

#### Sizes Resource
- **URI Template**: `sizes://all`
- **Description**: Returns all available droplet sizes.
- **Supported Operations**:
  - List all droplet sizes.

#### Specific Size Resource
- **URI Template**: `sizes://{slug}`
- **Description**: Returns information about a specific droplet size.
- **Supported Operations**:
  - Get details of a droplet size by slug.

#### Account Resource
- **URI Template**: `account://current`
- **Description**: Returns account information.
- **Supported Operations**:
  - Get account details.

#### Balance Resource
- **URI Template**: `balance://current`
- **Description**: Returns balance information.
- **Supported Operations**:
  - Get current balance details.

#### Action Resource
- **URI Template**: `actions://{id}`
- **Description**: Returns information about a specific action.
- **Supported Operations**:
  - Get action details by ID.

---

### Tools

Tools are used to perform actions on DigitalOcean resources, such as creating, updating, or deleting resources.

#### Apps Tools
- **Create App**:
  - **Tool Name**: `digitalocean-app-create`
  - **Description**: Creates a new app.
  - **Parameters**:
    - `Name` (string, required): Name of the app.
    - `Region` (string, required): Region of the app.
    - `Tier` (string, required): Tier of the app.

- **Delete App**:
  - **Tool Name**: `digitalocean-app-delete`
  - **Description**: Deletes an app.
  - **Parameters**:
    - `ID` (string, required): ID of the app to delete.

#### Droplets Tools
- **Create Droplet**:
  - **Tool Name**: `digitalocean-droplet-create`
  - **Description**: Creates a new droplet.
  - **Parameters**:
    - `Name` (string, required): Name of the droplet.
    - `Size` (string, required): Slug of the droplet size (e.g., `s-1vcpu-1gb`).
    - `ImageID` (number, required): ID of the image to use.
    - `Region` (string, required): Slug of the region (e.g., `nyc3`).
    - `Backup` (boolean, optional): Whether to enable backups (default: `false`).
    - `Monitoring` (boolean, optional): Whether to enable monitoring (default: `false`).

- **Delete Droplet**:
  - **Tool Name**: `digitalocean-droplet-delete`
  - **Description**: Deletes a droplet.
  - **Parameters**:
    - `ID` (number, required): ID of the droplet to delete.

- **Power Cycle Droplet**:
  - **Tool Name**: `digitalocean-droplet-power-cycle`
  - **Description**: Power cycles a droplet.
  - **Parameters**:
    - `ID` (number, required): ID of the droplet to power cycle.

- **Power On Droplet**:
  - **Tool Name**: `digitalocean-droplet-power-on`
  - **Description**: Powers on a droplet.
  - **Parameters**:
    - `ID` (number, required): ID of the droplet to power on.

- **Power Off Droplet**:
  - **Tool Name**: `digitalocean-droplet-power-off`
  - **Description**: Powers off a droplet.
  - **Parameters**:
    - `ID` (number, required): ID of the droplet to power off.

- **Resize Droplet**:
  - **Tool Name**: `digitalocean-droplet-resize`
  - **Description**: Resizes a droplet.
  - **Parameters**:
    - `ID` (number, required): ID of the droplet to resize.
    - `Size` (string, required): Slug of the new size (e.g., `s-1vcpu-1gb`).
    - `ResizeDisk` (boolean, optional): Whether to resize the disk (default: `false`).

- **Snapshot Droplet**:
  - **Tool Name**: `digitalocean-droplet-snapshot`
  - **Description**: Takes a snapshot of a droplet.
  - **Parameters**:
    - `ID` (number, required): ID of the droplet.
    - `Name` (string, required): Name for the snapshot.

- **Enable Backups**:
  - **Tool Name**: `digitalocean-droplet-enable-backups`
  - **Description**: Enables backups on a droplet.
  - **Parameters**:
    - `ID` (number, required): ID of the droplet.

- **Disable Backups**:
  - **Tool Name**: `digitalocean-droplet-disable-backups`
  - **Description**: Disables backups on a droplet.
  - **Parameters**:
    - `ID` (number, required): ID of the droplet.

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/mcp-digitalocean.git
   cd mcp-digitalocean
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the project:
   ```bash
   go build ./...
   ```

---

## Usage

1. Configure your DigitalOcean API token in the environment:
   ```bash
   export DO_TOKEN=your_token
   ```

2. Run the MCP server:
   ```bash
   go run cmd/mcp.go
   ```

3. Use the provided tools and resources to interact with DigitalOcean.

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.
