package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type MangaLink struct {
	ID      int    `json:"id"`
	Chapter string `json:"chapter"`
	URL     string `json:"url"`
}

func scrapeChapters(url string, ch chan<- []MangaLink) {
	var links []MangaLink
	var m MangaLink

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.SetRequestTimeout(120 * time.Second)

	c.OnHTML(".list-chapter table tbody tr td", func(e *colly.HTMLElement) {
		m.Chapter = e.Text
		m.URL = e.ChildAttr("a", "href")
		m.ID = e.Index
		links = append(links, m)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error - chapters:", e)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("User-Agent", "Mozilla/5.0")
	})

	c.Visit(url)
	ch <- links
}
