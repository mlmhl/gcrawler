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