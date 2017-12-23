package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/mlmhl/gcrawler/handler"
	"github.com/mlmhl/gcrawler/request"
	"github.com/mlmhl/gcrawler/response"
	"github.com/mlmhl/gcrawler/spider"
	"github.com/mlmhl/gcrawler/types"

	"github.com/golang/glog"
	"github.com/mlmhl/gcrawler/storage"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type state interface {
	Parse(n *html.Node) (string, state)
}

type rootState struct{}

func (s rootState) Parse(n *html.Node) (string, state) {
	if isRepoListNode(n) {
		return "", repoListNodeState{}
	}
	return "", s
}

func isRepoListNode(n *html.Node) bool {
	if n.Data == atom.Div.String() {
		for _, a := range n.Attr {
			if a.Key == atom.Id.String() && a.Val == "user-repositories-list" {
				return true
			}
		}
	}
	return false
}

// <div id="user-repositories-list">
type repoListNodeState struct{}

func (s repoListNodeState) Parse(n *html.Node) (string, state) {
	if isRepoList(n) {
		return "", repoListState{}
	}
	return "", nil
}

func isRepoList(n *html.Node) bool {
	return n.Data == atom.Ul.String()
}

// <ul data-filterable-for="your-repos-filter" data-filterable-type="substring">
type repoListState struct{}

func (s repoListState) Parse(n *html.Node) (string, state) {
	if isRepoItem(n) {
		return "", repoNodeState{}
	}
	return "", nil
}

func isRepoItem(n *html.Node) bool {
	return n.Data == atom.Li.String()
}

// <li class="col-12 d-block width-full py-4 border-bottom public fork" itemprop="owns" itemscope itemtype="http://schema.org/Code">
type repoNodeState struct{}

func (s repoNodeState) Parse(n *html.Node) (string, state) {
	if isRepo(n) {
		return "", repoState{}
	}
	return "", nil
}

func isRepo(n *html.Node) bool {
	if n.Data == atom.Div.String() {
		for _, a := range n.Attr {
			if a.Key == atom.Class.String() && a.Val == "d-inline-block mb-1" {
				return true
			}
		}
	}
	return false
}

// <div class="d-inline-block mb-1">
type repoState struct{}

func (s repoState) Parse(n *html.Node) (string, state) {
	if isLink(n) {
		return "", linkState{}
	}
	return "", nil
}

func isLink(n *html.Node) bool {
	return n.Data == atom.H3.String()
}

// <h3> <a href="XXX">
type linkState struct{}

func (s linkState) Parse(n *html.Node) (string, state) {
	for _, attr := range n.Attr {
		if attr.Key == atom.Href.String() {
			return "http://github.com" + attr.Val, nil
		}
	}
	return "", nil
}

func newParser() *parser {
	return &parser{state: rootState{}}
}

type parser struct {
	state state
	items []types.Item
}

func (parser *parser) Parse(n *html.Node) {
	defer func(state state) { parser.state = state }(parser.state)
	repo, state := parser.state.Parse(n)
	if len(repo) > 0 {
		parser.items = append(parser.items, item(repo))
	}
	if state == nil {
		return
	}
	parser.state = state
	for children := n.FirstChild; children != nil; children = children.NextSibling {
		parser.Parse(children)
	}
}

type item string

func (i item) Content() string {
	return string(i)
}

type gitHubRepoHandler struct {
}

func newGitHubRepoHandler() handler.Handler {
	return &gitHubRepoHandler{}
}

func (processor *gitHubRepoHandler) Handle(req request.Request, resp response.Response) ([]types.Item, []request.Request) {
	if resp.Error() != nil {
		glog.Errorf("Crawl request %s failed: %v", req.Description(), resp.Error())
		// If we wanted to retry this request, we can return the request as a successor.
		// return nil, []request.Request{req}
		return nil, nil
	}

	httpResp := resp.Response()
	if httpResp.Body == nil {
		// We can do nothing.
		return nil, nil
	}

	node, err := html.Parse(httpResp.Body)
	if err != nil {
		glog.Errorf("Parse response failed: %s, %v", req.Description(), err)
		return nil, nil
	}

	p := newParser()
	p.Parse(node)
	return p.items, nil
}

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
