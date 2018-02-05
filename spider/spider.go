package spider

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/mlmhl/gcrawler/request"
	"github.com/mlmhl/gcrawler/response"

	"github.com/golang/glog"
)

type stopReason int

func (r stopReason) Reason() string {
	switch r {
	case timeout:
		return "Timeout"
	case exhausted:
		return "Exhausted"
	case manualStop:
		return "Manual stop"
	}
	return ""
}

const (
	timeout stopReason = iota
	exhausted
	manualStop
)

// Spider is an independent instance of a web crawler, it receives a list of
// bootstrap requests, executes each requests and use a Handler to process the responses.
type Spider struct {
	*Options

	cond        *sync.Cond
	lock        *sync.Mutex
	requests    []request.Request
	pendingJobs int32
	runningJobs int32

	stopChan   chan struct{}
	stopReason stopReason
}

func New(options *Options) (*Spider, error) {
	if options.Handler == nil {
		return nil, errors.New("empty handler")
	}
	if len(options.Storages) == 0 {
		return nil, errors.New("no storage specified")
	}
	if options.Client == nil {
		options.Client = http.DefaultClient
	}

	spider := &Spider{
		Options:  options,
		lock:     &sync.Mutex{},
		requests: options.Bootstraps,
		stopChan: make(chan struct{}),
	}
	spider.cond = sync.NewCond(spider.lock)
	// Make bootstrap requests can be GC.
	options.Bootstraps = nil
	return spider, nil
}

func (spider *Spider) Run() {
	go func() {
		if spider.LifeTime > time.Duration(0) {
			time.AfterFunc(spider.LifeTime, func() { spider.stop(timeout) })
		}
	}()

	glog.Infof("Start crawling...")

	spider.lock.Lock()
	defer spider.lock.Unlock()
	for {
		for len(spider.requests) == 0 && !spider.stopped() && spider.unfinished() {
			spider.cond.Wait()
		}
		if spider.stopped() {
			break
		}
		if len(spider.requests) == 0 && !spider.unfinished() {
			// If no pending requests and running jobs, we can just stop this spider.
			spider.doStop(exhausted)
			break
		}
		spider.pendingJobs++
		go spider.process(spider.requests[0])
		spider.requests = spider.requests[1:]
	}

	glog.Infof("Stop crawling: %s", spider.stopReason.Reason())
}

func (spider *Spider) unfinished() bool {
	return spider.pendingJobs > 0 || spider.runningJobs > 0
}

func (spider *Spider) Stop() {
	spider.stop(manualStop)
}

func (spider *Spider) appendRequests(requests []request.Request) {
	spider.lock.Lock()
	defer spider.lock.Unlock()
	spider.requests = append(spider.requests, requests...)
	spider.cond.Broadcast()
}

func (spider *Spider) process(req request.Request) {
	spider.lock.Lock()
	for !spider.stopped() && spider.reachLimit() {
		spider.cond.Wait()
	}
	if spider.stopped() {
		glog.V(4).Infof("Spider stopped, ignore request: %s", req.Description())
		spider.lock.Unlock()
		return
	}
	spider.pendingJobs--
	spider.runningJobs++
	spider.lock.Unlock()

	glog.V(4).Infof("Process request: %s", req.Description())

	var resp *http.Response
	r, err := http.NewRequest(req.Method(), req.URL().String(), nil)
	if err == nil {
		resp, err = spider.Client.Do(r)
	}
	items, successors := spider.Handler.Handle(req, response.NewResponse(resp, err))
	for _, s := range spider.Storages {
		for _, item := range items {
			if err := s.Put(item); err != nil {
				glog.Warningf("Put item to storage %s failed: %v", s.Name(), err)
			}
		}
	}
	spider.lock.Lock()
	defer spider.lock.Unlock()
	spider.runningJobs--
	spider.requests = append(spider.requests, successors...)
	spider.cond.Broadcast()
}

func (spider *Spider) stop(r stopReason) {
	spider.lock.Lock()
	defer spider.lock.Unlock()
	spider.doStop(r)
}

func (spider *Spider) doStop(r stopReason) {
	close(spider.stopChan)
	spider.stopReason = r
	spider.cond.Broadcast()
}

func (spider *Spider) stopped() bool {
	select {
	case <-spider.stopChan:
		return true
	default:
	}
	return false
}

func (spider *Spider) reachLimit() bool {
	return spider.Concurrency > 0 && spider.runningJobs >= spider.Concurrency
}
