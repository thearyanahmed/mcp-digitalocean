# DBaaS MCP Tools

This directory contains tools and resources for managing DigitalOcean managed database resources via the MCP Server. These tools enable you to create, modify, and query clusters, users, firewalls, configuration, topics, and other database-related resources.

---

## Supported Tools

### Cluster Tools

- **`digitalocean-dbaas-cluster-list`**

  - Get list of clusters.
  - **Arguments:**
    - `page` (optional, integer as string): Page number for pagination
    - `per_page` (optional, integer): Number of results per page

- **`digitalocean-dbaas-cluster-get`**

  - Get a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The ID of the cluster to retrieve

- **`digitalocean-dbaas-cluster-get-ca`**

  - Get the CA certificate for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The ID of the cluster to retrieve the CA for

- **`digitalocean-dbaas-cluster-create`**

  - Create a new database cluster.
  - **Arguments:**
    - `name` (required): The name of the cluster
    - `engine` (required): The engine slug (e.g., valkey, pg, mysql, etc.)
    - `version` (required): The version of the engine
    - `region` (required): The region slug (e.g., nyc1)
    - `size` (required): The size slug (e.g., db-s-2vcpu-4gb)
    - `num_nodes` (required): The number of nodes
    - `tags` (optional): Comma-separated tags to apply to the cluster

- **`digitalocean-dbaas-cluster-delete`**

  - Delete a database cluster by its ID.
  - **Arguments:**
    - `ID` (required): The ID of the cluster to delete

- **`digitalocean-dbaas-cluster-resize`**

  - Resize a database cluster by its ID. At least one of size, num_nodes, or storage_size_mib must be provided.
  - **Arguments:**
    - `ID` (required): The ID of the cluster to resize
    - `size` (optional): The new size slug (e.g., db-s-2vcpu-4gb)
    - `num_nodes` (optional): The new number of nodes
    - `storage_size_mib` (optional): The new storage size in MiB

- **`digitalocean-dbaas-cluster-list-backups`**

  - List backups for a database cluster by its ID.
  - **Arguments:**
    - `ID` (required): The ID of the cluster to list backups for
    - `page` (optional, integer as string): Page number for pagination
    - `per_page` (optional, integer): Number of results per page

- **`digitalocean-dbaas-cluster-list-options`**

  - List available database options (engines, versions, sizes, regions, etc) for DigitalOcean managed databases.
  - **Arguments:** None

- **`digitalocean-dbaas-cluster-upgrade-major-version`**

  - Upgrade the major version of a database cluster by its ID. Requires the target version.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `version` (required): The target major version to upgrade to (e.g., 15 for PostgreSQL)

- **`digitalocean-dbaas-cluster-start-online-migration`**

  - Start an online migration for a database cluster by its ID. Accepts source_json (DatabaseOnlineMigrationConfig as JSON, required), disable_ssl (optional, bool as boolean), and ignore_dbs (optional, comma-separated).
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `source_json` (required): DatabaseOnlineMigrationConfig as JSON
    - `disable_ssl` (optional): Disable SSL for migration (bool as boolean)
    - `ignore_dbs` (optional): Comma-separated list of DBs to ignore

- **`digitalocean-dbaas-cluster-stop-online-migration`**

  - Stop an online migration for a database cluster by its ID and migration_id.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `migration_id` (required): The migration ID to stop

- **`digitalocean-dbaas-cluster-get-online-migration-status`**
  - Get the online migration status for a database cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

### Firewall Tools

- **`digitalocean-dbaas-cluster-get-firewall-rules`**

  - Get the firewall rules for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`digitalocean-dbaas-cluster-update-firewall-rules`**
  - Update the firewall rules for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `rules_json` (required): JSON array of firewall rules to set

### Kafka Tools

- **`digitalocean-dbaas-cluster-list-topics`**

  - List topics for a database cluster by its ID (Kafka clusters). Supports all ListOptions: page, per_page, with_projects, only_deployed, public_only, usecases (comma-separated).
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `page` (optional, integer as string): Page number for pagination
    - `per_page` (optional, integer): Number of results per page
    - `with_projects` (optional, bool as string): Whether to include project_id fields
    - `only_deployed` (optional, bool as string): Only list deployed agents
    - `public_only` (optional, bool as string): Include only public models
    - `usecases` (optional): Comma-separated usecases to filter

