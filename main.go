package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"path"
	"time"
)

const Purpose = "https://barkingstore.ru"

var fileName string = "quote.csv"

func main() {
	file, err := os.Create(path.Join("xxx", fileName))
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

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"))
	c.SetRequestTimeout(5 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("%s %d %s",
			r.Request.URL.String(),
			r.StatusCode,
			err)
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			log.Printf("%s %d", r.Request.URL.String(), r.StatusCode)
		}
	})

	c.OnHTML(".index_coll-text", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.ChildText(".h4-like")

		err = writer.Write([]string{
			e.Request.URL.JoinPath(link).String(),
			title,
		})
		if err != nil {
			fmt.Printf("can't handle %q %q\n", link, title)
		}
	})

	err = c.Visit(Purpose)
	if err != nil {
		log.Fatal(err)
	}
}
