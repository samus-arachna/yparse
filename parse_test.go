package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestProductFindTitle(t *testing.T) {
	bytes, err := os.Open("test/product1.html")
	if err != nil {
		t.Fatal("cant open product1.html")
	}

	doc, err := goquery.NewDocumentFromReader(bytes)
	if err != nil {
		t.Fail()
	}

	title := doc.Find(".product_overview > h1").Text()
	if len(title) == 0 {
		t.Fatal("Title not found.")
	}
}

func TestOpenSitemap(t *testing.T) {
	bytes, err := ioutil.ReadFile("test/sitemap.xml")
	if err != nil {
		t.Fail()
	}

	sitemapFile := string(bytes)

	if len(sitemapFile) == 0 {
		t.Fail()
	}
}

func TestLoadSiteMap(t *testing.T) {
	bytes, err := os.Open("test/sitemap.xml")
	if err != nil {
		t.Fail()
	}

	doc, err := goquery.NewDocumentFromReader(bytes)
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	locs := doc.Find("loc")

	if len(locs.Text()) == 0 {
		t.Fail()
	}
}
