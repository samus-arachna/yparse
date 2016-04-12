package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

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

// parse single product
func parseProduct(productURL string) (map[string]string, error) {
	doc, err := goquery.NewDocument(productURL)
	if err != nil {
		log.Fatal(err)
	}

	product := map[string]string{}

	product["url"] = productURL

	productCodeWrap := doc.Find(".ref").Text()
	if len(productCodeWrap) == 0 {
		return nil, errors.New("No product code wrap was found")
	}
	productCode := parseCode(productCodeWrap)
	product["code"] = productCode

	productTitle := doc.Find(".product_overview > h1").Text()
	if len(productTitle) == 0 {
		return nil, errors.New("No product title was found.")
	}
	product["title"] = productTitle

	productDesc := doc.Find(".product_overview .baseline a").Text()
	if len(productDesc) == 0 {
		return nil, errors.New("No product description was found.")
	}
	product["desc"] = productDesc

	productImg, _ := doc.Find("img#product_slider_image").Attr("src")
	if len(productImg) == 0 {
		return nil, errors.New("No product image was found.")
	}
	product["img"] = productImg

	productCurrentPriceWrap := doc.Find(".product_overview .price").Text()
	if len(productCurrentPriceWrap) == 0 {
		return nil, errors.New("No product current price was found.")
	}
	productCurrentPrice := parsePrice(productCurrentPriceWrap)
	product["price"] = productCurrentPrice

	productOldPriceWrap := doc.Find(".product_overview .striped_price").Text()
	productOldPrice := parsePrice(productOldPriceWrap)
	product["priceOld"] = productOldPrice

	return product, nil
}

// parse product price out of string
func parsePrice(wrap string) string {
	priceTrimmed := strings.TrimSpace(wrap)
	re := regexp.MustCompile("([0-9]+)|([,.][0-9]+)")
	numbers := re.FindAllString(priceTrimmed, -1)
	price := strings.Join(numbers, "")

	return price
}

// parse product code out of string
func parseCode(wrap string) string {
	wrapTrimmed := strings.TrimSpace(wrap)
	wrapSplitted := strings.Split(wrapTrimmed, "-")

	re := regexp.MustCompile("[0-9]+")
	code := re.FindAllString(wrapSplitted[0], -1)[0]
	if len(code) == 0 {
		return "0"
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
		if strings.Contains(loc, "/p/") && !strings.Contains(loc, "//p/") {
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
