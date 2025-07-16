# Droplet MCP Tools

This directory contains tools for managing DigitalOcean Droplets, Images, and Sizes via the MCP Server. All operations are exposed as tools with argument-based input—no resource URIs are used. Pagination and filtering are supported where applicable.

---

## Supported Tools

### Droplet Tools

- **digitalocean-droplet-create**  
  Create a new Droplet.  
  **Arguments:**  
  - `Name` (string, required): Name of the Droplet  
  - `Size` (string, required): Slug of the Droplet size (e.g., `s-1vcpu-1gb`)  
  - `ImageID` (number, required): ID of the image to use  
  - `Region` (string, required): Slug of the region (e.g., `nyc3`)  
  - `Backup` (boolean, optional, default: false): Enable backups  
  - `Monitoring` (boolean, optional, default: false): Enable monitoring

- **digitalocean-droplet-delete**  
  Delete a Droplet.  
  **Arguments:**  
  - `ID` (number, required): ID of the Droplet to delete

- **digitalocean-droplet-get**  
  Get information about a specific Droplet by its ID.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID

---

### Droplet Actions Tools

- **digitalocean-droplet-action-get**  
  Get information about a specific action performed on a Droplet.  
  **Arguments:**  
  - `DropletID` (number, required): Droplet ID  
  - `ActionID` (number, required): Action ID

- **digitalocean-droplet-action-reboot**  
  Reboot a Droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID

- **digitalocean-droplet-action-password-reset**  
  Reset password for a Droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID

- **digitalocean-droplet-action-rename**  
  Rename a Droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID  
  - `Name` (string, required): New name

- **digitalocean-droplet-action-change-kernel**  
  Change a Droplet's kernel.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID  
  - `KernelID` (number, required): Kernel ID

- **digitalocean-droplet-action-enable-ipv6**  
- **digitalocean-droplet-action-enable-private-net**  
- **digitalocean-droplet-action-enable-backups**  
- **digitalocean-droplet-action-disable-backups**  
  Enable/disable features on a Droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID

#### Tag-based Bulk Actions

- **digitalocean-droplet-action-power-cycle-by-tag**  
- **digitalocean-droplet-action-power-on-by-tag**  
- **digitalocean-droplet-action-power-off-by-tag**  
- **digitalocean-droplet-action-shutdown-by-tag**  
- **digitalocean-droplet-action-enable-backups-by-tag**  
- **digitalocean-droplet-action-disable-backups-by-tag**  
- **digitalocean-droplet-action-snapshot-by-tag**  
- **digitalocean-droplet-action-enable-ipv6-by-tag**  
- **digitalocean-droplet-action-enable-private-net-by-tag**  
  All require:  
  - `Tag` (string, required): Tag of the droplets  
  Some require:  
  - `Name` (string, required): Name for the snapshot (for snapshot-by-tag)

---

### Additional Droplet Actions Tools

- **digitalocean-droplet-action-rebuild-by-slug**  
  Rebuild a droplet using an image slug.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID  
  - `ImageSlug` (string, required): Slug of the image to rebuild from

- **digitalocean-droplet-action-power-cycle**  
  Power cycle a droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID

- **digitalocean-droplet-action-power-on**  
  Power on a droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID

- **digitalocean-droplet-action-power-off**  
  Power off a droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID

- **digitalocean-droplet-action-shutdown**  
  Shutdown a droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID

- **digitalocean-droplet-action-restore**  
  Restore a droplet from a backup/snapshot.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID  
  - `ImageID` (number, required): ID of the backup/snapshot image

- **digitalocean-droplet-action-resize**  
  Resize a droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID  
  - `Size` (string, required): Slug of the new size (e.g., s-1vcpu-1gb)  
  - `ResizeDisk` (boolean, optional, default: false): Whether to resize the disk

- **digitalocean-droplet-action-rebuild**  
  Rebuild a droplet from an image.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID  
  - `ImageID` (number, required): ID of the image to rebuild from

- **digitalocean-droplet-action-snapshot**  
  Take a snapshot of a droplet.  
  **Arguments:**  
  - `ID` (number, required): Droplet ID  
  - `Name` (string, required): Name for the snapshot

---

### Image Tools

- **digitalocean-image-list**  
  List all available distribution images. Supports pagination.  
  **Arguments:**  
  - `Page` (number, default: 1): Page number  
  - `PerPage` (number, default: 50): Items per page

- **digitalocean-image-get**  
  Get a specific image by its numeric ID.  
  **Arguments:**  
  - `ID` (number, required): Image ID

---

### Size Tools

- **digitalocean-size-list**  
  List all available Droplet sizes. Supports pagination.  
  **Arguments:**  
  - `Page` (number, default: 1): Page number  
  - `PerPage` (number, default: 50): Items per page

---

## Example Usage

- **Create a new Droplet:**  
  Tool: `digitalocean-droplet-create`  
  Arguments:  
    - `Name`: `"web-1"`  
    - `Size`: `"s-1vcpu-1gb"`  
    - `ImageID`: `123456`  
    - `Region`: `"nyc3"`  
    - `Backup`: `true`  
    - `Monitoring`: `true`

- **Get a Droplet by ID:**  
  Tool: `digitalocean-droplet-get`  
  Arguments:  
    - `ID`: `12345`

- **Get a Droplet action:**  
  Tool: `digitalocean-droplet-action-get`  
  Arguments:  
    - `DropletID`: `12345`  
    - `ActionID`: `67890`

- **Reboot a Droplet:**  
  Tool: `digitalocean-droplet-action-reboot`  
  Arguments:  
    - `ID`: `12345`

- **Reset password for a Droplet:**  
  Tool: `digitalocean-droplet-action-password-reset`  
  Arguments:  
    - `ID`: `12345`

- **Rename a Droplet:**  
  Tool: `digitalocean-droplet-action-rename`  
  Arguments:  
    - `ID`: `12345`  
    - `Name`: `"new-name"`

- **Change a Droplet's kernel:**  
  Tool: `digitalocean-droplet-action-change-kernel`  
  Arguments:  
    - `ID`: `12345`  
    - `KernelID`: `67890`

- **Perform bulk actions by tag:**  
  Tool: `digitalocean-droplet-action-power-cycle-by-tag`  
  Arguments:  
    - `Tag`: `"web"`

---

## Notes

- All tools use argument-based input; do not use resource URIs.
- Pagination is supported for list endpoints via `Page` and `PerPage` arguments.
- Tag-based tools allow you to perform bulk actions on all Droplets with a given tag.
- All responses are returned in JSON format for easy parsing and integration.
- For endpoints that require an ID or tag, provide the appropriate value in your query.
