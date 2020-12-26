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

// New creates a new handlers
func New(l *log.Logger, s files.StorageIO) *File {
	return &File{l, s}
}
