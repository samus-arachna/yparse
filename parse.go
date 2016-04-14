package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func runParse(products []string, connections int) []map[string]string {
	// stats for debug
	overall := len(products)
	real := 0

	// get first slice == number of connections
	parsed := []map[string]string{}
	pool := products[0:connections]
	products = products[connections:]
	var wg sync.WaitGroup

	// run connections
	for len(products) > 0 {
		wg.Add(len(pool))

		for _, product := range pool {
			real++
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
			if reflect.DeepEqual(pool, products) {
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

	fmt.Print("overall products to parse: ")
	fmt.Println(overall)
	fmt.Print("real number of products parsed: ")
	fmt.Println(real)

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

	// seeking product code (id)
	productCodeWrap := doc.Find(".ref").Text()
	if len(productCodeWrap) == 0 {
		return nil, errors.New("No product code wrap was found")
	}
	productCode := parseCode(productCodeWrap)
	product["code"] = productCode
	// END seeking product code (id)

	// seeking product title (name)
	productTitle := doc.Find(".product_overview > h1").Text()
	if len(productTitle) == 0 {
		return nil, errors.New("No product title was found.")
	}
	product["title"] = productTitle
	// END seeking product title (name)

	// seeking product description in two places
	productDesc := doc.Find(".product_overview .baseline a").Text()
	if len(productDesc) == 0 {
		productDesc, err = doc.Find(".txt_2column p").Html()
		if err != nil {
			return nil, errors.New("No product description was found.")
		}
		productDesc = strings.Split(productDesc, "<br")[0]
		if len(productDesc) == 0 {
			return nil, errors.New("No second product description was found.")
		}
	}
	product["desc"] = productDesc
	// END seeking product description in two places

	// seeking product preview image
	productImg, _ := doc.Find("img#product_slider_image").Attr("src")
	if len(productImg) == 0 {
		return nil, errors.New("No product image was found.")
	}
	product["img"] = productImg
	// END seeking product preview image

	// seeking product main (current) price
	productCurrentPriceWrap := doc.Find(".product_overview .price").Text()
	if len(productCurrentPriceWrap) == 0 {
		return nil, errors.New("No product current price was found.")
	}
	productCurrentPrice := parsePrice(productCurrentPriceWrap)
	product["price"] = productCurrentPrice
	// END seeking product main (current) price

	// seeking product old price
	productOldPriceWrap := doc.Find(".product_overview .striped_price").Text()
	productOldPrice := parsePrice(productOldPriceWrap)
	product["priceOld"] = productOldPrice
	// END seeking product old price

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
	code := re.FindAllString(wrapSplitted[0], -1)
	if len(code) == 0 {
		return "0"
	}

	return code[0]
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
