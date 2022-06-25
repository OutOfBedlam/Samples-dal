
1. install postgres

> reference https://hub.docker.com/_/postgres

```sh
docker run \
  --name my-postgres \
  -e POSTGRES_PASSWORD=my-secret \
  -e POSTGRES_USER=my-user \
  -e POSTGRES_DB=my-db \
  -p 50000:5432 \
  -d \
  postgres
```

2. build

```sh
make
```

3. run

```sh
./tmp/sample-dal -h

Usage: sample-dal <command>

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  check-conn
    connect test

  create-schema
    create tables

  collections
    list collections

  select
    select table

Run "sample-dal <command> --help" for more information on a command.

```