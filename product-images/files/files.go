package files

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

// FileStorage defines the location and allocate the size of storage to store data file
type FileStorage struct {
	basePath    string
	maxFileSize int
}

// NewFileStorage return a new FileStorage
// with intended location and maximum file size allowed
func NewFileStorage(basePath string, maxSize int) (*FileStorage, error) {
	path, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &FileStorage{basePath: path}, nil
}

// Save is a thing
// It implements StorageIO interface
func (fs *FileStorage) Write(path string, contents io.Reader) error {

	fullPath := fs.getFullPath(path)
	dir := filepath.Dir(fullPath)
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	_, err = os.Stat(fullPath)
	if err == nil {
		err = os.Remove(fullPath)
		if err != nil {
			return xerrors.Errorf("%w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("%w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("%w", err)
	}

	return nil
}

// returns absolute path
func (fs *FileStorage) getFullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(fs.basePath, path)
}
