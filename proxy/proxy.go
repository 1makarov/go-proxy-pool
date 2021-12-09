package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const defaultTimeout = time.Second * 5

type Client struct {
	TestURL     string
	TestTimeout time.Duration
}

type Proxy struct {
	URL *url.URL
}

func NewClient(testURL string, testTimeout time.Duration) *Client {
	if testTimeout == 0 {
		testTimeout = defaultTimeout
	}

	return &Client{
		TestURL:     testURL,
		TestTimeout: testTimeout,
	}
}

func (c *Client) NewProxy(raw string) (*Proxy, error) {
	proxy, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if err = c.request(proxy); err != nil {
		return nil, err
	}

	return &Proxy{URL: proxy}, nil
}

func (c *Client) request(proxy *url.URL) error {
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	client := http.Client{
		Transport: transport,
		Timeout:   c.TestTimeout,
	}

	resp, err := client.Get(c.TestURL)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	}

	return fmt.Errorf("error request")
}
