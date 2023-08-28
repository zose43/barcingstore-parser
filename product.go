package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type Goods struct {
	Status string     `json:"status"`
	Items  []*Product `json:"products"`
}

type Product struct {
	Title       string `json:"title" xml:"title"`
	Description string `json:"short_description,omitempty" xml:"short_description,omitempty"`
	Url         string `json:"url" xml:"url"`
	Images      []*Images
	Options     []Option         `json:"option_names,omitempty" xml:"option_names,omitempty"`
	Variants    []*ProductOption `json:"variants,omitempty" xml:"variants,omitempty"`
}

func (p *Product) cleanDescription() {
	reader := strings.NewReader(p.Description)
	sel, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Printf("can't clean description %s\n", err)
	}
	p.Description = sel.Text()
}

type Images struct {
	Id       int    `json:"id" xml:"id"`
	Url      string `json:"original_url" xml:"original_url"`
	Position int    `json:"position" xml:"position"`
}

type Option struct {
	Title string `json:"title" xml:"title"`
}

type ProductOption struct {
	Title     string `json:"title" xml:"title"`
	SKU       string `json:"sku" xml:"SKU"`
	Available bool   `json:"available" xml:"available"`
	Price     string `json:"price" xml:"price"`
	Amount    int    `json:"quantity" xml:"quantity"`
	ImageID   int    `json:"image_id" xml:"image_id"`
}
