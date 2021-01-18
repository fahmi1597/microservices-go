package data

import (
	"fmt"
	"testing"
)

func TestGetProductList(t *testing.T) {
	p := ProductsDB{}
	p.GetProductList("")
}

func TestGetProduct(t *testing.T) {
	p := ProductsDB{}
	p.GetProductByID(1, "")
}

func TestAddProduct(t *testing.T) {
	p := Product{
		Name:        "aaa",
		Description: "aaa",
		Price:       1.23,
		SKU:         "aaa-bbb-aaa",
	}

	pdb := ProductsDB{}
	pdb.AddProduct(p)
	ps, _ := pdb.GetProductList("")
	fmt.Printf("%#v \n", ps)
	p1, _ := pdb.GetProductByID(len(ps), "")
	fmt.Printf("%#v \n", p1)

}

func TestUpdateProduct(t *testing.T) {
	p := Product{
		ID:          1,
		Name:        "aaa",
		Description: "aaa",
		Price:       1.23,
		SKU:         "aaa-bbb-aaa",
	}
	pdb := ProductsDB{}
	pdb.GetProductList("")
	pdb.UpdateProduct(p)
	pdb.GetProductList("")
}
