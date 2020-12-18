package data

import (
	"encoding/json"
	"io"
)

// ToJSON used for converting the struct of Product
// to JSON. It has better performance than json.marshal
// since it doesn't have to buffer the output into memory
func ToJSON(p interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// FromJSON used for converting JSON
// formatted data to struct of Product
func FromJSON(p interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
