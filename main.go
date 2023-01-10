package main

import (
	"log"
	"mycrawler/collect"
	"time"
)

func init() {
	log.SetPrefix("HTTP DEV: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	url := "https://book.douban.com/subject/1007305/"

	var f collect.Fetcher
	f = collect.BrowserFetch{
		Timeout: 500 * time.Millisecond,
	}

	body, err := f.Get(url)
	if err != nil {
		log.Println("Failed to fetch url:", err)
		return
	}

	log.Printf("Body:\n %s", string(body))

}
