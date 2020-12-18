package handlers

import (
	"fmt"
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

// GetListProduct is a handler that return list of products
func (p *Products) GetListProduct(resp http.ResponseWriter, req *http.Request) {
	p.l.Println("[DEBUG] Retrieve product list")

	// Retrieve products
	lprod := data.GetListProduct()

	// Serialize products to JSON
	err := data.ToJSON(lprod, resp)
	if err != nil {
		p.l.Println("[ERROR] Failed to serialize data", http.StatusInternalServerError)
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(resp, http.StatusInternalServerError)
		return
	}
}

// GetProduct is a handler that return a product
func (p *Products) GetProduct(resp http.ResponseWriter, req *http.Request) {
	id := p.getProductID(req)

	p.l.Println("[DEBUG] Retrieve product", id)

	// Retrieve product
	prod, err := data.GetProduct(id)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprint(resp, http.StatusNotFound)
		return
	}
	// Serialize product to JSON
	err = data.ToJSON(prod, resp)
	if err != nil {
		p.l.Println("[ERROR] Failed to serialize data", http.StatusInternalServerError)
		resp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(resp, http.StatusInternalServerError)
		return
	}
}
