package main

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {

	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	// run parse pool
	locs := productLocations[252:255]

	runParse(locs, 2)
	// running parse and getting product data
	// products, _ := runParse(locs, 2)
	// getting xml document output
	// getXMLDocument(products)

}
