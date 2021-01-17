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
$ protoc -I protos/ protos/currency.proto --go_out=protos/currency --go-grpc_out=protos/currency --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative
```

## Interact with Server using gRPCurl

https://github.com/fullstorydev/grpcurl

```
$ go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

### List services
Request:
```
$ grpcurl.exe --plaintext localhost:9002 list
```

Response:
```
currency.Currency
grpc.reflection.v1alpha.ServerReflection
```

### List methods
Request:
```
$ grpcurl.exe --plaintext localhost:9002 list currency.Currency
```

Response:
```
currency.Currency.GetRate
currency.Currency.SubscribeRates
```

### Describe GetRate method
Request:
```
$ grpcurl.exe --plaintext localhost:9002 describe currency.Currency.GetRate
```

Response:
```
currency.Currency.GetRate is a method:
rpc GetRate ( .currency.RateRequest ) returns ( .currency.RateResponse );
```

### Describe RateRequest detail
Request:
```
$ grpcurl.exe --plaintext --msg-template localhost:9002 describe currency.RateRequest
```

Response:
```
currency.RateRequest is a message:
message RateRequest {
  .currency.Currencies Base = 1 [json_name = "base"];
  .currency.Currencies Destination = 2 [json_name = "destination"];
}

Message template:
{
  "base": "EUR",
  "destination": "EUR"
}
```

### Send a request to GetRate method
Request:
```
$ grpcurl.exe --plaintext -d '{"base":"EUR", "destination":"JPY"}' localhost:9002 currency.Currency.GetRate
```

Response
```
{
  "rate": 126.69172026319903
}
```

### Send request to SubscripeRates method
Request:
```
$ grpcurl.exe --plaintext -d @ localhost:9002 currency.Currency.SubscribeRates

{ "base" : "JPY", "destination" : "IDR"}
```

Response:
```
{
  "base": "JPY",
  "destination": "IDR",
  "rate": 125.58273518485274
}
{
  "base": "JPY",
  "destination": "IDR",
  "rate": 128.32494481468817
}
...
```