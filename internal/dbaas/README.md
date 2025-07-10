# DBaaS MCP Tools

This directory contains tools and resources for managing DigitalOcean managed database resources via the MCP Server. These tools enable you to create, modify, and query clusters, users, firewalls, configuration, topics, and other database-related resources.

---

## Supported Tools

### Cluster Tools

- **`do-dbaas-cluster-list`**
  - Get list of clusters.
  - **Arguments:**
    - `page` (optional, integer as string): Page number for pagination
    - `per_page` (optional, integer as string): Number of results per page

- **`do-dbaas-cluster-get`**
  - Get a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The ID of the cluster to retrieve

- **`do-dbaas-cluster-get-ca`**
  - Get the CA certificate for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The ID of the cluster to retrieve the CA for

- **`do-dbaas-cluster-create`**
  - Create a new database cluster.
  - **Arguments:**
    - `name` (required): The name of the cluster
    - `engine` (required): The engine slug (e.g., valkey, pg, mysql, etc.)
    - `version` (required): The version of the engine
    - `region` (required): The region slug (e.g., nyc1)
    - `size` (required): The size slug (e.g., db-s-2vcpu-4gb)
    - `num_nodes` (required): The number of nodes
    - `tags` (optional): Comma-separated tags to apply to the cluster

- **`do-dbaas-cluster-delete`**
  - Delete a database cluster by its ID.
  - **Arguments:**
    - `ID` (required): The ID of the cluster to delete

- **`do-dbaas-cluster-resize`**
  - Resize a database cluster by its ID. At least one of size, num_nodes, or storage_size_mib must be provided.
  - **Arguments:**
    - `ID` (required): The ID of the cluster to resize
    - `size` (optional): The new size slug (e.g., db-s-2vcpu-4gb)
    - `num_nodes` (optional): The new number of nodes
    - `storage_size_mib` (optional): The new storage size in MiB

- **`do-dbaas-cluster-list-backups`**
  - List backups for a database cluster by its ID.
  - **Arguments:**
    - `ID` (required): The ID of the cluster to list backups for
    - `page` (optional, integer as string): Page number for pagination
    - `per_page` (optional, integer as string): Number of results per page

- **`do-dbaas-cluster-list-options`**
  - List available database options (engines, versions, sizes, regions, etc) for DigitalOcean managed databases.
  - **Arguments:** None

- **`do-dbaas-cluster-upgrade-major-version`**
  - Upgrade the major version of a database cluster by its ID. Requires the target version.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `version` (required): The target major version to upgrade to (e.g., 15 for PostgreSQL)

- **`do-dbaas-cluster-start-online-migration`**
  - Start an online migration for a database cluster by its ID. Accepts source_json (DatabaseOnlineMigrationConfig as JSON, required), disable_ssl (optional, bool as string), and ignore_dbs (optional, comma-separated).
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `source_json` (required): DatabaseOnlineMigrationConfig as JSON
    - `disable_ssl` (optional): Disable SSL for migration (bool as string)
    - `ignore_dbs` (optional): Comma-separated list of DBs to ignore

- **`do-dbaas-cluster-stop-online-migration`**
  - Stop an online migration for a database cluster by its ID and migration_id.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `migration_id` (required): The migration ID to stop

- **`do-dbaas-cluster-get-online-migration-status`**
  - Get the online migration status for a database cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

### Firewall Tools

- **`do-dbaas-cluster-get-firewall-rules`**
  - Get the firewall rules for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`do-dbaas-cluster-update-firewall-rules`**
  - Update the firewall rules for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `rules_json` (required): JSON array of firewall rules to set

### Kafka Tools

- **`do-dbaas-cluster-list-topics`**
  - List topics for a database cluster by its ID (Kafka clusters). Supports all ListOptions: page, per_page, with_projects, only_deployed, public_only, usecases (comma-separated).
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `page` (optional, integer as string): Page number for pagination
    - `per_page` (optional, integer as string): Number of results per page
    - `with_projects` (optional, bool as string): Whether to include project_id fields
    - `only_deployed` (optional, bool as string): Only list deployed agents
    - `public_only` (optional, bool as string): Include only public models
    - `usecases` (optional): Comma-separated usecases to filter

