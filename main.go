package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const Purpose = "https://barkingstore.ru"

var prods []*Product

func main() {
	ft := convertTo()
	file, err := os.Create(filename(ft))
	if err != nil {
		log.Fatal("can't create file")
	}
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	err = writer.Write([]string{"Url", "Title"})
	if err != nil {
		log.Print("can't write csv file")
	}
	defer writer.Flush()

	c := getCollector()
	c2 := clone()

	c.OnRequest(requesting)
	c.OnError(processError)
	c.OnResponse(response)
	// handle menu pages
	c.OnHTML(".subcol1 .menu-link", func(e *colly.HTMLElement) {
		grabCategories(e, c2)
	})
	c2.OnError(processError)
	// create product-json endpoints list and handle json
	c2.OnHTML(".product-item form", func(e *colly.HTMLElement) {
		grabProducts(e, c2)
	})
	c2.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			fmt.Printf("%s %d", r.Request.URL, r.StatusCode)
		}
		if val := r.Headers.Get("Content-Type"); strings.Contains(val, "application/json") {
			handleProductJSON(r.Body)
		}
	})

	err = c.Visit(Purpose)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("done, %d items\n", len(prods))
}

func getCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
		colly.AllowedDomains("barkingstore.ru"),
		colly.MaxDepth(1))
	c.SetRequestTimeout(8 * time.Second)
	return c
}

func requesting(r *colly.Request) {
	fmt.Println("Visiting:", r.URL)
}

func processError(r *colly.Response, err error) {
	fmt.Printf("%s %d %s",
		r.Request.URL,
		r.StatusCode,
		err)
}

func response(r *colly.Response) {
	if r.StatusCode != 200 {
		fmt.Printf("%s %d", r.Request.URL, r.StatusCode)
	}
}

func convertTo() string {
	format := flag.String("f", "csv", "format of document")
	return *format
}

func grabProducts(e *colly.HTMLElement, c *colly.Collector) {
	const endpoint = "https://barkingstore.ru/products_by_id"
	id := e.Attr("data-product-id")
	prodPage := fmt.Sprintf("%s/%s.json", endpoint, id)

	err := c.Visit(prodPage)
	visitError(err, e.Request.URL)
}

func filename(format string) string {
	return path.Join("xxx", fmt.Sprintf("export.%s", format))
}

func visitError(err error, u any) {
	if err != nil {
		fmt.Printf("can't visit %s %s\n", u, err)
	}
}

func handleProductJSON(data []byte) {
	//todo unmarshall empty struct
	var prod Product
	err := json.Unmarshal(data, &prod)
	if err != nil {
		fmt.Printf("can't unmarshall data %s\n", err)
		return
	}
	prods = append(prods, &prod)
}

func grabCategories(e *colly.HTMLElement, c *colly.Collector) {
	//todo each iterate only 1 item
	var cats []string
	e.DOM.Each(func(_ int, s *goquery.Selection) {
		cat := e.Request.AbsoluteURL(s.AttrOr("href", ""))
		cats = append(cats, cat)
	})
	for _, cat := range cats {
		if len(cat) > 0 {
			err := c.Visit(cat)
			visitError(err, cat)
		}
	}
}

func clone() *colly.Collector {
	c := getCollector()
	return c.Clone()
}
