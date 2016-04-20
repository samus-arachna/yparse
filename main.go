package main

import "log"

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {

	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	if len(productLocations) == 0 {
		log.Fatal("There is no locations of products to parse.")
	}

	// run parse pool
	// locs := productLocations[252:255] // example partly parsing
	locs := productLocations // full parse

	// running parse and getting product data
	products, categories, _ := runParse(locs, 3)
	// getting xml document output
	getXMLDocument(products, categories)
}
