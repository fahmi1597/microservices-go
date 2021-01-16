package data

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidData(t *testing.T) {
	data := `{"name":"abc","price":1.22,"sku":"abc-123-abc"}`
	p := Product{}

	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		fmt.Print(err)
		return
	}

	v := NewValidator()

	// Going to panic since valid data will not return a ValidationErrors, it return nil value instead
	// It's okay.
	errs := v.Validate(&p)
	assert.Len(t, errs, 1)
}

func TestProductMissingNameReturnsErr(t *testing.T) {
	data := `{"price":1.22, "sku":"abc-123-abc"}`
	p := Product{}

	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		fmt.Print(err)
		return
	}

	v := NewValidator()
	vErr := v.Validate(&p)

	errMsg := fmt.Errorf("%s", vErr.Errors())
	if errMsg != nil {
		fmt.Println(errMsg)
		t.Fail()
		return
	}

	assert.Len(t, errMsg, 1)
}

func TestProductMissingPriceReturnsErr(t *testing.T) {
	data := `{"name":"abc","sku":"abc-abc-abc"}`
	p := Product{}
	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		fmt.Print(err)
		return
	}

	v := NewValidator()
	vErr := v.Validate(&p)

	errMsg := fmt.Errorf("%s", vErr.Errors())
	if errMsg != nil {
		fmt.Println(errMsg)
		t.Fail()
		return
	}

	assert.Len(t, v.Validate(&p), 1)
}

func TestProductMissingSKUReturnsErr(t *testing.T) {
	data := `{"name":"abc", "price":1.22}`
	p := Product{}
	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		fmt.Print(err)
		return
	}

	v := NewValidator()
	vErr := v.Validate(&p)

	errMsg := fmt.Errorf("%s", vErr.Errors())
	if errMsg != nil {
		fmt.Println(errMsg)
		t.Fail()
		return
	}

	assert.Len(t, v.Validate(&p), 1)
}

func TestProductInvalidSKUReturnsErr(t *testing.T) {
	data := `{
		"name":  "abc",
		"price": 1.22,
		"sku": "abc-1-ab"
		
	}`
	p := Product{}

	if err := json.Unmarshal([]byte(data), &p); err != nil {
		return
	}

	v := NewValidator()
	vErr := v.Validate(&p)

	errMsg := fmt.Errorf("%s", vErr.Errors())
	if errMsg != nil {
		fmt.Println(errMsg)
		t.Fail()
		return
	}

	assert.Len(t, errMsg, 1)

}
