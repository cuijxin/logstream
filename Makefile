build: server proxy client

client:
	go build -o bin/client ./cmd/client

server:
	go build -o bin/server ./cmd/server

proxy:
	go build -o bin/proxy ./cmd/proxy

clean:
	rm -rf ./bin
install:
	go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

generate:
	protoc pkg/proto/rpc.proto --go_out=plugins=grpc:.
