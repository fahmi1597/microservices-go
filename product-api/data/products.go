package data

import (
	"context"
	"fmt"

	"github.com/fahmi1597/microservices-go/currency/protos/currency"
	protogc "github.com/fahmi1597/microservices-go/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

// ErrProductNotFound is an error message when client requesting a product that can not be found or doesn't exist
var ErrProductNotFound = fmt.Errorf("Product not found")

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
	Price float64 `json:"price" validate:"gt=0,required"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[0-9]+-[a-z]+
	SKU string `json:"sku" validate:"sku,required"`
}

// Products is a reference to struct of Product
type Products []*Product

// ProductsDB is used to query data with
type ProductsDB struct {
	log      hclog.Logger
	currency currency.CurrencyClient
}

// NewProductsDB do something
func NewProductsDB(log hclog.Logger, cc currency.CurrencyClient) *ProductsDB {
	return &ProductsDB{
		log:      log,
		currency: cc,
	}
}

// GetProductList do something
func (p *ProductsDB) GetProductList(currency string) (Products, error) {

	// Just return list product if there's no currency in request
	if currency == "" {
		return productList, nil
	}

	rateRatio, err := p.getRateRatio(currency)
	if err != nil {
		p.log.Error("Unable to get exchange rate", "currency", currency, "error", err)
		return nil, err
	}
	// Create a new slice of Product to mutate the product price with requested currency rate
	newListProduct := Products{}
	for _, product := range productList {
		tempProducts := *product
		tempProducts.Price = tempProducts.Price * rateRatio

		// Append the mutated Products to newProducts
		newListProduct = append(newListProduct, &tempProducts)
	}

	return newListProduct, nil

}

// GetProductByID return a product
func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findProductByID(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}

	if currency == "" {
		return productList[i], nil
	}

	rateRatio, err := p.getRateRatio(currency)
	if err != nil {
		p.log.Error("Unable to get exchange rate", "currency", currency, "error", err)
		return nil, err
	}

	newProduct := *productList[i]
	newProduct.Price = newProduct.Price * rateRatio

	return &newProduct, nil

}

// AddProduct is a function to add the requested product
func AddProduct(product Product) {

	// Get latest item[position] id
	id := productList[len(productList)-1].ID
	product.ID = id + 1
	productList = append(productList, &product)

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

func (p *ProductsDB) getRateRatio(dest string) (float64, error) {

	// Base currency always USD
	rateRequest := &protogc.RateRequest{
		Base:        protogc.Currencies(protogc.Currencies_value["EUR"]),
		Destination: protogc.Currencies(protogc.Currencies_value[dest]),
	}

	// Get exchange rate
	rateResponse, err := p.currency.GetRate(context.Background(), rateRequest)
	if err != nil {
		p.log.Error("Unable to get rate ratio", "error", err.Error())
		return 0, nil
	}

	return rateResponse.Rate, nil
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
