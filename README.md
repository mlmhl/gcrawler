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
