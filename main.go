package main

import "fmt"

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {

	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	// run parse pool
	locs := productLocations[250:255]
	products, _ := runParse(locs, 3)

	for _, product := range products {
		getXMLProduct(product)
		fmt.Println(" ")
		fmt.Println(" ")
	}
}
