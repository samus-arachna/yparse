package main

import "fmt"

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {
	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	product := parseProduct(productLocations[251])

	// TODO parse first 50 products in a series of 5
	fmt.Println(product)
}
