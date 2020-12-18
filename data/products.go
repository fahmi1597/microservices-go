package data

import (
	"fmt"
	"time"
)

// Product is a struct of Product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products is a reference to struct of Product
type Products []*Product

// ErrProductNotFound is an error notification for product that doesn't exist
var ErrProductNotFound = fmt.Errorf("Product not found")

// GetProduct return a product
func GetProduct(id int) (p *Product, err error) {
	prod, _, err := findProduct(id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return prod, nil
}

// GetListProduct return a list of products
func GetListProduct() Products {
	return productList
}

// AddProduct is a function to add the requested product
func AddProduct(p Product) {

	// Get latest item[position] id
	cID := productList[len(productList)-1].ID
	p.ID = cID + 1
	productList = append(productList, &p)
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

// DeleteProduct is a function to delete the requested product
func DeleteProduct(id int) error {
	_, pos, _ := findProduct(id)
	if pos != -1 {
		productList = append(productList[:pos], productList[pos+1])
		return nil
	}

	return ErrProductNotFound
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

// Dummy data
var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frosty milky coffee",
		Price:       2.45,
		SKU:         "abc-123-aa",
		CreatedOn:   time.Now().Local().String(),
		UpdatedOn:   time.Now().Local().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffe without milk",
		Price:       1.29,
		SKU:         "asd-321-bb",
		CreatedOn:   time.Now().Local().String(),
		UpdatedOn:   time.Now().Local().String(),
	},
}
