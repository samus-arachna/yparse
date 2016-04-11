package main

import (
	"fmt"
	"sync"
)

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {
	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	fmt.Println(productLocations[253])
	product := parseProduct(productLocations[253])
	fmt.Println(product)

	// TODO parse first 15 products in a series of 5
	//products := productLocations[250:260]
	//runParse(products, 5)
}

func runParse(products []string, connections int) {
	// get first slice == number of connections
	parsed := []map[string]string{}
	pool := products[0:connections]
	products = products[connections:]
	var wg sync.WaitGroup

	// run connections
	fmt.Println(len(pool))
	wg.Add(len(pool))
	for _, product := range pool {
		go func(product string) {
			parsed = append(parsed, parseProduct(product))
			defer wg.Done()
		}(product)
	}
	wg.Wait()

	for _, item := range parsed {
		fmt.Println(item)
		fmt.Println("")
	}

	// get another slice, until no more products

	// signalize complete
}