- **`do-dbaas-cluster-create-topic`**
  - Create a topic for a Kafka database cluster by its ID. Accepts name (required), partition_count, replication_factor, and config_json (TopicConfig as JSON, all optional).
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `name` (required): The topic name to create
    - `partition_count` (optional, integer as string): Number of partitions
    - `replication_factor` (optional, integer as string): Replication factor
    - `config_json` (optional): TopicConfig as JSON

- **`do-dbaas-cluster-get-topic`**
  - Get a topic for a Kafka database cluster by its ID and topic name.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `name` (required): The topic name to get

- **`do-dbaas-cluster-delete-topic`**
  - Delete a topic for a Kafka database cluster by its ID and topic name.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `name` (required): The topic name to delete

- **`do-dbaas-cluster-update-topic`**
  - Update a topic for a Kafka database cluster by its ID and topic name. Accepts partition_count, replication_factor, and config_json (TopicConfig as JSON, all optional).
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `name` (required): The topic name to update
    - `partition_count` (optional, integer as string): Number of partitions
    - `replication_factor` (optional, integer as string): Replication factor
    - `config_json` (optional): TopicConfig as JSON

- **`do-dbaas-cluster-get-kafka-config`**
  - Get the Kafka config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`do-dbaas-cluster-update-kafka-config`**
  - Update the Kafka config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the KafkaConfig to set

### Mongo Tools

- **`do-dbaas-cluster-get-mongodb-config`**
  - Get the MongoDB config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`do-dbaas-cluster-update-mongodb-config`**
  - Update the MongoDB config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the MongoDBConfig to set

### MySQL Tools

- **`do-dbaas-cluster-get-mysql-config`**
  - Get the MySQL config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`do-dbaas-cluster-update-mysql-config`**
  - Update the MySQL config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the MySQLConfig to set

- **`do-dbaas-cluster-get-sql-mode`**
  - Get the SQL mode for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`do-dbaas-cluster-set-sql-mode`**
  - Set the SQL mode for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `modes` (required): Comma-separated SQL modes to set

### Opensearch Tools

- **`do-dbaas-cluster-get-opensearch-config`**
  - Get the Opensearch config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`do-dbaas-cluster-update-opensearch-config`**
  - Update the Opensearch config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the OpensearchConfig to set

### Postgres Tools

- **`do-dbaas-cluster-get-postgresql-config`**
  - Get the PostgreSQL config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`do-dbaas-cluster-update-postgresql-config`**
  - Update the PostgreSQL config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the PostgreSQLConfig to set

### Redis Tools

- **`do-dbaas-cluster-get-redis-config`**
  - Get the Redis config for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster UUID

- **`do-dbaas-cluster-update-redis-config`**
  - Update the Redis config for a cluster by its ID. Accepts a JSON string for the config.
  - **Arguments:**
    - `ID` (required): The cluster UUID
    - `config_json` (required): JSON for the RedisConfig to set

### User Tools

- **`do-dbaas-cluster-get-user`**
  - Get a database user by cluster ID and user name.
  - **Arguments:**
    - `ID` (required): The cluster ID (UUID)
    - `user` (required): The user name

- **`do-dbaas-cluster-list-users`**
  - List database users for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster ID (UUID)
    - `page` (optional): Page number for pagination
    - `per_page` (optional): Number of results per page

- **`do-dbaas-cluster-create-user`**
  - Create a database user for a cluster by its ID.
  - **Arguments:**
    - `ID` (required): The cluster ID (UUID)
    - `name` (required): The user name
    - `mysql_auth_plugin` (optional): MySQL auth plugin (e.g., mysql_native_password)
    - `settings_json` (optional): Raw JSON for DatabaseUserSettings

- **`do-dbaas-cluster-update-user`**
  - Update a database user for a cluster by its ID and user name.
  - **Arguments:**
    - `ID` (required): The cluster ID (UUID)
    - `user`