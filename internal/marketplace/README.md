# Marketplace MCP Tools

This directory contains tools for managing DigitalOcean Marketplace services via the MCP Server. All operations are exposed as tools with argument-based input—no resource URIs are used.

---

## Supported Tools

### 1-Click Application Tools

- **digitalocean-1-click-list**  
  List available 1-click applications from the DigitalOcean marketplace.  
  **Arguments:**  
  - `type` (string, optional, default: "droplet"): Type of 1-click apps to list (e.g., "droplet", "kubernetes")

- **digitalocean-1-click-kuberenetes-app-install**  
  Install 1-click applications on a Kubernetes cluster.  
  **Arguments:**  
  - `cluster_uuid` (string, required): UUID of the Kubernetes cluster to install apps on  
  - `app_slugs` (array, required): Array of app slugs to install

---

## Example Usage

- **List all droplet 1-click apps:**  
  Tool: `digitalocean-1-click-list`  
  Arguments: `{}`

- **List Kubernetes 1-click apps:**  
  Tool: `digitalocean-1-click-list`  
  Arguments:  
  - `type`: `"kubernetes"`

- **Install single app on Kubernetes cluster:**  
  Tool: `digitalocean-1-click-kuberenetes-app-install`  
  Arguments:  
  - `cluster_uuid`: `"k8s-1234567890abcdef"`  
  - `app_slugs`: `["wordpress"]`

- **Install multiple apps on Kubernetes cluster:**  
  Tool: `digitalocean-1-click-kuberenetes-app-install`  
  Arguments:  
  - `cluster_uuid`: `"k8s-1234567890abcdef"`  
  - `app_slugs`: `["wordpress", "mysql", "redis"]`

---

## JSON-RPC Examples

- **List droplet 1-click apps:**

  ```json
  {"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"digitalocean-1-click-list","arguments":{}}}
  ```

- **List Kubernetes 1-click apps:**

  ```json
  {"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"digitalocean-1-click-list","arguments":{"type":"kubernetes"}}}
  ```

- **Install apps on Kubernetes cluster:**

  ```json
  {"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"digitalocean-1-click-kuberenetes-app-install","arguments":{"cluster_uuid":"k8s-1234567890abcdef","app_slugs":["wordpress","nginx"]}}}
  ```

---

## Notes

- For Kubernetes app installation, you need the UUID of an existing Kubernetes cluster.
- Use the `digitalocean-1-click-list` tool with `type: "kubernetes"` to see available Kubernetes 1-click apps.
- A valid DigitalOcean API token is required for all operations.
