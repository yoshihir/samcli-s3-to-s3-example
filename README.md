### 初期設定

```shell
brew tap aws/tap
brew install aws-sam-cli
```

### docker

```shell
env TMPDIR=/private$TMPDIR docker-compose up -d
```

### unit test

```shell
cd src
go test
```

### build
```shell
make build
```

### integration-testing
```shell
make integration-testing
```


### validate template
```shell
sam validate --template ./template/local.yaml
```


### deploy
ソースを置くbucketは別途作成要(今回はpackage-bucket-example)
```shell
make deploy
```

### You should know

```shell
# sample code create
sam init --runtime go
# create event_file.json
sam local generate-event s3 put > event_file.json
# delete sam cli
aws cloudformation delete-stack --stack-name test-sam-cli
```