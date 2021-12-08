package proxypool

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	setting Setting
	storage *Storage
}

type Setting struct {
	MaxCountConn int
	TestURL      string
	Timeout      time.Duration
}

func New(setting Setting) (*Client, error) {
	if setting.TestURL == "" {
		return nil, fmt.Errorf("empty test url")
	}

	return &Client{
		setting: setting,
		storage: newStorage(setting.MaxCountConn),
	}, nil
}

func (c *Client) Add(proxyRaw string) error {
	proxy, err := c.validate(proxyRaw)
	if err != nil {
		return fmt.Errorf("error add proxy: %s, %w", proxyRaw, err)
	}

	c.storage.add(proxy)

	return nil
}

func (c *Client) Get() (*Proxy, error) {
	for {
		proxy, err := c.storage.get()
		if err != nil {
			time.Sleep(c.setting.Timeout)
			continue
		}

		return proxy, nil
	}
}

func (c *Client) Close(proxy *Proxy) {
	c.storage.close(proxy)
}

func (c *Client) validate(proxyRaw string) (*Proxy, error) {
	proxy, err := url.Parse(proxyRaw)
	if err != nil {
		return nil, err
	}

	if err = c.request(proxy); err != nil {
		return nil, err
	}

	return &Proxy{url: proxy}, nil
}

func (c *Client) request(proxy *url.URL) error {
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	client := http.Client{
		Transport: transport,
		Timeout:   c.setting.Timeout,
	}

	resp, err := client.Get(c.setting.TestURL)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	}

	return fmt.Errorf("error request")
}
