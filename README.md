
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

