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

### 記事の内容を考える

- package構成
- tddで開発した話
- 内容(gzipを取り出して編集してgzipにしてupload)
- go testのTestMainの話
- templateの話(環境変数、http://localstack:4572とか)
- makefileの話


以下は公式で生成されたもの

# sam-app

This is a sample template for sam-app - Below is a brief explanation of what we have generated for you:

```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── hello-world                 <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda function code
│   └── main_test.go            <-- Unit tests
└── template.yaml
```