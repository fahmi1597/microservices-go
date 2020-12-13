package handlers

import (
	"log"
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

// Products struct implements the http.handler
type Products struct {
	l *log.Logger
}

// NewProduct creates a handler where logger is injected
func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

// Satisfy the ServeHTTP from http.handler interface
func (p *Products) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	// Handle the GET method
	if req.Method == http.MethodGet {
		p.getProduct(resp, req)
		return
	}

	// Handle the others method and write log
	resp.WriteHeader(http.StatusMethodNotAllowed)
	p.l.Println("Bad method from", req.RemoteAddr)
}

func (p *Products) getProduct(resp http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle products requests")

	// Retrieve products
	lp := data.GetProducts()

	// Serialize products to JSON
	err := lp.ToJSON(resp)
	if err != nil {
		http.Error(resp, "Unable to produce json", http.StatusInternalServerError)
		p.l.Println("Failed to response products requests", http.StatusInternalServerError)
	}
}
