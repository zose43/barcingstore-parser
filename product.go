package main

type Goods struct {
	Status string     `json:"status"`
	Items  []*Product `json:"products"`
}

type Product struct {
	Title       string `json:"title"`
	Description string `json:"short_description,omitempty"`
	Url         string `json:"url"`
	Images      []*Images
	Options     []Option         `json:"option_names,omitempty"`
	Variants    []*ProductOption `json:"variants,omitempty"`
}

type Images struct {
	Url      string `json:"original_url"`
	Position int    `json:"position"`
}

type Option struct {
	Title string `json:"title"`
}

type ProductOption struct {
	Title     string `json:"title"`
	SKU       string `json:"sku"`
	Available bool   `json:"available"`
	Price     string `json:"price"`
	Amount    int    `json:"quantity"`
	ImageID   int    `json:"image_id"`
}
