# Replicating Postgres inside Docker — The How To

[Original Artical On Medium](https://medium.com/@2hamed/replicating-postgres-inside-docker-the-how-to-3244dc2305be)

>Let me tell you this. Docker is fascinating. Not only that it is a great development tool, it is also perfect for production.
In this post I’m going to explain how I set up a postgresql replication using docker.

Replication is a crucial part of every modern infrastructure these days and with scalability in mind it plays an important role on delivering high quality software to billions of people across the globe.
PostgreSql with built-in replication support is a great choice for a reliable and scalable database engine. The setup that we’re about to discuss consists of 1 master instance which handles both read and write operations and many slave instance which will only serve read requests and that is because writing data is whole different story than reading them.
Talk is cheap, show me the code…
We’ll start by creating the Master instance. This is what goes into the Dockerfile.

```docker
FROM postgres:9.6-alpine
COPY ./setup-master.sh /docker-entrypoint-initdb.d/setup-master.sh
RUN chmod 0666 /docker-entrypoint-initdb.d/setup-master.sh
```

You’ll notice that there is a setup-master.sh file which needs to be copied to that image which makes the Postgres ready for being a master in the replication process.

```bash
#!/bin/bash
echo "host replication all 0.0.0.0/0 md5" >> "$PGDATA/pg_hba.conf"
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
CREATE USER $PG_REP_USER REPLICATION LOGIN CONNECTION LIMIT 100 ENCRYPTED PASSWORD '$PG_REP_PASSWORD';
EOSQL
cat >> ${PGDATA}/postgresql.conf <<EOF
wal_level = hot_standby
archive_mode = on
archive_command = 'cd .'
max_wal_senders = 8
wal_keep_segments = 8
hot_standby = on
EOF
```

Don’t worry about all those variables there. We’ll discuss them later on.
Now we need a Dockerfile for slave instances as well.

```docker
FROM postgres:9.6-alpine
ENV GOSU_VERSION 1.10
ADD ./gosu /usr/bin/
RUN chmod +x /usr/bin/gosu
RUN apk add --update iputils
COPY ./docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["gosu", "postgres", "postgres"]
```

1. We need the gosu binary to execute the postgres as root
2. We also need the iputils package to be able to ping the master. You’ll see in a bit.

The next step is to prepare the slave images to actually be slaves. We use a `docker-entrypoin.sh` file which will be the first thing that docker execute upon creating the container.

```bash
#!/bin/bash
if [ ! -s "$PGDATA/PG_VERSION" ]; then
echo "*:*:*:$PG_REP_USER:$PG_REP_PASSWORD" > ~/.pgpass
chmod 0600 ~/.pgpass
until ping -c 1 -W 1 pg_master
do
echo "Waiting for master to ping..."
sleep 1s
done
until pg_basebackup -h pg_master -D ${PGDATA} -U ${PG_REP_USER} -vP -W
do
echo "Waiting for master to connect..."
sleep 1s
done
echo "host replication all 0.0.0.0/0 md5" >> "$PGDATA/pg_hba.conf"
set -e
cat > ${PGDATA}/recovery.conf <<EOF
standby_mode = on
primary_conninfo = 'host=pg_master port=5432 user=$PG_REP_USER password=$PG_REP_PASSWORD'
trigger_file = '/tmp/touch_me_to_promote_to_me_master'
EOF
chown postgres. ${PGDATA} -R
chmod 700 ${PGDATA} -R
fi
sed -i 's/wal_level = hot_standby/wal_level = replica/g' ${PGDATA}/postgresql.conf
exec "$@"
```

1. On the first line we check that this instance has already been set up or not by checking the   `PG_VERSION` file in the `PG_DATA` path so as not to do it on every startup of the container.
2. We put the replication user and password in the `.pgpass` so postgres can access it.
3. We start pinging the Master to make sure that it’s already up and running.
4. And we put in the necessary configuration in place for the slave servers.

Ok, it’s almost done.
Now we need only to create a docker-compose.yml file to start our database containers.

```yaml
version: "3"
services:
 pg_master:
  build: ./master
  volumes:
   - pg_data:/var/lib/postgresql/data
  environment:
   - POSTGRES_USER=hamed
   - POSTGRES_PASSWORD=123456
   - POSTGRES_DB=hamed
   - PG_REP_USER=rep
   - PG_REP_PASSWORD=123456
  
  networks:
   default:
   aliases:
    - pg_cluster
 pg_slave:
  build: ./slave
  environment:
   - POSTGRES_USER=hamed
   - POSTGRES_PASSWORD=123456
   - POSTGRES_DB=hamed
   - PG_REP_USER=rep
   - PG_REP_PASSWORD=123456
  networks:
   default:
   aliases:
    - pg_cluster
volumes:
 pg_data:
```

Keep in mind that we have put the Dockerfile and `setup-master.sh` file the master directory and the slave’s Dockerfile and `docker-entrypoint.sh `and the gosu binary inside the slave directory.
We now have 2 ways of running the this setup.

1. Using docker compose:
    ```
    docker-compose up
    ```
2. Using docker swarm:

To be able to run these container in a swarm we need to build those images and put on a public docker repository which I have already done and uploaded to Docker Hub.

1. [https://hub.docker.com/r/2hamed/pg_master/](https://hub.docker.com/r/2hamed/pg_master/)
2. [https://hub.docker.com/r/2hamed/pg_slave/](https://hub.docker.com/r/2hamed/pg_slave/)

Just change the line build: `./master` and build: `./slave` to image: `2hamed/pg_master` and image: `2hamed/pg_slave` respectively and run:

```
$ docker stack deploy -c docker-compose.yml my_pg_replication
```

This way you can easily scale up your slaves to as many as you need but keep in mind that you need to keep `max_wal_senders = 8` in your master up to date.

**Bonus**: If you have noticed we have set 2 network aliases `pg_master` and `pg_cluster` . In your application you can do the write operations on `pg_master` and the read operations on `pg_cluster` and docker will automatically route your requests to the correct container instance.

All the code in this post is available at [https://github.com/2hamed/docker-pg-replication](https://github.com/2hamed/docker-pg-replication).