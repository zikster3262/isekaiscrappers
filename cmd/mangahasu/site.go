package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type Manga struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func scrapeMangaSite(url string, ch chan<- []Manga) {

	collection := make([]Manga, 0)
	var m Manga

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.SetRequestTimeout(120 * time.Second)

	c.OnHTML(".list-top-of-week li .div_item", func(e *colly.HTMLElement) {
		m.Description = e.ChildText(".info-manga h3")
		m.ID = e.Index
		m.URL = e.ChildAttr(".name-manga", "href")
		collection = append(collection, m)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error - site:", e)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("User-Agent", "Mozilla/5.0")
	})

	c.Visit(url)
	ch <- collection
}
