package main

import (
	"encoding/json"
	"encoding/xml"
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

var (
	prods      []*Product
	catCounter uint
	format     string
)

func main() {
	//todo download images, output dir flag
	mustConvertTo()
	fname := filename()
	file, err := os.Create(fname)
	if err != nil {
		log.Fatal("can't create export file")
	}
	defer func() { _ = file.Close() }()

	c := getCollector()
	c2 := clone()

	c.OnError(processError)
	c.OnResponse(response)
	// handle menu pages
	c.OnHTML(".header-middle .subcol1", func(e *colly.HTMLElement) {
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
			handleJSON(r.Body)
		}
	})

	err = c.Visit(Purpose)
	if err != nil {
		log.Fatal(err)
	}

	err = exportFile(file)
	if err != nil {
		fmt.Printf("can't export document %s", err)
	}
	fmt.Printf("done, %d processed categories, %d saved products\n", catCounter, len(prods))
}

func getCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
		colly.AllowedDomains("barkingstore.ru"),
		colly.MaxDepth(1))
	c.SetRequestTimeout(8 * time.Second)
	return c
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

func mustConvertTo() {
	flag.StringVar(&format, "f", "json", "document format")
	flag.Parse()
	if format != "json" && format != "xml" {
		log.Fatalf("not suppored document format %s", format)
	}
}

func grabProducts(e *colly.HTMLElement, c *colly.Collector) {
	const endpoint = "https://barkingstore.ru/products_by_id"
	id := e.Attr("data-product-id")
	prodPage := fmt.Sprintf("%s/%s.json", endpoint, id)

	repeat, _ := c.HasVisited(prodPage)
	if !repeat {
		err := c.Visit(prodPage)
		visitError(err, e.Request.URL)
	}
}

func filename() string {
	dt := time.Now().Format(time.DateOnly)
	return path.Join("export", fmt.Sprintf("export-%s.%s", dt, format))
}

func visitError(err error, u any) {
	if err != nil {
		fmt.Printf("can't visit %s %s\n", u, err)
	}
}

func handleJSON(data []byte) {
	var goods Goods
	err := json.Unmarshal(data, &goods)
	if err != nil {
		fmt.Printf("can't unmarshall data %s\n", err)
		return
	}

	if goods.Status != "ok" {
		fmt.Printf("invalid JSON status")
		return
	}
	for _, prod := range goods.Items {
		prod.cleanDescription()
		prods = append(prods, prod)
	}
}

func grabCategories(e *colly.HTMLElement, c *colly.Collector) {
	var subcats []string
	e.DOM.Find(".menu-link").
		Each(func(_ int, s *goquery.Selection) {
			cat := e.Request.AbsoluteURL(s.AttrOr("href", ""))
			subcats = append(subcats, cat)
		})
	for _, cat := range subcats {
		if len(cat) > 0 {
			fmt.Println("visiting category:", cat)
			err := c.Visit(cat)
			visitError(err, cat)
			catCounter++
		}
	}
}

func clone() *colly.Collector {
	c := getCollector()
	return c.Clone()
}

func exportFile(file *os.File) error {
	switch format {
	case "xml":
		return saveInXML(file)
	case "json":
		return saveInJSON(file)
	}
	return nil
}

func saveInXML(f *os.File) error {
	rawData, err := xml.MarshalIndent(prods, "", " ")
	if err != nil {
		return err
	}

	_, _ = f.WriteString("<Products>")
	_, err = f.Write(rawData)
	if err != nil {
		return err
	}
	_, _ = f.WriteString("</Products>")
	return nil
}

func saveInJSON(f *os.File) error {
	rawData, err := json.Marshal(prods)
	if err != nil {
		return err
	}

	_, err = f.Write(rawData)
	if err != nil {
		return err
	}
	return nil
}
