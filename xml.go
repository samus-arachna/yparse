package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func getXMLDocument(products []map[string]string, categories map[string]category) {
	tpl := getXMLTemplate()

	// setting date in xml file
	timeXML := time.Now().Format("02/01/2006 15:04")
	tpl = strings.Replace(tpl, "%DATE%", timeXML, 1)
	// END setting date in xml file

	// writing categories to xml file
	categoriesXML := ""
	for _, val := range categories {
		categoriesXML += "\n" + getXMLCategory(val)
	}
	tpl = strings.Replace(tpl, "%CATEGORIES%", categoriesXML, 1)
	// END writing categories to xml file

	// writing products to xml file
	productsXML := ""
	for _, product := range products {
		productsXML += "\n" + getXMLProduct(product)
	}
	tpl = strings.Replace(tpl, "%OFFERS%", productsXML, 1)
	// END writing products to xml file

	// writing down xml file to filesystem
	writeXML(tpl)
}

func writeXML(xml string) {
	ioutil.WriteFile("catalog.xml", []byte(xml), 0644)
}

func getXMLTemplate() string {
	path := "data/template.xml"
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(dat)
}

// getting xml categories list
func getXMLCategory(cat category) string {
	type Category struct {
		XMLName  xml.Name `xml:"category"`
		ID       string   `xml:"id,attr"`
		ParentID string   `xml:"parentId,attr,omitempty"`
		Value    string   `xml:",chardata"`
	}

	v := &Category{
		ID:       cat.id,
		ParentID: cat.parentID,
		Value:    cat.name,
	}

	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	enc.Indent(" ", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return b.String()
}

// getting xml products list
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
		Picture:     product["img"],
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
