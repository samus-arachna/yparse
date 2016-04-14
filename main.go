package main

import "fmt"

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {
	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	// run parse pool
	locs := productLocations[250:255]
	products, count := runParse(locs, 3)

	for _, e := range products {
		fmt.Println(e)
		fmt.Println("")
	}
	fmt.Print("Parsed products: ")
	fmt.Println(count)
	/*
		link := "http://www.yves-rocher.ru/c/kol-e-i-braslet--beskonecnost--/p/yr.rukz.RN81907"
		product, err := parseProduct(link, true)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(product)
	*/
}
