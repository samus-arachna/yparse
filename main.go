package main

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
