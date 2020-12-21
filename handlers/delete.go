package handlers

import (
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

// swagger:route DELETE /products/{id} products DeleteProduct
// Delete a product
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// DeleteProduct is a handler to delete a product
func (p *Products) DeleteProduct(resp http.ResponseWriter, req *http.Request) {

	id := p.getProductID(req)
	p.l.Println("[DEBUG] Delete product id:", id)

	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] Deleting record id", id, "(does not exist)")
		http.Error(resp, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		p.l.Println("[ERROR] Deleting record")
		resp.WriteHeader(http.StatusInternalServerError)
		http.Error(resp, "Internal server error ", http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusNoContent)
}
