package handlers

import (
	"net/http"
	"strconv"

	"github.com/fahmi1597/microservices-go/data"
	"github.com/gorilla/mux"
)

// UpdateProduct is a handler to update a product
func (p *Products) UpdateProduct(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(resp, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT requests", id)

	prod := req.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(resp, "Product not found!", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(resp, "Internal server error ", http.StatusInternalServerError)
		return
	}
}
