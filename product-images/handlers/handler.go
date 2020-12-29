package handlers

import (
	"log"

	"github.com/fahmi1597/microservices-go/product-images/files"
)

// File is a thing
type File struct {
	log   *log.Logger
	store files.StorageIO
}

// NewFileHandler creates a new handlers
func NewFileHandler(l *log.Logger, s files.StorageIO) *File {
	return &File{l, s}
}
