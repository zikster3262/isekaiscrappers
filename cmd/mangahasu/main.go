package main

import "time"

// Manga Site

func main() {

	var urls = []string{"https://mangahasu.se/"}

	ch := make(chan []Manga)
	mch := make(chan []MangaLink)

	for _, url := range urls {
		go scrapeMangaSite(url, ch)
	}

	for _, v := range <-ch {
		go scrapeChapters(v.URL, mch)
	}
	close(ch)

	var list []string
	for {
		msg := <-mch
		for _, v := range msg {
			if v.URL != "" {
				list = append(list, v.URL)
			}

		}
		for _, web := range list {
			scrapeImages(web)
			time.Sleep(3 * time.Second)
		}
	}

}
