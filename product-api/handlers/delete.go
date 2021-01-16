package handlers

import (
	"net/http"

	"github.com/fahmi1597/microservices-go/product-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// DeleteProduct is a handler to delete a product
func (p *Products) DeleteProduct(resp http.ResponseWriter, req *http.Request) {
	id := p.getProductID(req)
	p.log.Debug("Deleting product", "id", id)

	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.log.Error("Failed to delete product (does not exist) ", "id", id)

		http.Error(resp, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		p.log.Error("Failed to delete product", "id", id)

		resp.WriteHeader(http.StatusInternalServerError)
		http.Error(resp, "Internal server error ", http.StatusInternalServerError)
		return
	}

	p.log.Info("Deleted product", "id", id)

	resp.WriteHeader(http.StatusNoContent)
}
