package spider

import (
	"net/http"
	"time"

	"github.com/mlmhl/gcrawler/handler"
	"github.com/mlmhl/gcrawler/request"
	"github.com/mlmhl/gcrawler/storage"
)

// Options is the set of all required and optional parameters to initialize a Spider.
type Options struct {
	Client      *http.Client
	Handler     handler.Handler
	LifeTime    time.Duration
	Storages    []storage.Storage
	Bootstraps  []request.Request
	Concurrency int32
}
