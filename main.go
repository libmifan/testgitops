package main

import (
	"log"
	"mycrawler/collect"
	"mycrawler/proxy"
	"time"
)

func init() {
	log.SetPrefix("HTTP DEV: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	//proxyURLs := []string{"https://sig.dahuangsz.com:443"}
	proxyURLs := []string{"https://whatname:USA*678.com@sig.dahuangsz.com:443"}
	url := "https://www.google.com"

	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		log.Println("RoundRobinProxySwitcher failed")
	}

	var f collect.Fetcher
	f = collect.BrowserFetch{
		Timeout: 5000 * time.Millisecond,
		Proxy:   p,
	}

	body, err := f.Get(url)
	if err != nil {
		log.Printf("Failed to fetch url: %v", err)
		//return
	}
	body, _ = f.Get(url)

	log.Printf("Body:\n %s", string(body))
	//doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
}
