# gcrawler

Package gcrawler provides a very simple and flexible web crawler that supports three running modes:

* Running a fixed period of time
* Running forever until no more pages to crawled
* Running forever until explicitly stopped by the `Stop` method

## Installation

To install, simply run in a terminal:

    go get github.com/mlmhl/gcrawler

## Usage

The following example (taken from /example/github_repo/main.go) shows how to create and start a `Spider`, running forever until no more pages to crawled:

```go
func main() {
	flag.Parse()

	u, err := url.Parse("https://github.com/mlmhl?tab=repositories")
	if err != nil {
		log.Fatalf("Parse start url failed: %v", err)
	}
	s, err := spider.New(&spider.Options{
		Handler:    newGitHubRepoHandler(),
		Storages:   []storage.Storage{storage.NewConsoleStorage()},
		Bootstraps: []request.Request{request.New(u, "GET")},
	})
	if err != nil {
		glog.Fatalf("Create spider failed: %v", err)
	}
	s.Run()
}
```

## Spider

A `Spider` is an independent instance of a web crawler, it receives a list of bootstrap requests, executes each requests and use a `Handler` to process the responses.

You can create a new `Spider` with an `Options`. `Options` contains a series of parameters to create a custom `Spider`:

- Client: The http client `Spider` used to make requests. By default, the net/http default Client will be used.

- LifeTime: If set, the `Spider` will only runs for such duration.

- Bootstraps: Bootstrap requests to be executed.

- Concurrency: If set, no more than such requests will be executed concurrently.

## Requset

A `Request` is an interface that tells `Spider` which URL to fetch, and whitch HTTP method to use(i.e. "GET","HEAD",...).

It is defined like so:

```go
type Request interface {
	URL() *url.URL
	Method() string
	Description() string
}
```

We recognizes a number of interfaces that the `Request` may implement, for more advanced needs.

- `HeaderProvider`: Implement this interface to specify the headers of the `Request`.

- `CookiesProvider`: Implement this interface to specify the cookies of the `Request`.

- `ReaderProvider`: Implement this interface to specify the body of the `Request` via a `io.Reader`.

## Response

A `Response` is the response of a `Request`, implemented as simple wrapper of `*http.Response` with an `Error` method. It is defined like so:

```go
type Response interface {
	Response() *http.Response
	Error() error
}
```

## Handler

A `Handler` responds to a `Response`, it has only one method as follows:

```go
type Handler interface {
	Handle(req request.Request, resp response.Response) (items []types.Item, successors []request.Request)
}
```

The `Handle` method usually parses the returned data, and generates a series of `Item` and `Request`. The spider will goon processing these successors if any. An `Item` can be anything that the user wants to crawl.

## Storage

A `Storage` is used to store the `Item` crawled by `Spider`. We provide two simple `Storage`: `ConsoleStorage` and `FileStorage`, which usually used for debugging.
Just as their name implies, a `ConsoleStorage` prints all `Item` to the console, and a `FileStorage` write all `Item` to a file on local disk.
Usually you should provide your own `Storage` in production environment.
