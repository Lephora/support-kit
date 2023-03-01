# support-kit
Support kit for lephora dev, test and CI/CD

1. openapi-validator

supporting openapi schema validate

进入openapi-validator目录，执行

```shell
# linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o validate main.go

#mac
go build -o validate main.go
```


2. acceptance-test

supporting E2E test and contract test

进入acceptance-test/core目录，执行

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o acceptance-test main.go

#mac
go build -o acceptance-test main.go
```
