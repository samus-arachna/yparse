package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRunParse(t *testing.T) {
	locations := getLocations(sitemapLocation)
	productLocations := getProductLocations(locations)
	if len(productLocations) == 0 {
		t.Fatal("Locations length was 0")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	randStart := rand.Intn(len(productLocations))
	randEnd := 0

	if (randStart + 10) > len(productLocations) {
		randStart = len(productLocations) - 15
		randEnd = len(productLocations) - 5
	} else {
		randEnd = randStart + 10
	}

	randLoc := productLocations[randStart:randEnd]

	_, count := runParse(randLoc, 3)

	if len(randLoc) != count {
		t.Fatal("You parsed less products, then you initially had.")
	}
}

func TestParseProduct2(t *testing.T) {
	pr := "data/product2.html"

	_, err := parseProduct(pr, false)
	if err != nil {
		return
	}

	t.Fatal("This product have no current price, so it can't be parsed.")
}

func TestParseProduct1(t *testing.T) {
	pr := "data/product1.html"

	product, err := parseProduct(pr, false)
	if err != nil {
		t.Fatal("Cant parse product")
	}

	if product["code"] != "82283" {
		t.Fatal("Code is not valid")
	}

	if product["title"] != "Крем-Заполнитель Глубоких Морщин" {
		t.Fatal("Title is not valid")
	}

	if product["desc"] !=
		"Мгновенно сокращает даже самые глубокие морщины и заполняет их изнутри" {
		t.Fatal("Desc is not valid")
	}

	if product["img"] !=
		"http://w23.yves-rocher-statics.com/medias/sys_master/retina/images/h54/h34/8864784318494.jpg" {
		t.Fatal("Image path is not valid")
	}

	if product["price"] != "1,049" {
		t.Fatal("Price is not valid")
	}

	if product["priceOld"] != "1,390" {
		t.Fatal("Old price is not valid")
	}
}

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

	fifth := parsePrice("")
	if fifth != "0" {
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
