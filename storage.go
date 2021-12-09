package proxypool

import (
	"sync"
	"time"
)

type item struct {
	proxy *Proxy
	count int
}

type storage struct {
	available    map[string]*item
	disabled     map[string]*item
	mu           sync.Mutex
	maxCountConn int
}

func newStorage(maxCountConn int) *storage {
	return &storage{
		available:    make(map[string]*item),
		disabled:     make(map[string]*item),
		maxCountConn: maxCountConn,
	}
}

func (s *storage) add(proxy *Proxy) {
	defer s.mu.Unlock()
	s.mu.Lock()

	s.available[proxy.URL.String()] = &item{proxy: proxy}
}

func (s *storage) get() *Proxy {
	for {
		s.mu.Lock()

		for _, i := range s.available {
			i.count += 1

			if i.count >= s.maxCountConn {
				delete(s.available, i.proxy.URL.String())
				s.disabled[i.proxy.URL.String()] = i
			}

			s.mu.Unlock()
			return i.proxy
		}

		s.mu.Unlock()
		time.Sleep(time.Second)
	}
}

func (s *storage) close(proxy *Proxy) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if i, ok := s.disabled[proxy.URL.String()]; ok {
		delete(s.disabled, proxy.URL.String())
		i.count -= 1
		s.available[proxy.URL.String()] = i
		return
	}

	i := s.available[proxy.URL.String()]
	if i.count > 0 {
		i.count -= 1
	}

	return
}
