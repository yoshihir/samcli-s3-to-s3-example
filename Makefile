.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./main/main
	
build:
	GOOS=linux GOARCH=amd64 go build -o src/main ./src