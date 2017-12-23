package storage

import (
	"fmt"
	"sync"

	"github.com/mlmhl/gcrawler/types"
)

const consoleStorageName = "Console"

var _ Storage = &consoleStorage{}

type consoleStorage struct {
	lock sync.Mutex
}

func NewConsoleStorage() Storage {
	return &consoleStorage{}
}

func (s *consoleStorage) Name() string {
	return consoleStorageName
}

func (s *consoleStorage) Put(item types.Item) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	fmt.Println(item.Content())
	return nil
}
