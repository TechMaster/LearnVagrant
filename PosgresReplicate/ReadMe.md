# Replicate Postgresql

Hãy đọc kỹ 4 bài viết này

1. [PostgreSQL Replication with Docker](ReplicatePostgresql.md)
2. [Replicating Postgres inside Docker — The How To](ReplicatePostgresInsideDocker.md)
3. [An Easy Recipe for Creating a PostgreSQL Cluster with Docker Swarm](https://blog.crunchydata.com/blog/an-easy-recipe-for-creating-a-postgresql-cluster-with-docker-swarm)
4. [SETTING UP POSTGRESQL STREAMING REPLICATION](https://www.cybertec-postgresql.com/en/setting-up-postgresql-streaming-replication/)

Cấu hình [docker-compose-replication.yml](docker-compose-replication.yml) tại Bitnami
```yaml
version: '2'

services:
  postgresql-master:
    image: docker.io/bitnami/postgresql:11
    ports:
      - '5432'
    volumes:
      - 'postgresql_master_data:/bitnami/postgresql'
    environment:
      - POSTGRESQL_REPLICATION_MODE=master
      - POSTGRESQL_REPLICATION_USER=repl_user
      - POSTGRESQL_REPLICATION_PASSWORD=repl_password
      - POSTGRESQL_USERNAME=postgres
      - POSTGRESQL_PASSWORD=my_password
      - POSTGRESQL_DATABASE=my_database
      - ALLOW_EMPTY_PASSWORD=yes
  postgresql-slave:
    image: docker.io/bitnami/postgresql:11
    ports:
      - '5432'
    depends_on:
      - postgresql-master
    environment:
      - POSTGRESQL_REPLICATION_MODE=slave
      - POSTGRESQL_REPLICATION_USER=repl_user
      - POSTGRESQL_REPLICATION_PASSWORD=repl_password
      - POSTGRESQL_MASTER_HOST=postgresql-master
      - POSTGRESQL_PASSWORD=my_password
      - POSTGRESQL_MASTER_PORT_NUMBER=5432
      - ALLOW_EMPTY_PASSWORD=yes

volumes:
  postgresql_master_data:
    driver: local
```

