package proxypool

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type proxyClient struct {
	TestURL     string
	TestTimeout time.Duration
}

func newProxyClient(testURL string, testTimeout time.Duration) *proxyClient {
	return &proxyClient{TestURL: testURL, TestTimeout: testTimeout}
}

type Proxy struct {
	URL *url.URL
}

func (c *proxyClient) new(rawProxy string) (*Proxy, error) {
	proxy, err := url.Parse(rawProxy)
	if err != nil {
		return nil, err
	}

	if err = c.checkOnRequest(proxy); err != nil {
		return nil, err
	}

	return &Proxy{URL: proxy}, nil
}

func (c *proxyClient) checkOnRequest(proxy *url.URL) error {
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
