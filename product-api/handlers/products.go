package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fahmi1597/microservices-go/product-api/data"
	"github.com/gorilla/mux"
)

// Products is a handler for getting and updating products
type Products struct {
	l *log.Logger
	v *data.Validation
}

// NewProduct creates a handler where logger is injected
func NewProduct(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}

// KeyProduct is a key of product used in request context
type KeyProduct struct{}

func (p *Products) getProductID(req *http.Request) int {

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.l.Println("[ERROR] Converting product id", id)
		panic(err)
	}

	return id
}

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}
