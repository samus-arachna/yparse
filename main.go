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
	locs := productLocations[595:597] // partly parse
	//locs := productLocations // full parse

	// running parse and getting product data
	products, categories, _ := runParse(locs, 1)

	// DEBUG ON
	for _, v := range categories {
		fmt.Println(v)
	}
	fmt.Println("")
	for _, p := range products {
		fmt.Println(p)
	}
	// DEBUG OFF

	// getting xml document output
	getXMLDocument(products, categories)
}
