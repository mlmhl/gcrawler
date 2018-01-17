package handler

import (
	"github.com/mlmhl/gcrawler/request"
	"github.com/mlmhl/gcrawler/response"
	"github.com/mlmhl/gcrawler/types"
)

// A Handler responds to a Response, always provided by user.
type Handler interface {
	// Handle usually parse the returned data in Response, and generates a series
	// of Item and Request. The spider will goon processing these successors if any.
	Handle(req request.Request, resp response.Response) (items []types.Item, successors []request.Request)
}
