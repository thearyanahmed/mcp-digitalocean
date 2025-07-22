# DigitalOcean Kubernetes (DOKS) Tools

This directory provides tools for managing DigitalOcean Kubernetes clusters and node pools via the MCP Server. All operations are exposed as tools with argument-based input—no resource URIs are used. Pagination and filtering are supported where applicable.

---

## Supported Tools

### Cluster Tools

- **digitalocean-doks-get-cluster**  
  Get information about a specific Kubernetes cluster.  
  **Arguments:**
    - `ClusterID` (string, required): ID of the cluster

- **digitalocean-doks-list-clusters**  
  List all Kubernetes clusters.  
  **Arguments:**
    - `Page` (number, default: 1): Page number
    - `PerPage` (number, default: 20): Items per page

- **digitalocean-doks-create-cluster**  
  Create a new Kubernetes cluster.  
  **Arguments:**
    - See schema in `spec/cluster-create-schema.json`

- **digitalocean-doks-update-cluster**  
  Update a Kubernetes cluster.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID
    - `Name` (string, optional): New name
    - `MaintenancePolicy` (object, optional): Maintenance window
    - `AutoUpgrade` (boolean, optional): Enable auto-upgrade
    - `SurgeUpgrade` (boolean, optional): Enable surge upgrades
    - `Tags` (array, optional): Tags

- **digitalocean-doks-delete-cluster**  
  Delete a Kubernetes cluster.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID

- **digitalocean-doks-upgrade-cluster**  
  Upgrade a Kubernetes cluster.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID
    - `VersionSlug` (string, required): Kubernetes version

- **digitalocean-doks-get-cluster-upgrades**  
  Get available upgrades for a cluster.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID

- **digitalocean-doks-get-kubeconfig**  
  Get kubeconfig for a cluster.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID

- **digitalocean-doks-get-credentials**  
  Get credentials for a cluster.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID

---

### Node Pool Tools

- **digitalocean-doks-create-nodepool**  
  Create a new node pool in a cluster.  
  **Arguments:**
    - See schema in `spec/node-pool-create-schema.json`

- **digitalocean-doks-get-nodepool**  
  Get a node pool in a cluster.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID
    - `NodePoolID` (string, required): Node pool ID

- **digitalocean-doks-list-nodepools**  
  List all node pools in a cluster.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID

- **digitalocean-doks-update-nodepool**  
  Update a node pool in a cluster.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID
    - `NodePoolID` (string, required): Node pool ID
    - `Name` (string, optional): New name
    - `Count` (number, optional): Number of nodes
    - `Tags` (array, optional): Tags
    - `Labels` (object, optional): Kubernetes labels
    - `Taints` (array, optional): Kubernetes taints
    - `AutoScale` (boolean, optional): Enable auto-scaling
    - `MinNodes` (number, optional): Minimum nodes
    - `MaxNodes` (number, optional): Maximum nodes

- **digitalocean-doks-delete-nodepool**  
  Delete a node pool in a cluster.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID
    - `NodePoolID` (string, required): Node pool ID

- **digitalocean-doks-delete-node**  
  Delete a node from a node pool.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID
    - `NodePoolID` (string, required): Node pool ID
    - `NodeID` (string, required): Node ID
    - `SkipDrain` (boolean, optional): Skip draining
    - `Replace` (boolean, optional): Replace node

- **digitalocean-doks-recycle-nodes**  
  Recycle specific nodes in a node pool.  
  **Arguments:**
    - `ClusterID` (string, required): Cluster ID
    - `NodePoolID` (string, required): Node pool ID
    - `NodeIDs` (array, required): List of node IDs

---

## Example Usage

- **Get a cluster:**  
  Tool: `digitalocean-doks-get-cluster`  
  Arguments:
    - `ClusterID`: `"abcd-1234"`

- **List clusters:**  
  Tool: `digitalocean-doks-list-clusters`  
  Arguments:
    - `Page`: `1`
    - `PerPage`: `20`

- **Create a node pool:**  
  Tool: `digitalocean-doks-create-nodepool`  
  Arguments:
    - See `spec/node-pool-create-schema.json`

- **Recycle nodes:**  
  Tool: `digitalocean-doks-recycle-nodes`  
  Arguments:
    - `ClusterID`: `"abcd-1234"`
    - `NodePoolID`: `"np-5678"`
    - `NodeIDs`: `["node-1", "node-2"]`

---

## Notes

- All tools use argument-based input; do not use resource URIs.
- Pagination is supported for list endpoints via `Page` and `PerPage` arguments.
- All responses are returned in JSON format for easy parsing and integration.
- For endpoints that require an ID, provide the appropriate value in your query.
- Schemas for cluster and node pool creation are found in the `spec/` directory.