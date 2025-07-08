# Droplet MCP Tools

This directory contains tools and resources for managing DigitalOcean Droplets via the MCP Server. These tools enable you to create, modify, control, and query Droplets and related resources such as images and sizes.

## Supported Tools

### Droplet Lifecycle

- **`digitalocean-droplet-create`**  
  Create a new Droplet.  
  **Arguments:**  
  - `Name` (string, required): Name of the Droplet  
  - `Size` (string, required): Slug of the Droplet size (e.g., `s-1vcpu-1gb`)  
  - `ImageID` (number, required): ID of the image to use  
  - `Region` (string, required): Slug of the region (e.g., `nyc3`)  
  - `Backup` (boolean, optional, default: false): Enable backups  
  - `Monitoring` (boolean, optional, default: false): Enable monitoring

- **`digitalocean-droplet-delete`**  
  Delete a Droplet.  
  **Arguments:**  
  - `ID` (number, required): ID of the Droplet to delete

### Power & State

- **`digitalocean-droplet-power-cycle`**  
  Power cycle a Droplet.  
  - `ID` (number, required): Droplet ID

- **`digitalocean-droplet-power-on`**  
  Power on a Droplet.  
  - `ID` (number, required): Droplet ID

- **`digitalocean-droplet-power-off`**  
  Power off a Droplet.  
  - `ID` (number, required): Droplet ID

- **`digitalocean-droplet-shutdown`**  
  Shutdown a Droplet.  
  - `ID` (number, required): Droplet ID

- **`digitalocean-droplet-reboot`**  
  Reboot a Droplet.  
  - `ID` (number, required): Droplet ID

### Image & Disk

- **`digitalocean-droplet-restore`**  
  Restore a Droplet from a backup/snapshot.  
  - `ID` (number, required): Droplet ID  
  - `ImageID` (number, required): Backup/snapshot image ID

- **`digitalocean-droplet-resize`**  
  Resize a Droplet.  
  - `ID` (number, required): Droplet ID  
  - `Size` (string, required): New size slug  
  - `ResizeDisk` (boolean, optional, default: false): Whether to resize the disk

- **`digitalocean-droplet-rebuild`**  
  Rebuild a Droplet from an image.  
  - `ID` (number, required): Droplet ID  
  - `ImageID` (number, required): Image ID

- **`digitalocean-droplet-rebuild-by-slug`**  
  Rebuild a Droplet using an image slug.  
  - `ID` (number, required): Droplet ID  
  - `ImageSlug` (string, required): Image slug

- **`digitalocean-droplet-snapshot`**  
  Take a snapshot of a Droplet.  
  - `ID` (number, required): Droplet ID  
  - `Name` (string, required): Name for the snapshot

### Networking & Features

- **`digitalocean-droplet-enable-ipv6`**  
  Enable IPv6 on a Droplet.  
  - `ID` (number, required): Droplet ID

- **`digitalocean-droplet-enable-private-net`**  
  Enable private networking on a Droplet.  
  - `ID` (number, required): Droplet ID

- **`digitalocean-droplet-enable-backups`**  
  Enable backups on a Droplet.  
  - `ID` (number, required): Droplet ID

- **`digitalocean-droplet-disable-backups`**  
  Disable backups on a Droplet.  
  - `ID` (number, required): Droplet ID

### Miscellaneous

- **`digitalocean-droplet-get-neighbors`**  
  Get neighbors of a Droplet.  
  - `ID` (number, required): Droplet ID

- **`digitalocean-droplet-get-kernels`**  
  Get available kernels for a Droplet.  
  - `ID` (number, required): Droplet ID

- **`digitalocean-droplet-change-kernel`**  
  Change a Droplet's kernel.  
  - `ID` (number, required): Droplet ID  
  - `KernelID` (number, required): Kernel ID

- **`digitalocean-droplet-rename`**  
  Rename a Droplet.  
  - `ID` (number, required): Droplet ID  
  - `Name` (string, required): New name

- **`digitalocean-droplet-password-reset`**  
  Reset password for a Droplet.  
  - `ID` (number, required): Droplet ID

### Tag-based Bulk Actions

- **`digitalocean-droplet-power-cycle-by-tag`**  
- **`digitalocean-droplet-power-on-by-tag`**  
- **`digitalocean-droplet-power-off-by-tag`**  
- **`digitalocean-droplet-shutdown-by-tag`**  
- **`digitalocean-droplet-enable-backups-by-tag`**  
- **`digitalocean-droplet-disable-backups-by-tag`**  
- **`digitalocean-droplet-snapshot-by-tag`**  
- **`digitalocean-droplet-enable-ipv6-by-tag`**  
- **`digitalocean-droplet-enable-private-net-by-tag`**  
  All require:  
  - `Tag` (string, required): Tag of the droplets  
  Some require:  
  - `Name` (string, required): Name for the snapshot (for snapshot-by-tag)

## Supported Resources

- **`droplets://{id}`**  
  Returns information about a specific Droplet by its numeric ID.

- **`droplets://{id}/actions/{action_id}`**  
  Returns information about a specific action performed on a Droplet.

- **`images://distribution`**  
  Returns all available distribution images.

- **`images://{id}`**  
  Returns information about a specific image by its numeric ID.

- **`sizes://all`**  
  Returns all available Droplet sizes.

## Example Queries Using Droplet MCP Tools

- Create a new Droplet named "web-1" in region "nyc3" with size "s-1vcpu-1gb" and image ID 123456.
- Delete Droplet with ID 987654.
- Power cycle Droplet 12345.
- Resize Droplet 12345 to "s-2vcpu-2gb" and resize the disk.
- Take a snapshot of Droplet 12345 named "pre-upgrade".
- Enable IPv6 on Droplet 12345.
- List all available Droplet sizes.
- Show all distribution images.
- Get information about Droplet 12345.
- Get details of action 67890 on Droplet 12345.
- Power off all Droplets tagged "staging".
- Enable backups for all Droplets tagged "production".

## Tool Usage Examples

- **Create a new Droplet:**  
  Tool: `digitalocean-droplet-create`  
  Arguments:  
    - `Name`: `"web-1"`  
    - `Size`: `"s-1vcpu-1gb"`  
    - `ImageID`: `123456`  
    - `Region`: `"nyc3"`  
    - `Backup`: `true`  
    - `Monitoring`: `true`

- **Resize a Droplet:**  
  Tool: `digitalocean-droplet-resize`  
  Arguments:  
    - `ID`: `12345`  
    - `Size`: `"s-2vcpu-2gb"`  
    - `ResizeDisk`: `true`

- **Take a snapshot by tag:**  
  Tool: `digitalocean-droplet-snapshot-by-tag`  
  Arguments:  
    - `Tag`: `"staging"`  
    - `Name`: `"before-maintenance"`

## Notes

- All Droplet and image identifiers are numeric IDs unless otherwise specified.
- Tag-based tools allow you to perform bulk actions on all Droplets with a given tag.
- All responses are returned in JSON format for easy parsing and integration.
- For endpoints that require an ID or tag, replace the placeholder with the appropriate value in your query.
- Use the `sizes://all` and `images://distribution` resources to discover valid size slugs and image IDs for Droplet creation.

---