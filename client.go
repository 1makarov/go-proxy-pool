package proxypool

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
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
	if setting.MaxCountConn == 0 {
		setting.MaxCountConn = 3
	}
	if setting.Timeout == 0 {
		setting.Timeout = time.Second * 5
	}

	return &Client{
		setting: setting,
		storage: newStorage(setting.MaxCountConn),
	}, nil
}

func (c *Client) AddArray(proxyRaw []string, thread int) error {
	ch := make(chan error, 1)
	var stopped bool

	go func() {
		wg := &sync.WaitGroup{}

		for i, p := range proxyRaw {
			if stopped {
				return
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

func (c *Client) Get() *Proxy {
	return c.storage.get()
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
