package main

import (
	"fmt"
	"testing"

	"github.com/fahmi1597/microservices-go/ms-client/client"
	"github.com/fahmi1597/microservices-go/ms-client/client/products"
)

func TestClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:3000")
	c := client.NewHTTPClientWithConfig(nil, cfg)
	p := products.NewGetProductsParams()
	resp, err := c.Products.GetProducts(p)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", resp.GetPayload()[0])
	t.Fail()
}
