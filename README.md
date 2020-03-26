# logstream

Here is an example of streaming data from a server through a proxy to a client all using [gRPC](https://www.grpc.io/)

The server produces log messages. The proxy creates a stream to the server. The client creates a stream to the proxy and consumes data originally sent from the server.

### Running the system
* * *

To compile all processes, simply run:
```
make
```

The default values of each process will connect the client to the proxy, and the proxy to the server.

To run the system:
```
./bin/server
./bin/proxy
./bin/client
```

It is sometimes helpful to remove the proxy from the system. To do so, run:
```
./bin/server -addr :8000
./bin/client -raddr :8000
```

### Setting up for Development
* * *

Aside from an installation of [Go](https://golang.org/), you will also need three things:

1. the protoc compiler, Downland from the [releases page](https://github.com/protocolbuffers/protobuf/releases).

2. the protoc plugin for Go which includes support for gRPC, and
```    
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

3. the gRPC library in Go.
```
go get google.golang.org/grpc
```
