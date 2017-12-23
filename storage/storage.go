package storage

import "github.com/mlmhl/gcrawler/types"

type Storage interface {
	Name() string
	Put(item types.Item) error
}
