package handlers

import (
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

// swagger:route PUT /products products updateProduct
// Update a product details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// UpdateProduct is a handler for updating a product
func (p *Products) UpdateProduct(resp http.ResponseWriter, req *http.Request) {

	prod := req.Context().Value(KeyProduct{}).(data.Product)
	err := data.UpdateProduct(prod)
	p.l.Println("[DEBUG] Updating product id: ", prod.ID)

	if err == data.ErrProductNotFound {
		resp.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, resp)
		return
	}
	if err != nil {
		http.Error(resp, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}
