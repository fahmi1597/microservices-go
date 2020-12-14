package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	// Handle the GET request
	if req.Method == http.MethodGet {
		p.getProduct(resp, req)
		return
	}

	// Handle the POST request
	if req.Method == http.MethodPost {
		p.addProduct(resp, req)
		return
	}

	// Handle the PUT request
	if req.Method == http.MethodPut {

		// Use regex to construct the pattern we want to get
		re := regexp.MustCompile("/([0-9]+)")
		// Store match string to [][]string
		m := re.FindAllStringSubmatch(req.URL.Path, -1)

		// Filter URI from m
		// m = [[/id id]], so len(m) must be 1
		if len(m) != 1 {
			p.l.Println("Invalid request, found more than one id:", m)
			http.Error(resp, "Invalid URI", http.StatusBadRequest)
			return
		}

		//  m[0] = [/id id], so len must be 2
		if len(m[0]) != 2 {
			p.l.Println("Bad URI: ", m[0])
			http.Error(resp, "Bad Request", http.StatusBadRequest)
			return
		}

		// Get id
		idStr := m[0][1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			p.l.Println("Unable to convert URI to number", idStr)
			http.Error(resp, "Error", http.StatusInternalServerError)
			return
		}

		p.updateProduct(id, resp, req)

		return
	}

	// Handle the others method and write log
	resp.WriteHeader(http.StatusMethodNotAllowed)
	p.l.Println("Bad method from", req.RemoteAddr)
}

func (p *Products) getProduct(resp http.ResponseWriter, req *http.Request) {
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

	p.l.Println("Handle GET requests: success")
}

func (p *Products) addProduct(resp http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle POST requests")

	ap := &data.Product{}

	err := ap.FromJSON(req.Body)
	if err != nil {
		p.l.Println("Failed to decode JSON", http.StatusBadRequest)
		http.Error(resp, "Unable to decode data to json", http.StatusBadRequest)
		return
	}

	// p.l.Printf("data : %#v", ap)
	data.AddProduct(ap)

	p.l.Println("Handle POST requests: success")
}

func (p *Products) updateProduct(id int, resp http.ResponseWriter, req *http.Request) {

	p.l.Println("Handle PUT requests")

	up := &data.Product{}

	if err := up.FromJSON(req.Body); err != nil {
		p.l.Println("Failed to decode JSON", http.StatusBadRequest)
		http.Error(resp, "Unable to decode data to json", http.StatusBadRequest)
		return
	}

	err := data.UpdateProduct(id, up)
	if err == data.ErrProductNotFound {
		http.Error(resp, "Product not found!", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(resp, "Internal server error ", http.StatusInternalServerError)
		return
	}

	p.l.Println("Handle PUT requests: success")
}
