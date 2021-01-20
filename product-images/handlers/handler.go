package handlers

import (
	"github.com/fahmi1597/microservices-go/product-images/files"
	"github.com/hashicorp/go-hclog"
)

// File is a thing
type File struct {
	log   hclog.Logger
	store files.StorageIO
}

// NewFileHandler creates a new handler
func NewFileHandler(l hclog.Logger, s files.StorageIO) *File {
	return &File{l, s}
}
