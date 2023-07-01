package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"collymovie"

	"github.com/gocolly/colly"
)

func main() {

	t := time.Now()

	// Instantiate default collector
	c := colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
		//colly.AllowedDomains("movie.douban.com"),
		colly.UserAgent("Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Mobile Safari/537.36"),

		// colly.URLFilters(
		// 	regexp.MustCompile("https://movie\\.douban\\.com/subject/.+$"),
		// 	regexp.MustCompile("https://movie\\.douban\\.com/celebrity/.*$"),
		// ),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "",
		Delay:       1 * time.Second,
		RandomDelay: 5 * time.Second,
	})

	// 设置随机头
	//extensions.RandomUserAgent(c)
	//extensions.Referer(c)

	// Create another collector to scrape movie/celebrity details
	movieCollector := c.Clone()
	//colly.URLFilters(regexp.MustCompile("https://movie\\.douban\\.com/subject/[0-9]{8,10}/$"))
	//movieCollector.URLFilters = []*regexp.Regexp{regexp.MustCompile("https://movie\\.douban\\.com/subject/[0-9]{8,10}/$")}
	//celebrityCollector := c.Clone()

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

	// 文件存储
	fName := "douban_movie.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("创建文件失败 %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write CSV header
	writer.Write([]string{"名称", "年份", "详情", "评星", "链接"})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response", r.StatusCode, string(r.Body))
	})

	// 遍历页面中的链接并访问
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		targetURL := e.Request.AbsoluteURL(link)
		if strings.Contains(targetURL, "/subject/") {
			movieCollector.Visit(targetURL)
		}

		// 访问链接，使用e.Request.Visit会记录深度
		//c.Visit(e.Request.AbsoluteURL(link))
		e.Request.Visit(link)
		//fmt.Printf("Link found end: %q -> %s\n", e.Text, link)
	})

	movieCollector.OnHTML("div#content", func(e *colly.HTMLElement) {
		selection := e.DOM
		idx := selection.Find("div.top250 > span.top250-no").Text()
		title := selection.Find("h1 > span").First().Text()
		year := selection.Find("h1 > span.year").Text()
		info := selection.Find("div#info").Text()
		info = strings.ReplaceAll(info, " ", "")
		info = strings.ReplaceAll(info, "\n", "; ")
		rating := selection.Find("strong.rating_num").Text()
		movie := collymovie.Movie{
			Idx:    idx,
			Title:  title,
			Year:   year,
			Info:   info,
			Rating: rating,
			URL:    e.Request.URL.String(),
		}
		writer.Write([]string{
			title,
			year,
			info,
			rating,
			e.Request.URL.String(),
		})
		fmt.Printf("Movie found: %+v", movie)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// Start scraping on https://hackerspaces.org
	//c.Visit("https://movie.douban.com/top250")
	c.Visit("https://www.baidu.com")
	c.Wait()
	fmt.Printf("花费时间:%s", time.Since(t))
}
