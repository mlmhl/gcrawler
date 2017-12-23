package request

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Request interface {
	URL() *url.URL
	Method() string
	Description() string
}

func New(url *url.URL, method string) Request {
	return &minimumRequest{url, method}
}

type minimumRequest struct {
	url    *url.URL
	method string
}

func (req *minimumRequest) URL() *url.URL {
	return req.url
}

func (req *minimumRequest) Method() string {
	return req.method
}

func (req *minimumRequest) Description() string {
	return fmt.Sprintf("%s %s", req.method, req.url.String())
}

type HeaderProvider interface {
	Header() http.Header
}

func WithHeader(req Request, header http.Header) Request {
	return &headerRequest{req, header}
}

type headerRequest struct {
	Request
	header http.Header
}

func (req *headerRequest) Header() http.Header {
	return req.header
}

type CookiesProvider interface {
	Cookies() []*http.Cookie
}

func WithCookies(req Request, cookies []*http.Cookie) Request {
	return &cookiesRequest{req, cookies}
}

type cookiesRequest struct {
	Request
	cookies []*http.Cookie
}

func (req *cookiesRequest) Cookies() []*http.Cookie {
	return req.cookies
}

type ReaderProvider interface {
	Reader() io.Reader
}

func WithReader(req Request, reader io.Reader) Request {
	return &readerRequest{req, reader}
}

type readerRequest struct {
	Request
	reader io.Reader
}

func (req *readerRequest) Reader() io.Reader {
	return req.reader
}
