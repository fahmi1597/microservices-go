package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	p.l.Println("handle products requests")

	lp := data.GetProducts()
	d, err := json.Marshal(lp)
	if err != nil {
		http.Error(resp, "Unable to produce json ", http.StatusInternalServerError)
		p.l.Println("Failed to response products requests", http.StatusInternalServerError)
	}

	resp.Write(d)
}
