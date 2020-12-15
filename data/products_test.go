package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "test",
		Price: 1.00,
		SKU:   "aBc-123-ab",
	}

	if err := p.Validate(); err != nil {
		t.Fatal(err.Error())
	}

	// to do : test with multiple product format
}
