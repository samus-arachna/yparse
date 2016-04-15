package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

func getXMLDocument(products []map[string]string) {

}

func getXMLProduct(product map[string]string) {
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
		ID:          product["id"],
		Available:   "true",
		URL:         product["url"],
		Price:       product["price"],
		OldPrice:    product["oldprice"],
		CurrencyID:  "RUB",
		CategoryID:  "1337",
		Picture:     product["image"],
		Name:        product["title"],
		Vendor:      "ИВ РОШЕ",
		Description: product["desc"],
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent(" ", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func exampleEncoder() {
	type Address struct {
		City, State string
	}
	type Person struct {
		XMLName   xml.Name `xml:"person"`
		ID        int      `xml:"id,attr"`
		FirstName string   `xml:"name>first"`
		LastName  string   `xml:"name>last"`
		Age       int      `xml:"age"`
		Height    float32  `xml:"height,omitempty"`
		Married   bool
		Address
		Comment string `xml:",comment"`
	}

	v := &Person{ID: 13, FirstName: "John", LastName: "Doe", Age: 42}
	v.Comment = " Need more details. "
	v.Address = Address{"Hanga Roa", "Easter Island"}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("  ", "    ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
