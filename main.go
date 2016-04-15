package main

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {

	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	// run parse pool
	locs := productLocations[250:254]
	products, _ := runParse(locs, 2)

	getXMLDocument(products)
}
