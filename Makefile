vendor:
	go mod vendor

bin:
	mkdir bin
run:
	go run *.go -mod=vendor

build: bin
	go build -mod=vendor -o ./bin/

