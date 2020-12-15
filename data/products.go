package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
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

// Validate is a struct validator
func (p *Product) Validate() error {

	v := validator.New()
	v.RegisterValidation("sku", validateSKU)

	return v.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku format : abc-123-xy
	fs := fl.Field().String()
	re := regexp.MustCompile(`[a-z]{3}-[0-9]{3}-[a-z]{2}`)
	m := re.FindAllString(fs, -1)

	if len(m) != 1 {
		return false
	}

	return true

}

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

// Dummy data
var productList = Products{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frosty milky coffee",
		Price:       2.45,
		SKU:         "abc-123-aa",
		CreatedOn:   time.Now().Local().String(),
		UpdatedOn:   time.Now().Local().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffe without milk",
		Price:       1.29,
		SKU:         "asd-321-bb",
		CreatedOn:   time.Now().Local().String(),
		UpdatedOn:   time.Now().Local().String(),
	},
}
