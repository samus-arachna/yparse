package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func getXMLDocument(products []map[string]string, categories map[string]category) {
	tpl := getXMLTemplate()

	// writing products to xml file
	productsXML := ""
	for _, product := range products {
		productsXML += "\n" + getXMLProduct(product)
	}
	tpl = strings.Replace(tpl, "%OFFERS%", productsXML, 1)
	// END writing products to xml file

	fmt.Println(tpl)
}

// TODO
func getXMLCategory(categories map[string]category) string {
	type Category struct {
		XMLName  xml.Name `xml:"category"`
		ID       string   `xml:"id,attr"`
		ParentID string   `xml:"parentId,attr"`
		Value    string
	}

	v := &Category{
		ID:       "1337",
		ParentID: "7331",
		Value:    "yo",
	}

	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	enc.Indent(" ", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return b.String()
}

func getXMLTemplate() string {
	path := "data/template.xml"
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(dat)
}

func getXMLProduct(product map[string]string) string {
	type Product struct {
		XMLName     xml.Name `xml:"offer"`
		ID          string   `xml:"id,attr"`
		Available   string   `xml:"available,attr"`
		URL         string   `xml:"url"`
		Price       string   `xml:"price"`
		OldPrice    string   `xml:"oldprice"`
		CurrencyID  string   `xml:"currencyId"`
		CategoryID  string   `xml:"categoryId"`
		Picture     string   `xml:"picture"`
		Name        string   `xml:"name"`
		Vendor      string   `xml:"vendor"`
		Description string   `xml:"description"`
	}

	v := &Product{
		ID:          product["code"],
		Available:   "true",
		URL:         product["url"],
		Price:       product["price"],
		OldPrice:    product["oldprice"],
		CurrencyID:  "RUB",
		CategoryID:  product["categoryID"],
		Picture:     product["image"],
		Name:        product["title"],
		Vendor:      "ИВ РОШЕ",
		Description: product["desc"],
	}

	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	enc.Indent(" ", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return b.String()
}
