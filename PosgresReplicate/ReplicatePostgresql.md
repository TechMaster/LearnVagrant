# PostgreSQL Replication with Docker

[Original Artical on Medium](https://medium.com/swlh/postgresql-replication-with-docker-c6a904becf77)

> There are so many ways to setup replication for a PostgreSQL master, but when it comes to docker, it could waste your time. In this article I will tell you how to setup a PostgreSQL master first, then we will add a slave for it using streaming replication method, all in docker containers.

## Heads-Up !
Before going into the details I assume you are familiar with docker and docker-compose service, understand the basics and could work with terminal. Also it’s good to read my article, Tips on Using Docker, since I use them mostly in my configurations.

## Let’s setup a PostgreSQL database
For setting up a postgreSQL database with docker, you could go to their official docker hub page, check the various versions and configurations. I recommend to use docker-compose since it’s easier to manage. Here is a simple docker-compose file for running it.
```yaml
version: '3'
services:
  database:
    image: postgres:13
    container_name: my_postgres_database
    restart: always
    volumes:
        - ./data:/var/lib/postgresql/data
    ports:
      - "127.0.0.1:5432:5432"
    environment:
      - 'POSTGRES_PASSWORD=my_password'
      - 'POSTGRES_DB=my_default_database'
```
## Let’s tune it for production
Now the container should be running and you could access it easily and start using it, but this is not good enough for production. There are some configurations and performance tweaks which needs to be applied based on your hardware ! it’s not essential, but good to have it. You could get all the details you need from https://pgtune.leopard.in.ua website. Then you need to create a postgresql.conf file and add those lines. You should know that these are not all the required configurations for running postgres, these are just the performance tweaks. Now you wonder what are the other configurations? They have answered this question in their docker hub page, first run the command below to get the sample file and then add the configurations you got from pgtune website into it.
```
$ docker run -i --rm postgres cat /usr/share/postgresql/postgresql.conf.sample > my-postgres.conf
```

Here is a final `my-postgres.conf` sample
```
listen_addresses = '*'
# DB Version: 13
# OS Type: linux
# DB Type: web
# Total Memory (RAM): 62 GB
# CPUs num: 12
# Connections num: 1000
# Data Storage: ssd
# Performance tweaks
max_connections = 1000
shared_buffers = 15872MB
effective_cache_size = 47616MB
maintenance_work_mem = 2GB
checkpoint_completion_target = 0.7
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 4063kB
min_wal_size = 1GB
max_wal_size = 4GB
max_worker_processes = 12
max_parallel_workers_per_gather = 4
```

## Let’s setup a PostgreSQL master database !

There are couple of more configurations needed for a master database, first we need to add some lines to my-postgres.conf for replication.

```
# Replication
wal_level = replica
hot_standby = on
max_wal_senders = 10
max_replication_slots = 10
hot_standby_feedback = on
```

Also we need to tell postgres to let our replication user connect to that database and trust it, I assume my replication user is called replicator . There is another configuration file called pg_hba.conf which handles the accesses to the database. Let’s create a custome my-pg_hba.conf with these lines.

```
# TYPE  DATABASE        USER            ADDRESS                 METHOD
# "local" is for Unix domain socket connections only
local   all             all                                     trust
# IPv4 local connections:
host    all             all             127.0.0.1/32            trust
# IPv6 local connections:
host    all             all             ::1/128                 trust
# Allow replication connections from localhost, by a user with the
# replication privilege.
local   replication     all                                     trust
host    replication     all             127.0.0.1/32            trust
host    replication     all             ::1/128                 trust
host    replication     replicator      0.0.0.0/0               trust
host all all all md5

```
Now that we have these files, we need to mount them into the container and let the postgres service use them. Let’s add some more lines to our `docker-compose.yml` file.

```yaml
version: '3'
services:
  database:
    image: postgres:13
    container_name: my_postgres_database
    restart: always
    volumes:
        - ./data:/var/lib/postgresql/data
        - ./my-postgres.conf:/etc/postgresql/postgresql.conf
        - ./my-pg_hba.conf:/etc/postgresql/pg_hba.conf
    ports:
      - "127.0.0.1:5432:5432"
    environment:
      - 'POSTGRES_PASSWORD=my_password'
      - 'POSTGRES_DB=my_default_database'

```

## Let’s Get Ready For Slave !

We are almost there for setting up a replication, there are couple of steps which we need to take.
1. Create the replicator user on master
    ```
    $ CREATE USER replicator WITH REPLICATION ENCRYPTED PASSWORD 'my_replicator_password';
    ```
2. Create the physical replication slot on master
    ```
    $ SELECT * FROM pg_create_physical_replication_slot('replication_slot_slave1');
    ```
To see that the physical replication slot has been created successfully, you could run this query $ SELECT * FROM pg_replication_slots; and you should see something like this.

```
-[ RECORD 1 ]-------+------------------------
slot_name           | replication_slot_slave1
plugin              | 
slot_type           | physical
datoid              | 
database            | 
temporary           | f
active              | f
active_pid          | 
xmin                | 
catalog_xmin        | 
restart_lsn         | 
confirmed_flush_lsn |
```

You could see, since we are not running any slave for this slot, it’s not active yet.

3. We need to get a backup from our master database and restore it for the slave. The best way for doing this is to usepg_basebackup command. Here is the documentation for postgres version 13.
If you don’t want to read all the documentation to find out which flags you should use, just copy the command below and we will go through the flags in this command.
```
$ pg_basebackup -D /tmp/postgresslave -S replication_slot_slave1 -X stream -P -U replicator -Fp -R
```

**What are these flags??!!**

```
-D directory
--pgdata=directory
Sets the target directory to write the output to. pg_basebackup will create this directory (and any missing parent directories) if it does not exist. If it already exists, it must be empty.
-S slotname
--slot=slotname
```

This option can only be used together with `-X stream`. It causes WAL streaming to use the specified replication slot. If the base backup is intended to be used as a streaming-replication standby using a replication slot, the standby should then use the same replication slot name as `primary_slot_name`. This ensures that the primary server does not remove any necessary WAL data in the time between the end of the base backup and the start of streaming replication on the new standby.
The specified replication slot has to exist unless the option `-C` is also used.

If this option is not specified and the server supports temporary replication slots (version 10 and later), then a temporary replication slot is automatically used for WAL streaming.
```
X method
--wal-method=method
```

Includes the required WAL (write-ahead log) files in the backup. This will include all write-ahead logs generated during the backup. Unless the method none is specified, it is possible to start a postmaster in the target directory without the need to consult the log archive, thus making the output a completely standalone backup.
  ```
  -P
  --progress
  ```

Enables progress reporting.
```
-U username
--username=username
Specifies the user name to connect as.
-F format
--format=format
```

Selects the format for the output. format can be one of the following:
```
.p
plain
```

Write the output as plain files, with the same layout as the source server’s data directory and tablespaces. When the cluster has no additional tablespaces, the whole database will be placed in the target directory. If the cluster contains additional tablespaces, the main data directory will be placed in the target directory, but all other tablespaces will be placed in the same absolute path as they have on the source server.
This is the default format.
```
-R
--write-recovery-conf
```
Creates a standby.signal file and appends connection settings to the postgresql.auto.conf file in the target directory (or within the base archive file when using tar format). This eases setting up a standby server using the results of the backup.
The postgresql.auto.conf file will record the connection settings and, if specified, the replication slot that pg_basebackup is using, so that streaming replication will use the same settings later on.
After running this command, you could see there is postgresslave directory in the /tmp/ directory.

## Let’s Setup the Slave !

Now we need to setup another container for the slave, we are going to use the same docker-compose.yml file and just change the container name and the port.

```yaml
version: '3'
services:
  database:
    image: postgres:13
    container_name: my_postgres_database
    restart: always
    volumes:
        - ./data-slave:/var/lib/postgresql/data
        - ./my-postgres.conf:/etc/postgresql/postgresql.conf
       - ./my-pg_hba.conf:/etc/postgresql/pg_hba.conf
    ports:
      - "127.0.0.1:5432:5432"
    environment:
      - 'POSTGRES_PASSWORD=my_password'
      - 'POSTGRES_DB=my_default_database'
```
Now you just need to copy the  `/tmp/postgresslave` directory from master, to data-slave directory in your host machine. Let’s see what’s the final step before firing up the slave container.

## The Trick !
Since you have run the pg_basebackup inside a docker container and also asked for recovery config file, it has created a postgresql.auto.conf file inside the data-slave directory. In this file you should see something like this.

```
# Do not edit this file manually!
# It will be overwritten by the ALTER SYSTEM command.
primary_conninfo = 'user=replicator passfile=''/root/.pgpass'' channel_binding=prefer port=5432 sslmode=prefer sslcompression=0 ssl_min_protocol_version=TLSv1.2 gssencmode=prefer krbsrvname=postgres target_session_attrs=any'
primary_slot_name = 'replication_slot_proxy_slave1'
```

Now you could see the primary_conninfo which tells the slave how should connect to the master, but these configurations are not right. Let’s change the primy_conninfo and pass the correct information for connecting to master.
```
primary_conninfo = 'host=127.0.0.1 port=5432 user=replicator password=my_replicator_password'
```

Also we need to add a restore command which tells slave how to deal with this backup, so add this line as well.

```
restore_command = 'cp /var/lib/postgresql/data/pg_wal/%f "%p"'
```

Now it’s finished, you could fire up the slave container as well.
You could go to slave and run the 
```
$ SELECT * FROM pg_replication_slots;
```
query again.

```
-[ RECORD 1 ]-------+-----------------------------
slot_name           | replication_slot_slave1
plugin              | 
slot_type           | physical
datoid              | 
database            | 
temporary           | f
active              | t
active_pid          | 1332
xmin                | 20800
catalog_xmin        | 
restart_lsn         | 0/105AB6F8
confirmed_flush_lsn | 
wal_status          | reserved
safe_wal_size       |

```

Now you could see the slot is activated. You could also test the replication by creating a dummy table on master and check it on slave