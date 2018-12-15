PROJECT_NAME:= "localstack-example"

.PHONY: deps clean build integration-testing deploy

deps:
	go get -u ./...

clean: 
	rm -rf ./src/main
	
build:
	GOOS=linux GOARCH=amd64 go build -o src/main ./src

integration-testing: build
	docker-compose -p $(PROJECT_NAME) down
	env TMPDIR=/private$TMPDIR docker-compose -p $(PROJECT_NAME) up -d
	sleep 10s
	aws --endpoint-url=http://localhost:4572 s3 mb s3://bucket-example
	aws --endpoint-url=http://localhost:4572 s3 mb s3://bucket-example-convert
	aws --endpoint-url=http://localhost:4572 s3 cp ./testdata/example.json.gz s3://bucket-example/example.json.gz
	sam local invoke MainFunction --event event_file.json --template ./template/local.yaml \
	--docker-network $$(docker network ls -q -f name=$(PROJECT_NAME))

deploy: build
	sam package --template-file ./template/staging.yaml --s3-bucket package-bucket-example --output-template-file packaged.yaml
	sam deploy --template-file packaged.yaml --stack-name sam-cli-example --capabilities CAPABILITY_IAM