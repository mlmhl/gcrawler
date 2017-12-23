package handler

import (
	"github.com/mlmhl/gcrawler/request"
	"github.com/mlmhl/gcrawler/response"
	"github.com/mlmhl/gcrawler/types"
)

type Handler interface {
	Handle(req request.Request, resp response.Response) (items []types.Item, successors []request.Request)
}
