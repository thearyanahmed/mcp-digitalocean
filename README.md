# MCP DigitalOcean Integration

This project provides an integration with DigitalOcean's API, exposing resources and tools for managing droplets, CDNs, billing, and other DigitalOcean services. It is built using the [godo](https://github.com/digitalocean/godo) library and the [MCP framework](https://github.com/mark3labs/mcp-go).

## Features

### Resources

Resources are read-only entities that provide information about DigitalOcean resources.

| **Resource Name**       | **URI Template**                          | **Description**                                   | **Supported Operations**         |
|--------------------------|-------------------------------------------|-------------------------------------------------|----------------------------------|
| **Droplets Resource**    | `droplets://{id}`                        | Returns information about a droplet.            | Get droplet details by ID.       |
| **Droplet Actions**      | `droplets://{id}/actions/{action_id}`    | Returns information about a droplet action.     | Get action details.              |
| **Sizes Resource**       | `sizes://all`                            | Returns all available droplet sizes.            | List all droplet sizes.          |
| **Account Resource**     | `account://current`                      | Returns account information.                    | Get account details.             |
| **Balance Resource**     | `balance://current`                      | Returns balance information.                    | Get current balance details.     |
| **Action Resource**      | `actions://{id}`                         | Returns information about a specific action.     | Get action details by ID.        |
| **Images Resource**      | `images://distribution`                  | Returns all available distribution images.       | List all distribution images.    |
| **Specific Image**       | `images://{id}`                          | Returns information about a specific image.      | Get image details by ID.         |
| **CDN Resource**         | `cdn://{id}`                             | Returns information about a CDN.                | Get CDN details by ID.           |
| **Billing History**      | `billing://{last}`                      | Returns billing history for last n bills.                        | Get billing history.             |
| **Certificate Resource** | `certificates://{id}`                    | Returns information about a certificate.        | Get certificate details by ID.   |
| **Domains Resource**     | `domains://{name}`                       | Returns information about a domain.             | Get domain details by name.      |
| **Domain Record**        | `domains://{name}/records/{record_id}`   | Returns information about a domain record.      | Get domain record details.       |
| **Autoscale Pool**       | `autoscale://{id}`                       | Returns information about an autoscale pool.    | Get autoscale pool details.      |
| **Firewall Resource**    | `firewalls://{id}`                       | Returns information about a firewall.           | Get firewall details by ID.      |
| **SSH Key Resource**     | `keys://{id}`                            | Returns information about an SSH key.           | Get SSH key details by ID.       |
| **Regions Resource**     | `regions://all`                          | Returns all available regions.                   | List all regions.                |
| **Reserved IPv4 Resource** | `reserved_ips://{ip}`                  | Returns information about a reserved IPv4.       | Get reserved IPv4 details by IP. |
| **Reserved IPv6 Resource** | `reserved_ipv6://{ip}`                 | Returns information about a reserved IPv6.       | Get reserved IPv6 details by IP. |
| **Partner Attachment**   | `partner_attachment://{id}`              | Returns information about a partner attachment.  | Get partner attachment details.  |
| **VPC Resource**         | `vpcs://{id}`                            | Returns information about a VPC.                | Get VPC details by ID.           |

---

### Tools

Tools provide actions that can be performed on DigitalOcean resources.

