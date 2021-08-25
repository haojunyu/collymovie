package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	t := time.Now()

	// Instantiate default collector
	c := colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		//colly.Async(true),
		colly.AllowedDomains("movie.douban.com"),
		/*
			colly.URLFilters(
				regexp.MustCompile("https://movie\\.douban\\.com/subject/.+$"),
				regexp.MustCompile("https://movie\\.douban\\.com/celebrity/.*$"),
			),
		*/
	)

	// c.Limit(&colly.LimitRule{
	// 	DomainGlob:  "",
	// 	Delay:       1 * time.Second,
	// 	RandomDelay: 1 * time.Second,
	// })

	// 设置随机头
	//extensions.RandomUserAgent(c)
	//extensions.Referer(c)

	// mongodb 作数据存储
	/*
		storage := &mongo.Storage{
			Database: "colly",
			URI:      "mongodb://101.132.117.101:8017",
		}

		if err := c.SetStorage(storage); err != nil {
			panic(err)
		}
	*/

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://movie.douban.com/top250")
	c.Wait()
	fmt.Printf("花费时间:%s", time.Since(t))
}
