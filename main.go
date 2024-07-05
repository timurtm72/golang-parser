package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/xuri/excelize/v2"
)

type Link struct {
	Id  int
	Url string
}

// WebScrap выполняет веб-скрапинг и передает результаты в канал
func WebScrap(url string, ch chan<- []Link, wg *sync.WaitGroup) {
	defer wg.Done()

	var links []Link
	var idx int

	// Извлекаем домен из URL
	newUrl := url[8:] // Возможно, лучше использовать функцию для извлечения домена

	// Определяем домен
	domainParts := strings.Split(newUrl, "/")
	domain := domainParts[0]

	c := colly.NewCollector(
		colly.AllowedDomains(domain),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		if r.URL.String() != "" {
			links = append(links, Link{idx + 1, r.URL.String()})
			idx++
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	// Асинхронный визит
	c.Visit(url)

	// Передаем результаты в канал
	ch <- links
}

// SaveToFile сохраняет данные в файл Excel
func SaveToFile(filename string, internalLinks []Link) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for i, row := range internalLinks {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%v", i+1), strconv.Itoa(row.Id))
		f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i+1), row.Url)
	}

	if err := f.SaveAs(filename); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Data saved to file...")
	}
}

func main() {
	fmt.Println("Идет загрузка...")
	// Создаем канал для передачи данных
	ch := make(chan []Link)

	// Создаем WaitGroup для ожидания завершения горутины
	var wg sync.WaitGroup
	wg.Add(1)

	// Запускаем веб-скрапинг в горутине
	go WebScrap("https://technocom.site123.me", ch, &wg)

	// Ожидаем завершения горутины и закрытия канала
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Ожидаем данные из канала
	var links []Link
	for l := range ch {
		links = append(links, l...)
	}

	// Сохраняем данные в файл
	SaveToFile("data.xlsx", links)

	// Печатаем полученные ссылки
	printLinks(links)
}

// Функция для печати ссылок
func printLinks(links []Link) {
	for _, link := range links {
		fmt.Printf("%d => %s\n", link.Id, link.Url)
	}
}
