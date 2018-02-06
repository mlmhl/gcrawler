package storage

import "github.com/mlmhl/gcrawler/types"

// Storage is used to store the Item crawled by Spider.
type Storage interface {
	Name() string
	Put(item types.Item) error
}
