package main

import (
	"fmt"
	"github.com/1makarov/go-proxy-pool"
	"log"
	"time"
)

func main() {
	client, err := proxypool.New(proxypool.Setting{
		MaxCountConn: 3,
		TestTimeout:  5 * time.Second,
		TestURL:      "https://api.ip.sb/ip",
	})
	if err != nil {
		log.Fatalln(err)
	}

	if err = client.Add("http://user:password@host:port"); err != nil {
		log.Fatalln(err)
	}

	proxy := client.Get()
	fmt.Println(proxy)

	client.Close(proxy)
}
