package handlers

import (
	"net/http"

	"github.com/fahmi1597/microservices-go/product-api/data"
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
	resp.Header().Set("Content-Type", "application/json")

	prod := req.Context().Value(KeyProduct{}).(data.Product)
	p.log.Debug("Updating product", "id", prod.ID)

	err := data.UpdateProduct(prod)
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
