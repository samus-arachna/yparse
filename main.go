package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {
	siteMapXML := getMainSitemap(sitemapLocation)
	fmt.Println(siteMapXML)
}

// TODO
func getSecondarySitemaps() {

}

// getting main sitemap
func getMainSitemap(xmlURL string) string {
	resp, err := http.Get(xmlURL)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()

	return string(bytes)
}
