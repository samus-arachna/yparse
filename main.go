package main

import (
	"fmt"
	"sync"
)

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {
	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	/*
		fmt.Println(productLocations[258])
		fmt.Println(" --- ")
		product, err := parseProduct(productLocations[258])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(product)
	*/

	// run parse pool
	products := productLocations[250:255]
	runParse(products, 2)
}

func runParse(products []string, connections int) []map[string]string {
	// get first slice == number of connections
	parsed := []map[string]string{}
	pool := products[0:connections]
	products = products[connections:]
	var wg sync.WaitGroup

	// run connections
	for len(products) > 0 {
		wg.Add(len(pool))

		for _, product := range pool {
			go func(product string) {
				parsedProduct, err := parseProduct(product)
				if err != nil {
					fmt.Println(err.Error() + " on link " + product)
					fmt.Println("")
					wg.Done()
					return
				}
				parsed = append(parsed, parsedProduct)
				defer wg.Done()
			}(product)
		}

		wg.Wait()

		if len(products) > connections {
			pool = products[0:connections]
			products = products[connections:]
		} else {
			if len(pool) == len(products) {
				products = []string{}
			} else {
				pool = products[:]
			}
		}
	}

	for _, item := range parsed {
		fmt.Println(item)
		fmt.Println("")
	}

	return parsed
}