- **`digitalocean-dbaas-cluster-create-topic`**

  - Create a topic for a Kafka database cluster by its ID. Accepts name (required), partition_count, replication_factor, and config_json (TopicConfig as JSON, all optional).
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `name` (required): The topic name to create
    - `partition_count` (optional, integer as string): Number of partitions
    - `replication_factor` (optional, integer as string): Replication factor
    - `config_json` (optional): TopicConfig as JSON

- **`digitalocean-dbaas-cluster-get-topic`**

  - Get a topic for a Kafka database cluster by its ID and topic name.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `name` (required): The topic name to get

- **`digitalocean-dbaas-cluster-delete-topic`**

  - Delete a topic for a Kafka database cluster by its ID and topic name.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `name` (required): The topic name to delete

- **`digitalocean-dbaas-cluster-update-topic`**

  - Update a topic for a Kafka database cluster by its ID and topic name. Accepts partition_count, replication_factor, and config_json (TopicConfig as JSON, all optional).
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `name` (required): The topic name to update
    - `partition_count` (optional, integer as string): Number of partitions
    - `replication_factor` (optional, integer as string): Replication factor
    - `config_json` (optional): TopicConfig as JSON

- **`digitalocean-dbaas-cluster-get-kafka-config`**

  - Get the Kafka config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`digitalocean-dbaas-cluster-update-kafka-config`**
  - Update the Kafka config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the KafkaConfig to set

### Mongo Tools

- **`digitalocean-dbaas-cluster-get-mongodb-config`**

  - Get the MongoDB config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`digitalocean-dbaas-cluster-update-mongodb-config`**
  - Update the MongoDB config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the MongoDBConfig to set

### MySQL Tools

- **`digitalocean-dbaas-cluster-get-mysql-config`**

  - Get the MySQL config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`digitalocean-dbaas-cluster-update-mysql-config`**

  - Update the MySQL config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the MySQLConfig to set

- **`digitalocean-dbaas-cluster-get-sql-mode`**

  - Get the SQL mode for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`digitalocean-dbaas-cluster-set-sql-mode`**
  - Set the SQL mode for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `modes` (required): Comma-separated SQL modes to set

### Opensearch Tools

- **`digitalocean-dbaas-cluster-get-opensearch-config`**

  - Get the Opensearch config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`digitalocean-dbaas-cluster-update-opensearch-config`**
  - Update the Opensearch config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the OpensearchConfig to set

### Postgres Tools

- **`digitalocean-dbaas-cluster-get-postgresql-config`**

  - Get the PostgreSQL config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`digitalocean-dbaas-cluster-update-postgresql-config`**
  - Update the PostgreSQL config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the PostgreSQLConfig to set

### Redis Tools

- **`digitalocean-dbaas-cluster-get-redis-config`**

  - Get the Redis config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`digitalocean-dbaas-cluster-update-redis-config`**
  - Update the Redis config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the RedisConfig to set

### User Tools

- **`digitalocean-dbaas-cluster-get-user`**

  - Get a database user by cluster ID and user name.
  - **Arguments:**
    - `ID` (required): The cluster ID (UUID)
    - `user` (required): The user name

- **`digitalocean-dbaas-cluster-list-users`**

  - List database users for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster ID (UUID)
    - `page` (optional): Page number for pagination
    - `per_page` (optional): Number of results per page

- **`digitalocean-dbaas-cluster-create-user`**

  - Create a database user for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster ID (UUID)
    - `name` (required): The user name
    - `mysql_auth_plugin` (optional): MySQL auth plugin (e.g., mysql_native_password)
    - `settings_json` (optional): Raw JSON for DatabaseUserSettings

- **`digitalocean-dbaas-cluster-update-user`**
  - Update a database user for a cluster by its ID and user name.
  - **Arguments:**
    - `ID` (required): The cluster ID (UUID)
    - `user`

---

## Example queries using DBaaS MCP Tools

Below are some example natural language queries you might use, along with the corresponding tool and arguments. Each example shows how a user might phrase a request, and how it maps to the underlying tool and parameters.

### Clusters

