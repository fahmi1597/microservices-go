package handlers

import (
	"net/http"

	"github.com/fahmi1597/microservices-go/product-api/data"
)

// swagger:route GET /products products getProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// GetProductList is a handler that return list of products
func (p *Products) GetProductList(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	cur := req.URL.Query().Get("currency")
	//p.log.Debug("Retrieving product list", "currency", cur)

	// Retrieve product list
	prods, err := p.productDB.GetProductList(cur)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, resp)
		return
	}

	// Serialize products to JSON
	err = data.ToJSON(&prods, resp)
	if err != nil {
		p.log.Error("Failed to serialize data", "error", err)
		return
	}
}

// swagger:route GET /products/{id} products getProduct
// Return a product from the database
// responses:
//	200: productResponse
//	404: errorResponse

// GetProductByID is a handler that return a product
func (p *Products) GetProductByID(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	id := p.getProductID(req)
	cur := req.URL.Query().Get("currency")
	//p.log.Debug("Retrieving product", "id", id, "currency", cur)

	prod, err := p.productDB.GetProductByID(id, cur)
	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.log.Error("Failed to retrieve product id", "id", id, "error", err)

		resp.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, resp)
		return
	default:
		p.log.Error("Failed to retrieve product id", "id", id, "error", err)

		resp.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, resp)
		return
	}

	// Serialize product to JSON
	err = data.ToJSON(&prod, resp)
	if err != nil {
		p.log.Error("Failed to serialize data", "error", err)
		return
	}

}
