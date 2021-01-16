## Request method test
get:
	curl -s localhost:9000/products | jq
get_1:
	curl -s localhost:9000/products/1 | jq
post:
	curl -s localhost:9000/products -d '{"name":"tea", "description":"A cup of tea", "price":0.25, "sku":"aaa-111-aaa"}'
put:
	curl -sX PUT localhost:9000/products -d '{"id":1,"name":"milk", "description":"A cup of milk", "price":0.25, "sku":"bbb-222-bbb"}'
delete:
	curl -sX DELETE localhost:9000/products/1

## Swagger
swagger:
	@echo Ensure you have the swagger CLI or this command will fail.
	@echo You can install the swagger CLI with: go get -u github.com/go-swagger/go-swagger/cmd/swagger
	@echo ....
	swagger generate spec -o ./swagger.yaml --scan-models
	