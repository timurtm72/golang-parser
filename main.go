package main

import (
	"fmt"
	"log"
)

func main() {
	//url := "https://www.yandex.ru"
	url, _ := readUrl()
	doc, err := fetchPage(url)
	if err != nil {
		log.Fatalf("Error fetching page: %v", err)
	}
	if validateUrl(url) {
		links := extractLinks(doc)
		for _, link := range links {
			fmt.Println(link)
		}
	} else {
		log.Fatalf("Error validating URL: %v", url)
	}
}
