package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fahmi1597/microservices-go/product-api/data"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

// Products is a handler for getting and updating products
type Products struct {
	log       hclog.Logger
	validator *data.Validation
	productDB *data.ProductsDB
}

// KeyProduct is a key of product used in request context
type KeyProduct struct{}

// ErrInvalidProductPath is an error message when the product path is not valid
// reserved for 404 ?
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// NewProductHandler creates a new Products handler
func NewProductHandler(l hclog.Logger, v *data.Validation, pdb *data.ProductsDB) *Products {
	return &Products{l, v, pdb}
}

func (p *Products) getProductID(req *http.Request) int {

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.log.Error("Failed to convert product id", "id", id)
		panic(err)
	}

	return id
}
