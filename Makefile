.PHONY: deps clean build deploy

deps:
	go get -u ./...

clean: 
	rm -rf ./src/main
	
build:
	GOOS=linux GOARCH=amd64 go build -o src/main ./src

test-unit: build
	go test ./src/

deploy: build
	sam package --template-file ./template/staging.yaml --s3-bucket package-bucket-example --output-template-file packaged.yaml
	sam deploy --template-file packaged.yaml --stack-name sam-cli-example --capabilities CAPABILITY_IAM
