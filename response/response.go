package response

import "net/http"

type Response interface {
	Response() *http.Response
	Error() error
}

func NewResponse(resp *http.Response, err error) Response {
	return response{resp, err}
}

type response struct {
	resp *http.Response
	error error
}

func (resp response) Response() *http.Response {
	return resp.resp
}

func (resp response) Error() error {
	return resp.error
}
