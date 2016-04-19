package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type category struct {
	id       string
	parentID string
	name     string
}

func runParse(products []string, connections int) ([]map[string]string, int) {
	// how much products was parsed
	count := 0

	// get first slice == number of connections
	parsed := []map[string]string{}
	pool := products[0:connections]
	products = products[connections:]

	// init categories
	categories := map[string]category{}

	// syncing connections
	var wg sync.WaitGroup

	// run connections
	for len(products) > 0 {
		wg.Add(len(pool))

		for _, product := range pool {
			count++
			go func(product string) {
				parsedProduct, err := parseProduct(product, true, &categories)
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

	// outputting parsed products
	/*
		for _, item := range parsed {
			fmt.Println(item)
			fmt.Println("")
		}
	*/

	// printing out categories
	fmt.Println(categories)
	fmt.Println("")

	return parsed, count
}

// parse single product
func parseProduct(productURL string, fromURL bool, categories *map[string]category) (map[string]string, error) {
	doc := getDocumentType(productURL, fromURL)

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
	productTitle := strings.TrimSpace(doc.Find(".product_overview > h1").Text())
	if len(productTitle) == 0 {
		return nil, errors.New("No product title was found.")
	}
	product["title"] = productTitle
	// END seeking product title (name)

	// seeking product description in two places
	productDesc := strings.TrimSpace(doc.Find(".product_overview .baseline a").Text())
	if len(productDesc) == 0 {
		productDescSecond, err := doc.Find(".txt_2column p").Html()
		if err != nil {
			return nil, errors.New("No product description was found.")
		}

		productDescSecond = strings.Split(productDescSecond, "<br")[0]
		if len(productDescSecond) == 0 {
			return nil, errors.New("No second product description was found.")
		}

		product["desc"] = productDescSecond
	} else {
		product["desc"] = productDesc
	}
	// END seeking product description in two places

	// seeking product preview image
	productImg, _ := doc.Find("img#product_slider_image").Attr("src")
	if len(productImg) == 0 {
		return nil, errors.New("No product image was found.")
	}
	product["img"] = productImg
	// END seeking product preview image

	// seeking product main (current) price
	productCurrentPriceWrap := strings.TrimSpace(doc.Find(".product_overview .price").Text())
	if len(productCurrentPriceWrap) == 0 {
		return nil, errors.New("No product current price was found.")
	}
	productCurrentPrice := parsePrice(productCurrentPriceWrap)
	product["price"] = productCurrentPrice
	// END seeking product main (current) price

	// seeking product old price
	productOldPriceWrap := strings.TrimSpace(doc.Find(".product_overview .striped_price").Text())
	productOldPrice := parsePrice(productOldPriceWrap)
	product["priceOld"] = productOldPrice
	// END seeking product old price

	// parsing single category to product
	cat := doc.Find(".crumbs a").Eq(-2)
	catName := cat.Text()
	product["categoryName"] = catName
	catHref, _ := cat.Attr("href")
	catID := parseCategoryID(catHref)
	product["categoryID"] = catID
	// END parsing single category to product

	// TODO parsing category tree
	parseCategory(doc, categories)
	fmt.Println("")
	// END parsing category tree

	return product, nil
}

// TODO parse category
func parseCategory(doc *goquery.Document, categories *map[string]category) {
	sel := doc.Find(".crumbs a")
	for i := range sel.Nodes {
		single := sel.Eq(i)

		// checking that href exist on breadcrumb item
		attr, exist := single.Attr("href")
		if !exist {
			return
		}

		// we need only categories, categories contain /c/
		if strings.Contains(attr, "/c/") {
			categoryID := parseCategoryID(attr)
			categoryName := single.Text()

			newCategory := category{
				id:       categoryID,
				parentID: "0",
				name:     categoryName,
			}

			// checking that category is not already presented
			_, exist := (*categories)[categoryID]
			if !exist {
				(*categories)[categoryID] = newCategory
			}
		}
	}
}

// parsing category id
func parseCategoryID(attr string) string {
	urlTrimmed := strings.TrimSpace(attr)
	if len(urlTrimmed) == 0 {
		return ""
	}
	catID := strings.Split(urlTrimmed, "/")
	if len(catID) == 0 {
		return ""
	}

	return string(catID[len(catID)-1])
}

func getDocumentType(productPath string, fromURL bool) *goquery.Document {
	if fromURL {
		resp := prepareClient(productPath, 120)

		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			log.Fatal(err)
		}

		return doc
	}

	bytes, err := os.Open(productPath)
	if err != nil {
		log.Fatal("File not found.")
	}

	doc, err := goquery.NewDocumentFromReader(bytes)
	if err != nil {
		log.Fatal("File not read.")
	}

	return doc
}

// prepare client with custom settings
// for this parse to work we need to set a custom timeout
// site is just too damn slow
func prepareClient(productPath string, setSeconds time.Duration) *http.Response {
	timeout := time.Duration(setSeconds * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(productPath)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

// parse product price out of string
func parsePrice(wrap string) string {
	priceTrimmed := strings.TrimSpace(wrap)
	re := regexp.MustCompile("([0-9]+)|([,.][0-9]+)")
	numbers := re.FindAllString(priceTrimmed, -1)
	price := strings.Join(numbers, "")
	if len(price) == 0 {
		return "0"
	}

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
