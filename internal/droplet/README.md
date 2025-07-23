# Droplet MCP Tools

This directory contains tools for managing Droplets, Images, and Sizes via the MCP Server. All operations are exposed as tools with argument-based input—no resource URIs are used. Pagination and filtering are supported where applicable.

---

## Supported Tools

### Droplet Tools

- **droplet-create**  
  Create a new Droplet.  
  **Arguments:**  
  - `Name` (string, required): Name of the Droplet  
  - `Size` (string, required): Slug of the Droplet size (e.g., `s-1vcpu-1gb`)  
  - `ImageID` (number, required): ID of the image to use  
  - `Region` (string, required): Slug of the region (e.g., `nyc3`)  
  - `Backup` (boolean, optional, default: false): Enable backups  
  - `Monitoring` (boolean, optional, default: false): Enable monitoring

- **droplet-delete**  
  Delete a Droplet.  
  **Arguments:**  
  - `ID` (number, required): ID of the Droplet to delete

- **droplet-get**  
  Get information about a specific Droplet by its ID.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID

- **droplet-list**  
  List all droplets for the user. Supports pagination.  
  **Arguments:**  
  - `Page` (number, default: 1): Page number  
  - `PerPage` (number, default: 50): Items per page

---

### Droplet Actions Tools

- **droplet-action**  
  Get information about a specific action performed on a Droplet.  
  **Arguments:**  
  - `DropletID` (number, required): Droplet ID  
  - `ActionID` (number, required): Action ID

- **droplet-reboot**  
  Reboot a Droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID

---

## Example Usage

- **Create a new Droplet:**  
  Tool: `droplet-create`  
  Arguments:  
    - `Name`: `"web-1"`  
    - `Size`: `"s-1vcpu-1gb"`  
    - `ImageID`: `123456`  
    - `Region`: `"nyc3"`  
    - `Backup`: `true`  
    - `Monitoring`: `true`

- **Get a Droplet by ID:**  
  Tool: `droplet-get`  
  Arguments:  
    - `ID`: `12345`

- **Reboot a Droplet:**  
  Tool: `droplet-reboot`  
  Arguments:  
    - `ID`: `12345`

---

## Notes

- All tools use argument-based input; do not use resource URIs.
- Pagination is supported for list endpoints via `Page` and `PerPage` arguments.
- Tag-based tools allow you to perform bulk actions on all Droplets with a given tag.
- All responses are returned in JSON format for easy parsing and integration.
- For endpoints that require an ID or tag, provide the appropriate value in your query.
