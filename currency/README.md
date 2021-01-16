# Currency Services

The currency service is a gRPC service which provides up to date exchange rates and currency conversion capabilities.

## Tools

### Protocol buffers

https://github.com/protocolbuffers/protobuf/releases/tag/v3.14.0 

- Windows - protoc-3.14.0-win64.zip
- Linux - protoc-3.14.0-linux-x86_64.zip

### gRPC plugins

https://grpc.io/docs/languages/go/quickstart/

```
$ export GO111MODULE=on  # Enable module mode
$ go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

## Generate Go Code

```
protoc -I protos/ protos/currency.proto --go_out=protos/currency --go-grpc_out=protos/currency --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative
```

## Interact with Server using gRPCurl

https://github.com/fullstorydev/grpcurl

```
go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

### List services

```
grpcurl.exe --plaintext localhost:9002 list
```

### List methods

```
grpcurl.exe --plaintext localhost:9002 list currency.Currency
```

### Describe GetRate method

```
grpcurl.exe --plaintext localhost:9002 describe currency.Currency.GetRate
```
### Describe RateRequest detail

```
grpcurl.exe --plaintext localhost:9002 describe currency.RateRequest
```

### Send a request to a method

```
grpcurl.exe --plaintext -d '{"base":"EUR", "destination":"JPY"}' localhost:9002 currency.Currency.GetRate
```