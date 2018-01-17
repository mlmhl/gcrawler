package response

import "net/http"

// Response is the HTTP response for a spider request.
type Response interface {
	// Response returns the HTTP response if succeeded, otherwise return nil.
	Response() *http.Response
	// Error returns the error if exist, otherwise return nil.
	Error() error
}

// NewResponse create a new Response with specified HTTP response and error.
func NewResponse(resp *http.Response, err error) Response {
	return response{resp, err}
}

type response struct {
	resp  *http.Response
	error error
}

func (resp response) Response() *http.Response {
	return resp.resp
}

func (resp response) Error() error {
	return resp.error
}
