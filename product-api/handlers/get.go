package handlers

import (
	"context"
	"net/http"

	protogc "github.com/fahmi1597/microservices-go/currency/protos/currency"
	"github.com/fahmi1597/microservices-go/product-api/data"
)

// swagger:route GET /products products getProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// GetProducts is a handler that return list of products
func (p *Products) GetProducts(resp http.ResponseWriter, req *http.Request) {
	p.l.Println("[INFO] Retrieve product list")
	resp.Header().Set("Content-Type", "application/json")

	// Retrieve products
	lprod := data.GetListProduct()

	// Serialize products to JSON
	err := data.ToJSON(lprod, resp)
	if err != nil {
		p.l.Println("[ERROR] Failed to serialize data", http.StatusInternalServerError)
		http.Error(resp, "Error fetching products", http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /products/{id} products getProduct
// Return a product from the database
// responses:
//	200: productResponse
//	404: errorResponse

// GetProduct is a handler that return a product
func (p *Products) GetProduct(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "application/json")
	id := p.getProductID(req)
	p.l.Println("[INFO] Retrieve product id:", id)

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
	// Dummy rate request
	rr := &protogc.RateRequest{
		Base:        protogc.Currencies(protogc.Currencies_value["IDR"]).String(),
		Destination: protogc.Currencies(protogc.Currencies_value["USD"]).String(),
	}
	// Get exchange rate
	excRate, err := p.cc.GetRate(context.Background(), rr)
	if err != nil {
		p.l.Println("[ERROR] Getting exchange rate")
		data.ToJSON(&GenericError{Message: err.Error()}, resp)
	}
	prod.Price = prod.Price * excRate.Rate

	// Serialize product to JSON
	err = data.ToJSON(prod, resp)
	if err != nil {
		p.l.Println("[ERROR] Failed to serialize data", http.StatusInternalServerError)
		http.Error(resp, "Error reading product", http.StatusInternalServerError)
		return
	}
}
