package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestParsePrice(t *testing.T) {
	first := parsePrice("   1770 рублей")
	if first != "1770" {
		t.Fatal("Price is not valid")
	}

	second := parsePrice("330,00 р")
	if second != "330,00" {
		t.Fatal("Price is not valid")
	}

	third := parsePrice("  1,770.00  ")
	if third != "1,770.00" {
		t.Fatal("Price is not valid")
	}

	fourth := parsePrice("  1000 dollars  ")
	if fourth != "1000" {
		t.Fatal("Price is not valid")
	}
}

func TestParseCode(t *testing.T) {
	first := parseCode("Код&nbsp;82283&nbsp;- Тюбик&nbsp;20&nbsp;мл")
	if first != "82283" {
		t.Fatal("Code is not valid")
	}

	second := parseCode("Код 04487 - Чтото")
	if second != "04487" {
		t.Fatal("Code is not valid")
	}

	third := parseCode("Code - What")
	if third != "0" {
		t.Fatal("Code is not valid")
	}
}

func TestProductLocations(t *testing.T) {
	locations := getLocations(sitemapLocation)
	productLocations := getProductLocations(locations)

	if len(productLocations) == 0 {
		t.Fatal("Locations length was 0")
	}

	for _, loc := range productLocations {
		if !strings.Contains(loc, "/p/") {
			t.Fatal("Some of the links didn't contained /p/")
		}

		if !strings.Contains(loc, "http://") {
			t.Fatal("Some of the links didn't contained http")
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	randLoc := productLocations[rand.Intn(len(productLocations))]

	response, err := http.Get(randLoc)
	if err != nil {
		t.Fatal("Can't get a random location via http " + err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal("Can't get a location body")
	}

	if len(string(body)) == 0 {
		t.Fatal("Location body is empty")
	}

	defer response.Body.Close()
}

func TestGetLocations(t *testing.T) {
	locations := getLocations(sitemapLocation)

	if len(locations) == 0 {
		t.Fatal("Locations length was 0")
	}

	for _, loc := range locations {
		if !strings.Contains(loc, "http://") {
			t.Fatal("Some of the links didn't contained http")
		}

		if !strings.Contains(loc, "yves-rocher.ru") {
			t.Fatal("Some of the links didn't contained yves-rocher.ru")
		}
	}
}
