package main

type Product struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url"`
	Images      []struct {
		Url string `json:"url"`
	} `json:"images,omitempty"`
	Options  []*Option `json:"option_names,omitempty"`
	Variants []*ProductOption
}

type Option struct {
	Title string `json:"title"`
}

type ProductOption struct {
	SKU       string `json:"sku"`
	Available bool   `json:"available"`
	Price     int    `json:"price"`
	Amount    int    `json:"amount"`
}
