package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

var url string = "https://www.idnes.cz/"

type News struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

func main() {
	allNews := make([]News, 0)

	c := colly.NewCollector()

	var n News
	c.SetRequestTimeout(120 * time.Second)

	c.OnHTML("#art-top2 .art", func(e *colly.HTMLElement) {
		n.Description = e.ChildText(".art-link h3")
		n.Link = e.ChildAttr(".art-link", "href")
		n.ID = e.Index
		allNews = append(allNews, n)

	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)
	for i, txt := range allNews {
		fmt.Println(i, txt)
	}
}
