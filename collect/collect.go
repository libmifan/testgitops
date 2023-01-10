package collect

import (
	"bufio"
	"context"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Fetcher interface {
	Get(url string) ([]byte, error)
}

type ChromedpFetcher struct{}

func (ChromedpFetcher) Get(url string) ([]byte, error) {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var body string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body > footer`),
	)
	if err != nil {
		return nil, err
	}
	return []byte(body), nil
}

type BrowserFetch struct{}

func (BrowserFetch) Get(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to do request", err)
		return nil, err
	}
	defer resp.Body.Close()

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
