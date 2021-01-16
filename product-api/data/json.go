package data

import (
	"encoding/json"
	"io"
)

// ToJSON used for converting the struct of Product
// to JSON. It has better performance than json.marshal
// since it doesn't have to buffer the output into memory
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// FromJSON used for converting JSON
// formatted data to struct of Product
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}
