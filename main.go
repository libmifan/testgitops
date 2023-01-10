package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func mainDpre() {
	url := "https://www.thepaper.cn/"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("fether url error:%v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error status code: %v", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll error:", err)
		return
	}

	//fmt.Println("body:", string(body))
	n := strings.Count(string(body), "<a")
	bn := bytes.Count(body, []byte("<a"))
	if n != bn {
		fmt.Printf("n: %d, bn: %d\n", n, bn)
	}
	fmt.Println("contains url: ", n)

}

func init() {
	log.SetPrefix("HTTP DEV: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	url := "https://www.thepaper.cn/"

	fetch, err := Fetch(url)
	if err != nil {
		log.Println("Failed to fetch url:", err)
		return
	}

	doc, err := htmlquery.Parse(bytes.NewReader(fetch))
	if err != nil {
		log.Printf("htmlquery.Parse failed:%v\n", err)
	}
	//nodes := htmlquery.Find(doc, `//div[@class="news_li"]/h2/a[@target="_blank"`)
	nodes := htmlquery.Find(doc, `//div[@class="small_toplink__GmZhY"]/a/h2`)

	for _, node := range nodes {
		log.Println("Fetch card ", node.FirstChild.Data)
	}
}

func Fetch(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("not the expected status:", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	e := DeterminEncoding(reader)
	utf8Reader := transform.NewReader(reader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

func DeterminEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)

	if err != nil {
		log.Println("fetch error:", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
