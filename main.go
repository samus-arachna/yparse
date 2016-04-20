package main

import (
	"fmt"
	"log"
)

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {

	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	if len(productLocations) == 0 {
		log.Fatal("There is no locations of products to parse.")
	}

	// initiating logging, removing old logs
	logInit()

	// run parse pool
	fmt.Println("Parsing started.")
	locs := productLocations // full parse

	// running parse and getting product data
	products, categories, _ := runParse(locs, 3)
	// getting xml document output
	getXMLDocument(products, categories)
}
