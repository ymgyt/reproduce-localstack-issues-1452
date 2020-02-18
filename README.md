# Steps

https://github.com/localstack/localstack/issues/1452


```console
# start docker containers
docker-compose up --detach

# run go
go run main.go
```

check local bucket

```console
AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy aws --endpoint=http://localhost:4572 s3 ls ymgyt-localstack-repro --recursive
```

