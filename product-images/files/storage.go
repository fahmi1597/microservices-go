package files

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

// StorageIO defines the behavior for file operations
// Implementations may be of the time local disk, or cloud storage, etc
type StorageIO interface {
	Write(path string, file io.Reader) error
}

// Storage is a thing
type Storage struct {
	basePath    string
	maxFileSize int
}

// New is a thing
func New(basePath string, maxSize int) (*Storage, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Storage{basePath: p}, nil
}

// Save is a thing
// It implements StorageIO interface
func (sl *Storage) Write(path string, contents io.Reader) error {

	fp := sl.fullPath(path)
	dir := filepath.Dir(fp)
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf(": %w", err)
	}

	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

// returns absolute path
func (sl *Storage) fullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(sl.basePath, path)
}
