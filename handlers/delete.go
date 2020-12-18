package handlers

import (
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

// DeleteProduct is a handler to delete a product
func (p *Products) DeleteProduct(resp http.ResponseWriter, req *http.Request) {

	id := p.getProductID(req)
	p.l.Println("[DEBUG] Retrieve product id", id)

	err := data.DeleteProduct(id)
	if err != data.ErrProductNotFound {
		p.l.Println("[ERROR] deleting record id does not exist")
		http.Error(resp, "Product not found!", http.StatusNotFound)
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		p.l.Println("[ERROR] deleting record")
		resp.WriteHeader(http.StatusInternalServerError)
		http.Error(resp, "Internal server error ", http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusNoContent)
}
