package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/fahmi1597/microservices-go/data"
	"github.com/gorilla/mux"
)

// Products struct implements the http.handler
type Products struct {
	l *log.Logger
}

// NewProduct creates a handler where logger is injected
func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

// GetProduct is a handler that return list of products
func (p *Products) GetProduct(resp http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle GET requests")

	// Retrieve products
	lp := data.GetProducts()

	// Serialize products to JSON
	err := lp.ToJSON(resp)
	if err != nil {
		p.l.Println("Failed to encode JSON", http.StatusInternalServerError)
		http.Error(resp, "Unable to encode data to json", http.StatusInternalServerError)
		return
	}
}

// AddProduct is a handler to adds a product
func (p *Products) AddProduct(resp http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle POST requests")

	ap := req.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&ap)
}

// UpdateProduct is a handler to update a product
func (p *Products) UpdateProduct(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(resp, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT requests", id)

	up := req.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &up)
	if err == data.ErrProductNotFound {
		http.Error(resp, "Product not found!", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(resp, "Internal server error ", http.StatusInternalServerError)
		return
	}
}

// KeyProduct is a key of product
type KeyProduct struct{}

// MiddlewareValidation is a midleware handler for validation
func (p *Products) MiddlewareValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		prod := data.Product{}

		err := prod.FromJSON(req.Body)

		if err != nil {
			p.l.Println("Error: deserialize product", http.StatusBadRequest)
			http.Error(resp, "Error reading product", http.StatusBadRequest)
			return
		}

		// Create context of validated product
		ctx := context.WithValue(req.Context(), KeyProduct{}, prod)
		// Send the validated product request to the next handler
		req = req.WithContext(ctx)
		next.ServeHTTP(resp, req)
	})
}
