package handlers

import (
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

// AddProduct is a handler to adds a product
func (p *Products) AddProduct(resp http.ResponseWriter, req *http.Request) {

	prod := req.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("[DEBUG] Inserting product: %#v\n", prod)
	data.AddProduct(prod)

}
