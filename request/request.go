package request

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// A Request is an interface that tells Spider which URL to fetch,
// and which HTTP method to use(i.e. "GET","HEAD",...).
type Request interface {
	// URL returns the url will be crawled.
	URL() *url.URL
	// Method returns the HTTP method used to send request.
	Method() string
	// Description returns a human readable unique key of the Request.
	Description() string
}

// New returns a new Request with specified url and HTTP method.
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

// HeaderProvider is an enhanced interface a Request may implement.
type HeaderProvider interface {
	// Header returns the header used to make HTTP request.
	Header() http.Header
}

// WithHeader returns a new Request base on the specified Request and header,
// the returned Request implements HeaderProvider interface.
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

// CookiesProvider is an enhanced interface a Request may implement.
type CookiesProvider interface {
	// Cookies returns the coolies used to make HTTP request.
	Cookies() []*http.Cookie
}

// WithCookies returns a new Request base on the specified Request and cookies,
// the returned Request implements CookiesProvider.
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

// ReaderProvider is an enhanced interface a Request may implement.
type ReaderProvider interface {
	// Reader return an io.Reader, which will be used as the post data in HTTP request.
	Reader() io.Reader
}

// WithCookies returns a new Request base on the specified Request and reader,
// the returned Request implements HeaderProvider ReaderProvider.
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
