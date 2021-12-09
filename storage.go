package proxypool

import (
	"github.com/1makarov/go-proxy-pool/types"
	"sync"
	"time"
)

type storage struct {
	proxies      map[string]*types.Proxy
	mu           sync.Mutex
	maxCountConn int
}

func newStorage(maxCountConn int) *storage {
	proxies := make(map[string]*types.Proxy)

	return &storage{proxies: proxies, maxCountConn: maxCountConn}
}

func (s *storage) add(proxy *types.Proxy) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.proxies[proxy.url.String()] = proxy
}

func (s *storage) get() *types.Proxy {
	for {
		s.mu.Lock()

		for _, proxy := range s.proxies {
			if proxy.count >= s.maxCountConn {
				continue
			}

			proxy.count += 1
			return proxy
		}

		s.mu.Unlock()
		time.Sleep(time.Second)
	}
}

func (s *storage) close(proxy *types.Proxy) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, ok := s.proxies[proxy.url.String()]; ok {
		v.count -= 1
	}
}
