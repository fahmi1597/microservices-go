.PHONY: protos

# Newer approach, grpc plugin generated file are separated, currency.pb.go and currency_grpc.pb.go 
# This one need a small modification in server/currency.go
protos:
	protoc -I protos/ protos/currency.proto --go_out=protos/currency --go-grpc_out=protos/currency --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative

# protoc -I protos/ protos/currency.proto --go_out=protos/currency --go-grpc_out=protos/currency \
# --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative

# 
# protos_plugin:
#	protoc -I protos/ protos/currency.proto --go_out=protos/currency --plugin=protoc-gen-go-grpc