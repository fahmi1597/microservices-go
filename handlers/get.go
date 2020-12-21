package handlers

import (
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

// swagger:route GET /products products GetProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// GetProducts is a handler that return list of products
func (p *Products) GetProducts(resp http.ResponseWriter, req *http.Request) {
	p.l.Println("[DEBUG] Retrieve product list")

	// Retrieve products
	lprod := data.GetListProduct()

	// Serialize products to JSON
	err := data.ToJSON(lprod, resp)
	if err != nil {
		p.l.Println("[ERROR] Failed to serialize data", http.StatusInternalServerError)
		http.Error(resp, "Error reading products", http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /products/{id} products GetProduct
// Return a product from the database
// responses:
//	200: productResponse
//	404: errorResponse

// GetProduct is a handler that return a product
func (p *Products) GetProduct(resp http.ResponseWriter, req *http.Request) {

	id := p.getProductID(req)
	p.l.Println("[DEBUG] Retrieve product id:", id)

	prod, err := data.GetProduct(id)
	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Println("[ERROR] Retrieving product", err)

		resp.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, resp)
		return
	default:
		p.l.Println("[ERROR] Retrieving product", err)

		resp.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, resp)
		return
	}

	// Serialize product to JSON
	err = data.ToJSON(prod, resp)
	if err != nil {
		p.l.Println("[ERROR] Failed to serialize data", http.StatusInternalServerError)
		http.Error(resp, "Error reading product", http.StatusInternalServerError)
		return
	}
}
