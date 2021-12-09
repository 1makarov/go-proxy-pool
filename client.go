package proxypool

import (
	"fmt"
	"time"
)

const (
	defaultMaxCountConn = 3
	defaultTestTimeout  = 5 * time.Second
	defaultTestURL      = "https://api.ip.sb/ip"
)

type Setting struct {
	TestURL      string
	TestTimeout  time.Duration
	MaxCountConn int
}

type Client struct {
	proxyClient  *proxyClient
	storage      *storage
	maxCountConn int
}

func New(s Setting) (*Client, error) {
	if s.MaxCountConn == 0 {
		s.MaxCountConn = defaultMaxCountConn
	}

	if s.TestTimeout == 0 {
		s.TestTimeout = defaultTestTimeout
	}

	if s.TestURL == "" {
		s.TestURL = defaultTestURL
	}

	return &Client{
		proxyClient: newProxyClient(s.TestURL, s.TestTimeout),
		storage:     newStorage(s.MaxCountConn),
	}, nil
}

func (c *Client) Add(rawProxy string) error {
	proxy, err := c.proxyClient.New(rawProxy)
	if err != nil {
		return fmt.Errorf("error add proxy: %s, %w", rawProxy, err)
	}

	c.storage.add(proxy)

	return nil
}

func (c *Client) Get() *Proxy {
	return c.storage.get()
}

func (c *Client) Close(proxy *Proxy) {
	c.storage.close(proxy)
}
