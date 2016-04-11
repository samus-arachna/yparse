package main

import "fmt"

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {
	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	//product := parseProduct(productLocations[251])

	// TODO parse first 15 products in a series of 3
	products := productLocations[0:15]
	runParse(products, 3)
}

func runParse(products []string, connections int) {
	// get first slice == number of connections
	pool := products[0:connections]
	products = products[connections:]

	fmt.Println(len(pool))
	fmt.Println(len(products))

	// run connections

	// wait till all complete

	// get another slice, until no more products

	// signalize complete
}
