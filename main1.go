package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
)

type Link struct {
	Id       int
	LinkPath string
}

var idx = 0
var links []Link

// WebScrap =============================================================================================
func WebScrap(url string) {
	newUrl := url[8:]
	c := colly.NewCollector(
		colly.AllowedDomains(newUrl)) //"technocom.site123.me"))

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	})

	c.OnRequest(func(r *colly.Request) {
		if r.URL.String() != "" {
			links = append(links, Link{idx + 1, r.URL.String()})
			fmt.Println(links[idx])
			idx++
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	c.OnResponse(func(r *colly.Response) {
	})
	c.Visit(url) //"https://technocom.site123.me/")
}

// -----------------------------------------------------------------------------------------
func SaveToFile(filename string, InternalLinks []Link) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for i, row := range InternalLinks {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%v", i+1), strconv.Itoa(row.Id))
		f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i+1), row.LinkPath)
		idx++
	}
	if err := f.SaveAs(filename); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Data saved to file...")
	}
	InternalLinks = nil
}

// -----------------------------------------------------------------------------------------
func (l Link) String() string {
	return fmt.Sprintf("{Id:%d, Link:%s}", l.Id, l.LinkPath)
}

// -----------------------------------------------------------------------------------------
func printLinks(links []Link) {
	for _, link := range links {
		fmt.Printf("%d => %s\n", link.Id, link.LinkPath)
	}
}

// -----------------------------------------------------------------------------------------
func main() {
	//=====================================================================================
	fmt.Println("Идет загрузка...")
	WebScrap("https://technocom.site123.me")
	SaveToFile("data.xlsx", links)
	//=====================================================================================

}
