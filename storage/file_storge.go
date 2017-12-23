package storage

import (
	"io"
	"os"

	"github.com/mlmhl/gcrawler/types"
)

const fileStorageName = "File"

var _ Storage = FileStorage{}

type FileStorage struct {
	file *os.File
}

func NewFileStorage(path string) (FileStorage, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return FileStorage{}, err
	}
	return FileStorage{file}, nil
}

func (s FileStorage) Close() {
	if s.file != nil {
		s.file.Close()
	}
}

func (s FileStorage) Put(item types.Item) error {
	content := item.Content() + "\n"
	n, err := s.file.WriteString(content)
	if err == nil && n < len(content) {
		err = io.ErrShortWrite
	}
	return err
}

func (s FileStorage) Name() string {
	return fileStorageName
}
