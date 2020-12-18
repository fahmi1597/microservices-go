package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fahmi1597/microservices-go/data"
	"github.com/gorilla/mux"
)

// Products struct implements the http.handler
type Products struct {
	l *log.Logger
	v *data.Validation
}

// NewProduct creates a handler where logger is injected
func NewProduct(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}

// KeyProduct is a key of product
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
