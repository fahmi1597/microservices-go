package data

import (
	"fmt"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"`

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+{3}-[0-9]+{3}-[a-z]+{2}
	SKU string `json:"sku" validate:"required,sku"`
}

// Products is a reference to struct of Product
type Products []*Product

// ErrProductNotFound is an error message when client requesting a product that can not be found or doesn't exist
var ErrProductNotFound = fmt.Errorf("Product not found")

// GetListProduct return a list of products
func GetListProduct() Products {
	return productList
}

// GetProduct return a product
func GetProduct(id int) (*Product, error) {

	i := findProductByID(id)
	if i != -1 {
		return productList[i], nil
	}

	return nil, ErrProductNotFound
}

// AddProduct is a function to add the requested product
func AddProduct(p Product) {
	// Get latest item[position] id
	cID := productList[len(productList)-1].ID
	p.ID = cID + 1
	productList = append(productList, &p)
}

// UpdateProduct is a function to update the requested product
func UpdateProduct(p Product) error {

	i := findProductByID(p.ID)
	if i != -1 {
		productList[i] = &p
		return nil
	}

	return ErrProductNotFound
}

// DeleteProduct is a function to delete the requested product
func DeleteProduct(id int) error {
	i := findProductByID(id)
	if i != -1 {
		productList = append(productList[:i], productList[i+1:]...)
		return nil
	}

	return ErrProductNotFound
}

// findProductByID finds the index of a product in the database
// returns -1 when no product can be found
func findProductByID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

// Dummy data
var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frosty milky coffee",
		Price:       2.45,
		SKU:         "abc-123-aa",
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffe without milk",
		Price:       1.29,
		SKU:         "asd-321-bb",
	},
}
