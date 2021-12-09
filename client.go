package proxypool

import (
	"fmt"
	"github.com/1makarov/go-proxy-pool/proxy"
	"sync"
)

type Client struct {
	proxyClient  *proxy.Client
	storage      *Storage
	maxCountConn int
}

func New(maxCountConn int, proxyClient *proxy.Client) (*Client, error) {
	if maxCountConn == 0 {
		maxCountConn = 3
	}

	return &Client{
		proxyClient:  proxyClient,
		storage:      newStorage(maxCountConn),
		maxCountConn: maxCountConn,
	}, nil
}

func (c *Client) AddArray(proxyRaw []string, thread int) error {
	ch := make(chan error, 1)
	var stopped bool

	go func() {
		wg := &sync.WaitGroup{}

		for i, p := range proxyRaw {
			if stopped {
				break
			}

			proxy := p

			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := c.Add(proxy); err != nil {
					ch <- err
				}
			}()

			if i == 0 {
				continue
			}

			if i%thread == 0 {
				wg.Wait()
			}
		}
		wg.Wait()

		close(ch)
	}()

	for err := range ch {
		stopped = true
		return err
	}

	return nil
}

func (c *Client) Add(proxyRaw string) error {
	proxy, err := c.validate(proxyRaw)
	if err != nil {
		return fmt.Errorf("error add proxy: %s, %w", proxyRaw, err)
	}

	c.storage.add(proxy)

	return nil
}

func (c *Client) Get() *types.Proxy {
	return c.storage.get()
}

func (c *Client) Close(proxy *types.Proxy) {
	c.storage.close(proxy)
}
