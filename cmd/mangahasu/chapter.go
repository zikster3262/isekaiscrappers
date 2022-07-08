package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func scrapeImages(url string) {

	// func main(url string) {

	dir := "manga"

	website := "https://mangahasu.se"
	findex := indexAt(url, "/", len(website))
	lindex := indexAt(url, "/", len("https://mangahasu.se"+"/"))

	folder := url[findex+1 : lindex]
	downFold := dir + "/" + folder
	i := strings.Index(url, "chapter")
	ix := len(url) - len(".html")
	chapter := url[i:ix]

	chapterfolder := dir + "/" + folder + "/" + chapter

	if err := os.MkdirAll(downFold, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(chapterfolder, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.SetRequestTimeout(120 * time.Second)

	c.OnHTML("#loadchapter img", func(e *colly.HTMLElement) {
		ch := e.Attr("src")
		path := chapterfolder + "/" + ch[len(ch)-7:]

		DownloadFile(path, ch)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("User-Agent", "Mozilla/5.0")
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(url)
}

func indexAt(s, sep string, n int) int {
	idx := strings.Index(s[n:], sep)
	if idx > -1 {
		idx += n
	}
	return idx
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
