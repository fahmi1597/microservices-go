package handlers

import (
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

// swagger:route POST /products products addProduct
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse

// AddProduct is a handler for creating a product
func (p *Products) AddProduct(resp http.ResponseWriter, req *http.Request) {

	prod := req.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("[DEBUG] Inserting product: %#v\n", prod)
	data.AddProduct(prod)

}
