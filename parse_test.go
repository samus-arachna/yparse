package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

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

	/*
		for i := range locs.Nodes {
			single := locs.Eq(i)
			fmt.Println(single.Text())
		}
	*/
}
