package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// fetchPage fetches the HTML document from the given URL
func fetchPage(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("error to fetchPage: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: %s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Errorf("error to parse page: %s", err)
		return nil, err
	}

	return doc, nil
}

// extractLinks extracts all the href attributes from <a> tags
func extractLinks(node *html.Node) []string {
	var links []string
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
				break
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		links = append(links, extractLinks(child)...)
	}
	return links
}

func readUrl() (string, error) {
	fmt.Print("Enter url to format https://www.google.com =>>")
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	url, err := reader.ReadString('\n')
	if err != nil {
		errStr := "An error occured while reading input. Please try again"
		fmt.Println(errStr, err)
		return errStr, err
	} else {
		// remove the delimeter from the string
		url = strings.TrimSuffix(url, "\n")
		fmt.Errorf("URL: %s", url)
	}
	return url, nil
}

func validateUrl(url string) bool {
	// Регулярное выражение для проверки формата URL
	// Оно проверяет, что строка начинается с "www.", затем идет имя домена (состоящее из букв, цифр или дефисов),
	// и заканчивается точкой и доменным расширением (состоящим из букв)
	regex := `^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`
	re := regexp.MustCompile(regex)

	// Проверяем, соответствует ли URL регулярному выражению
	if re.MatchString(url) {
		fmt.Println("URL is in the correct format.")
		return true
	} else {
		fmt.Println("URL is not in the correct format.")
		return false
	}
}
