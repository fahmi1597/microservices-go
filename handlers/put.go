package handlers

import (
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

// UpdateProduct is a handler to update a product
func (p *Products) UpdateProduct(resp http.ResponseWriter, req *http.Request) {
	prod := req.Context().Value(KeyProduct{}).(data.Product)

	err := data.UpdateProduct(prod)
	p.l.Println("[DEBUG] Updating product: ", prod.ID)

	if err == data.ErrProductNotFound {
		http.Error(resp, "Product not found!", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(resp, "Internal server error", http.StatusInternalServerError)
		return
	}
}
