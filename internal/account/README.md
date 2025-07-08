# Account MCP Tools

This directory contains tools and resources for interacting with DigitalOcean account-related features via the MCP Server. These tools allow you to query and manage account information, billing, invoices, SSH keys, and actions.

## Supported Tools

- `digitalocean-key-create`: Create a new SSH key for your account. This tool lets you add a new public SSH key, which can then be used for Droplet or App Platform deployments.
- `digitalocean-key-delete`: Delete an existing SSH key by its ID. Useful for removing unused or compromised keys from your account.

## Supported Resources

- `account://current`: Retrieve information about the current DigitalOcean account, such as email, UUID, status, and more.
- `balance://current`: Get the current account balance, including outstanding charges and account credits.
- `billing://{last}`: Fetch billing history for the last N months. Replace `{last}` with the number of months you want to retrieve.
- `invoice://{last}`: Retrieve invoice history for the last N months. Replace `{last}` with the number of months you want to retrieve.
- `actions://{id}`: Get information about a specific action (such as Droplet creation, deletion, etc.) by its numeric ID.
- `keys://{id}`: Retrieve information about a specific SSH key by its numeric ID.

## Example Queries Using Account MCP Tools

- Show me my current account information.
- What is my current DigitalOcean balance?
- List my billing history for the last 3 months.
- Get all invoices from the past 6 months.
- Show details for action ID 12345678.
- Retrieve information about SSH key with ID 987654.
- Add a new SSH key to my account.
- Remove SSH key with ID 12345.

## Notes

- Most resources are read-only and provide information about your account and its usage.
- The SSH key tools (`digitalocean-key-create` and `digitalocean-key-delete`) allow you to manage your account's SSH keys directly.
- For endpoints that require an ID or a count (such as `{id}` or `{last}`), replace the placeholder with the appropriate value in your query.