| Example Query                                 | Tool                              | Arguments                                                                                                            |
| --------------------------------------------- | --------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Could you please list all dbaas clusters?     | digitalocean-dbaas-cluster-list   | `{ "page": "1", "per_page": 10 }`                                                                                    |
| Show me details for cluster `<cluster-uuid>`  | digitalocean-dbaas-cluster-get    | `{ "ID": "<cluster-uuid>" }`                                                                                         |
| Create a dbaas cluster called 'my-db' in nyc1 | digitalocean-dbaas-cluster-create | `{ "name": "my-db", "engine": "mysql", "version": "8", "region": "nyc1", "size": "db-s-1vcpu-1gb", "num_nodes": 1 }` |
| Delete the cluster `<cluster-uuid>`           | digitalocean-dbaas-cluster-delete | `{ "ID": "<cluster-uuid>" }`                                                                                         |
| Resize cluster `<cluster-uuid>` to 2 nodes    | digitalocean-dbaas-cluster-resize | `{ "ID": "<cluster-uuid>", "num_nodes": 2 }`                                                                         |

### Users

| Example Query                                            | Tool                                   | Arguments                                                                  |
| -------------------------------------------------------- | -------------------------------------- | -------------------------------------------------------------------------- |
| List all users for cluster `<cluster-uuid>`              | digitalocean-dbaas-cluster-list-users  | `{ "ID": "<cluster-uuid>" }`                                               |
| Add a user named 'readonly' to cluster `<cluster-uuid>`  | digitalocean-dbaas-cluster-create-user | `{ "ID": "<cluster-uuid>", "name": "readonly" }`                           |
| Remove the user 'readonly' from cluster `<cluster-uuid>` | digitalocean-dbaas-cluster-delete-user | `{ "ID": "<cluster-uuid>", "user": "readonly" }`                           |
| Update user 'readonly' with new settings                 | digitalocean-dbaas-cluster-update-user | `{ "ID": "<cluster-uuid>", "user": "readonly", "settings_json": "{...}" }` |

### Firewalls

| Example Query                                             | Tool                                             | Arguments                                                        |
| --------------------------------------------------------- | ------------------------------------------------ | ---------------------------------------------------------------- |
| What are the firewall rules for cluster `<cluster-uuid>`? | digitalocean-dbaas-cluster-get-firewall-rules    | `{ "ID": "<cluster-uuid>" }`                                     |
| Update the firewall rules for cluster `<cluster-uuid>`    | digitalocean-dbaas-cluster-update-firewall-rules | `{ "ID": "<cluster-uuid>", "rules_json": "[ { ...rule... } ]" }` |

### Configuration

| Example Query                                             | Tool                                                | Arguments                                              |
| --------------------------------------------------------- | --------------------------------------------------- | ------------------------------------------------------ |
| Show me the MySQL config for cluster `<cluster-uuid>`     | digitalocean-dbaas-cluster-get-mysql-config         | `{ "ID": "<cluster-uuid>" }`                           |
| Update the MongoDB config for cluster `<cluster-uuid>`    | digitalocean-dbaas-cluster-update-mongodb-config    | `{ "ID": "<cluster-uuid>", "config_json": "{ ... }" }` |
| Get the Redis config for cluster `<cluster-uuid>`         | digitalocean-dbaas-cluster-get-redis-config         | `{ "ID": "<cluster-uuid>" }`                           |
| Update the PostgreSQL config for cluster `<cluster-uuid>` | digitalocean-dbaas-cluster-update-postgresql-config | `{ "ID": "<cluster-uuid>", "config_json": "{ ... }" }` |

### Kafka Topics

| Example Query                                                      | Tool                                    | Arguments                                        |
| ------------------------------------------------------------------ | --------------------------------------- | ------------------------------------------------ |
| List all topics in Kafka cluster `<cluster-uuid>`                  | digitalocean-dbaas-cluster-list-topics  | `{ "ID": "<cluster-uuid>" }`                     |
| Create a topic called 'my-topic' in Kafka cluster `<cluster-uuid>` | digitalocean-dbaas-cluster-create-topic | `{ "ID": "<cluster-uuid>", "name": "my-topic" }` |
| Delete the topic 'my-topic' from Kafka cluster `<cluster-uuid>`    | digitalocean-dbaas-cluster-delete-topic | `{ "ID": "<cluster-uuid>", "name": "my-topic" }` |

---

Feel free to use these queries as a starting point for interacting with the DBaaS MCP tools!
