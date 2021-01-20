package files

import (
	"io"
)

// StorageIO defines the behavior of file operations
type StorageIO interface {
	Write(path string, file io.Reader) error
}
