package proxypool

import (
	"fmt"
	"sync"
)

type Storage struct {
	proxies      map[string]*Proxy
	mu           sync.Mutex
	maxCountConn int
}

func newStorage(maxCountConn int) *Storage {
	proxies := make(map[string]*Proxy)

	return &Storage{proxies: proxies, maxCountConn: maxCountConn}
}

func (s *Storage) add(proxy *Proxy) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.proxies[proxy.url.String()] = proxy
}

func (s *Storage) get() (*Proxy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, proxy := range s.proxies {
		if proxy.count >= s.maxCountConn {
			continue
		}

		proxy.count += 1
		return proxy, nil
	}

	return nil, fmt.Errorf("empty proxy")
}

func (s *Storage) close(proxy *Proxy) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, ok := s.proxies[proxy.url.String()]; ok {
		v.count -= 1
	}
}
