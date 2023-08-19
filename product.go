package main

type Product struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url"`
	Images      []struct {
		Url string `json:"url"`
	} `json:"images,omitempty"`
	Options []struct {
		Title string `json:"title"`
	} `json:"option_names,omitempty"`
	Variants []ProductOption
}

//todo move option to separate struct

type ProductOption struct {
	SKU       string `json:"sku"`
	Available bool   `json:"available"`
	Price     int    `json:"price"`
	Amount    int    `json:"amount"`
}
