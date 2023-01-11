package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"mycrawler/collect"
	"mycrawler/log"
	"mycrawler/parse/doubangroup"
	"time"
)

func main() {
	////proxyURLs := []string{"https://sig.dahuangsz.com:443"}
	//proxyURLs := []string{"https://whatname:USA*678.com@sig.dahuangsz.com:443"}
	//url := "https://www.google.com"
	//
	//p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	//if err != nil {
	//	log.Println("RoundRobinProxySwitcher failed")
	//}
	//
	//var f collect.Fetcher
	//f = collect.BrowserFetch{
	//	Timeout: 5000 * time.Millisecond,
	//	Proxy:   p,
	//}
	//
	//body, err := f.Get(url)
	//if err != nil {
	//	log.Printf("Failed to fetch url: %v", err)
	//	return
	//}
	//
	//log.Printf("Body:\n %s", string(body))
	////doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))

	plugin, c := log.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	// douban cookie
	cookie := "viewed=\"1007305\"; bid=qibrMJRpIY4; gr_user_id=baadbf6a-ff0d-4a50-b89b-b9a2c5397576; __gads=ID=a2baf8eaa6209382-2228eb6339d900f9:T=1673339647:RT=1673339647:S=ALNI_MagnwXcDKP6ihN-CjoHj4Z2YiwKLQ; __utmc=30149280; __utmz=30149280.1673339650.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __gpi=UID=00000ba12fe4c036:T=1673339647:RT=1673432806:S=ALNI_MbJcM8oc8TTwFicRqSnrha7bn31zA; __yadk_uid=RtMVLztsoJP3GzpDvOUdRAmzjdJc6AMg; douban-fav-remind=1; dbcl2=\"41494661:STARhSgYVhY\"; ck=sQGR; push_noty_num=0; push_doumail_num=0; __utmv=30149280.4149; _pk_id.100001.8cb4=6e202c036ec9a16d.1673432805.2.1673446236.1673434741.; _pk_ses.100001.8cb4=*; __utma=30149280.501428854.1673339650.1673432806.1673446237.4; __utmt=1; __utmb=30149280.2.10.1673446237"

	var workList []*collect.Request
	for i := 25; i <= 100; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		workList = append(workList, &collect.Request{
			Url:       str,
			Cookie:    cookie,
			ParseFunc: doubangroup.ParseURL,
		})
	}

	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
	}

	for len(workList) > 0 {
		items := workList
		workList = nil

		for _, item := range items {
			body, err := f.Get(item)
			time.Sleep(1 * time.Second)

			if err != nil {
				logger.Error("read content failed", zap.Error(err))
				continue
			}
			res := item.ParseFunc(body, item)
			for _, item := range res.Items {
				logger.Info("get url", zap.String("url", item.(string)))
			}
			workList = append(workList, res.Requesrts...)
		}
	}
}
