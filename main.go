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
	proxyURLs := []string{"http://whatname:USA*678.com@43.159.63.218:443"}
	url := "https://www.google.com"

	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		log.Println("RoundRobinProxySwitcher failed")
	}

	var f collect.Fetcher
	f = collect.BrowserFetch{
		Timeout: 500 * time.Millisecond,
		Proxy:   p,
	}

	body, err := f.Get(url)
	if err != nil {
		log.Println("Failed to fetch url:", err)
		return
	}

	log.Printf("Body:\n %s", string(body))
	//doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
}