| **Tool Name**                          | **Description**                                   | **Parameters**                                                                                     |
|----------------------------------------|---------------------------------------------------|---------------------------------------------------------------------------------------------------|
| **digitalocean-key-create**            | Creates a new SSH key.                            | `Name` (string, required), `PublicKey` (string, required)                                         |
| **digitalocean-key-delete**            | Deletes an SSH key.                               | `ID` (number, required)                                                                          |
| **digitalocean-firewall-create**       | Creates a new firewall.                           | `Name`, `InboundProtocol`, `InboundPortRange`, `InboundSource`, `OutboundProtocol`, `OutboundPortRange`, `OutboundDestination`, `DropletIDs`, `Tags` |
| **digitalocean-firewall-delete**       | Deletes a firewall.                               | `ID` (string, required)                                                                          |
| **digitalocean-droplet-create**        | Creates a new droplet.                            | `Name`, `Size`, `ImageID`, `Region`, `Backup`, `Monitoring`                                      |
| **digitalocean-droplet-delete**        | Deletes a droplet.                                | `ID` (number, required)                                                                          |
| **digitalocean-cdn-create**            | Creates a new CDN.                                | `Origin`, `TTL`, `CustomDomain`                                                                  |
| **digitalocean-cdn-delete**            | Deletes a CDN.                                    | `ID` (string, required)                                                                          |
| **digitalocean-cdn-flush-cache**       | Flushes the cache of a CDN.                       | `ID` (string, required), `Files` (array of strings, required)                                    |
| **digitalocean-domain-create**         | Creates a new domain.                             | `Name`, `IPAddress`                                                                              |
| **digitalocean-domain-delete**         | Deletes a domain.                                 | `Name` (string, required)                                                                        |
| **digitalocean-domain-record-create**  | Creates a new domain record.                      | `Domain`, `Type`, `Name`, `Data`                                                                 |
| **digitalocean-domain-record-delete**  | Deletes a domain record.                          | `Domain`, `RecordID`                                                                             |
| **digitalocean-domain-record-edit**    | Edits a domain record.                            | `Domain`, `RecordID`, `Type`, `Name`, `Data`                                                     |
| **digitalocean-certificate-create**    | Creates a new certificate.                        | `Name`, `PrivateKey`, `LeafCertificate`, `CertificateChain`                                      |
| **digitalocean-certificate-delete**    | Deletes a certificate.                            | `ID` (string, required)                                                                          |
| **digitalocean-certificate-get**       | Retrieves a certificate by ID.                    | `ID` (string, required)                                                                          |
| **digitalocean-autoscale-create**      | Creates a new autoscale pool.                     | `Name`, `Config`, `DropletTemplate`                                                              |
| **digitalocean-autoscale-delete**      | Deletes an autoscale pool.                        | `ID` (string, required)                                                                          |
| **digitalocean-autoscale-update**      | Updates an autoscale pool.                        | `ID`, `Name`, `Config`, `DropletTemplate`                                                       |
| **digitalocean-reserved-ip-reserve**   | Reserves a new IPv4 or IPv6.                      | `Region` (string, required), `Type` (string, required, "ipv4" or "ipv6")                       |
| **digitalocean-reserved-ip-release**   | Releases a reserved IPv4 or IPv6.                 | `IP` (string, required), `Type` (string, required, "ipv4" or "ipv6")                          |
| **digitalocean-reserved-ip-assign**    | Assigns a reserved IP to a droplet.               | `IP` (string, required), `DropletID` (number, required), `Type` (string, required, "ipv4" or "ipv6") |
| **digitalocean-reserved-ip-unassign**  | Unassigns a reserved IP from a droplet.           | `IP` (string, required), `Type` (string, required, "ipv4" or "ipv6")                          |
| **digitalocean-partner-attachment-create** | Creates a new partner attachment.                | `Name` (string, required), `Region` (string, required), `Bandwidth` (number, required)           |
| **digitalocean-partner-attachment-delete** | Deletes a partner attachment.                    | `ID` (string, required)                                                                          |
| **digitalocean-partner-attachment-get-service-key** | Gets the service key of a partner attachment.    | `ID` (string, required)                                                                          |
| **digitalocean-partner-attachment-get-bgp-config** | Gets the BGP configuration of a partner attachment. | `ID` (string, required)                                                                       |
| **digitalocean-partner-attachment-update** | Updates a partner attachment.                    | `ID` (string, required), `Name` (string, required), `VPCIDs` (array of strings, required)        |
| **digitalocean-vpc-create**            | Creates a new VPC.                                | `Name` (string, required), `Region` (string, required)                                           |
| **digitalocean-vpc-list-members**      | Lists members of a VPC.                           | `ID` (string, required)                                                                          |
| **digitalocean-vpc-delete**            | Deletes a VPC.                                    | `ID` (string, required)                                                                          |

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
