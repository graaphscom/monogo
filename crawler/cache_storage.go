package crawler

import (
	"errors"
	"io"
	"os"
	"path"
)

var FileStorage = fileStorage{}

func (FS fileStorage) exists(key string) (bool, error) {
	if _, err := os.Stat(key); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
func (FS fileStorage) write(key string, value io.Reader) error {
	err := os.MkdirAll(path.Dir(key), 0750)
	if err != nil {
		return err
	}

	file, err := os.Create(key)
	defer func() { err = file.Close() }()
	if err != nil {
		return err
	}

	_, err = io.Copy(file, value)

	return err
}

type fileStorage struct{}

type CacheStorage interface {
	exists(key string) (bool, error)
	write(key string, value io.Reader) error
}
