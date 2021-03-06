module github.com/fahmi1597/microservices-go/product-api

go 1.15

replace github.com/fahmi1597/microservices-go/currency => /home/iamf/go/src/github.com/fahmi1597/microservices-go/currency

require (
	github.com/fahmi1597/microservices-go/currency v0.0.0-20210116052656-0f24482c0413
	github.com/go-openapi/errors v0.19.9
	github.com/go-openapi/runtime v0.19.24
	github.com/go-openapi/strfmt v0.19.11
	github.com/go-openapi/swag v0.19.12
	github.com/go-openapi/validate v0.20.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-hclog v0.15.0
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.34.0
)
