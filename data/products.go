package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Product is a struct of Product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products is a reference to struct of Product
type Products []*Product

// ErrProductNotFound is an error notification for product that doesn't exist
var ErrProductNotFound = fmt.Errorf("Product not found")

// GetProducts return a list of products
func GetProducts() Products {
	return productList
}

// AddProduct is a function to add the requested product
func AddProduct(p *Product) {
	p.ID = productNextID()
	productList = append(productList, p)
}

// UpdateProduct is a function to update the requested product
func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p

	return nil
}

func findProduct(id int) (p *Product, pos int, err error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func productNextID() int {
	cid := productList[len(productList)-1]
	return cid.ID + 1
}

// ToJSON used for converting the struct of Product
// to JSON. It has better performance than json.marshal
// since it doesn't have to buffer the output into memory
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// FromJSON used for converting JSON
// formatted data to struct of Products which refer to Product
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// List of Product is hard coded for testing purpose only
var productList = Products{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frosty milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedOn:   time.Now().Local().String(),
		UpdatedOn:   time.Now().Local().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffe without milk",
		Price:       1.29,
		SKU:         "asd321",
		CreatedOn:   time.Now().Local().String(),
		UpdatedOn:   time.Now().Local().String(),
	},
}
