package proxypool

import "net/url"

type Proxy struct {
	url   *url.URL
	count int
}
