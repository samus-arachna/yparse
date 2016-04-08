package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/html"
)

const sitemapLocation = "http://www.yves-rocher.ru/sitemap.xml"

func main() {
	sitemapLocations := getLocations(sitemapLocation)
	productLocations := getProductLocations(sitemapLocations)

	parseProduct(productLocations[251])
}

// TODO
// parse single product should return a map(?)
func parseProduct(productURL string) {
	doc, err := goquery.NewDocument(productURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("url: " + productURL)

	productCodeWrap := doc.Find(".ref").Text()
	productCode := parseCode(productCodeWrap)
	fmt.Print("code: ")
	fmt.Println(productCode)

	productTitle := doc.Find(".product_overview > h1").Text()
	fmt.Println("title: " + productTitle)

	productDesc := doc.Find(".product_overview .baseline a").Text()
	fmt.Println("desc: " + productDesc)

	productImg, _ := doc.Find("img#product_slider_image").Attr("src")
	fmt.Println("img: " + productImg)

	productCurrentPriceWrap := doc.Find(".inside .price").Text()
	productCurrentPrice := parsePrice(productCurrentPriceWrap)
	fmt.Println("price: " + productCurrentPrice)

	productOldPriceWrap := doc.Find(".inside .striped_price").Text()
	productOldPrice := parsePrice(productOldPriceWrap)
	fmt.Println("old price: " + productOldPrice)
}

// TODO
// parse product price out of string
func parsePrice(wrap string) string {
	priceTrimmed := strings.TrimSpace(wrap)
	re := regexp.MustCompile("[0-9]+([,.][0-9]+)")
	numbers := re.FindAllString(priceTrimmed, -1)
	price := strings.Join(numbers, "")

	return price
}

// parse product code out of string
func parseCode(wrap string) int {
	wrapTrimmed := strings.TrimSpace(wrap)
	wrapSplitted := strings.Split(wrapTrimmed, "-")

	re := regexp.MustCompile("[0-9]+")
	code, err := strconv.Atoi(re.FindAllString(wrapSplitted[0], -1)[0])
	if err != nil {
		return 0
	}

	return code
}

// getting all product locations
func getProductLocations(xmlData []string) []string {
	locations := []string{}

	for _, loc := range xmlData {
		locations = append(locations, getLocations(loc)...)
	}

	productLocations := []string{}

	for _, loc := range locations {
		if strings.Contains(loc, "/p/") {
			productLocations = append(productLocations, loc)
		}
	}

	return productLocations
}

// getting locations from xml
func getLocations(xmlURL string) []string {
	resp, err := http.Get(xmlURL)
	if err != nil {
		log.Fatal(err)
	}

	locations := []string{}
	token := html.NewTokenizer(resp.Body)
	depth := 0

	for {
		tt := token.Next()
		switch tt {
		case html.ErrorToken:
			resp.Body.Close()
			return locations
		case html.TextToken:
			if depth > 0 {
				locations = append(locations, string(token.Text()))
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := token.TagName()
			if string(tn) == "loc" {
				if tt == html.StartTagToken {
					depth++
				} else {
					depth--
				}
			}
		}
	}
}
