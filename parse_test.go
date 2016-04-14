package main

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestProductLocations(t *testing.T) {
	locations := getLocations(sitemapLocation)
	productLocations := getProductLocations(locations)

	if len(productLocations) == 0 {
		t.Fatal("Locations length was 0")
	}

	for _, loc := range productLocations {
		if !strings.Contains(loc, "http://") {
			t.Fatal("Some of the links didn't contained http")
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	randLoc := productLocations[rand.Intn(len(productLocations))]

	response, err := http.Get(randLoc)
	if err != nil {
		t.Fatal("Can't get a random location via http")
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
