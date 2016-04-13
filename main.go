package main

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {
	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	// run parse pool
	products := productLocations[250:260]
	runParse(products, 3)
}
