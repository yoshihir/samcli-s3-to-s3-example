.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./src/main
	
build:
	GOOS=linux GOARCH=amd64 go build -o src/main ./src

test-unit: build
	go test ./src/