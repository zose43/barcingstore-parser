package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

const Purpose = "https://barkingstore.ru"

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
	c.OnRequest(requesting)
	c.OnError(processError)
	c.OnResponse(response)
	// handle menu pages
	c.OnHTML(".subcol1 .menu-link", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		err := c.Visit(nextPage)
		visitError(err, e.Request.URL)
	})
	// create product-json endpoints list
	c.OnHTML(".product-item form", products)

	err = c.Visit(Purpose)
	if err != nil {
		log.Fatal(err)
	}
}

func getCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
		colly.AllowedDomains("barkingstore.ru"))
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
	if val := r.Headers.Get("Content-Type"); strings.Contains(val, "application/json") {
		handleProductJSON(r.Body)
	}
}

func convertTo() string {
	format := flag.String("f", "csv", "format of document")
	return *format
}

func products(e *colly.HTMLElement) {
	const endpoint = "https://barkingstore.ru/products_by_id"
	id := e.Attr("data-product-id")
	prodPage := fmt.Sprintf("%s/%s.json", endpoint, id)

	c := getCollector()
	err := c.Visit(prodPage)
	visitError(err, e.Request.URL)
}

func filename(format string) string {
	return path.Join("xxx", fmt.Sprintf("export.%s", format))
}

func visitError(err error, u *url.URL) {
	if err != nil {
		fmt.Printf("can't visit %s %s\n", u, err)
	}
}

func handleProductJSON(data []byte) {
	var prods []*Product
	err := json.Unmarshal(data, &prods)
	if err != nil {
		fmt.Printf("can't unmarshall data %s\n", err)
		return
	}

	fmt.Printf("done, %d items\n", len(prods))
}
