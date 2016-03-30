'use strict';

let request = require('request');
let cheerio = require('cheerio');

const siteMapRootUrl = 'http://www.yves-rocher.ru/sitemap.xml';

request(siteMapRootUrl, (error, response, html) => {
  if (error) throw new Error('Could not get siteMapRootUrl');

  let $ = cheerio.load(html);

  const locations = getSitemapsLocations($);

  let productsUrl = [];
  let locationsComplete = 0;
  const locationsLength = locations.length;

  locations.forEach((locationUrl) => {
    request(locationUrl, (error, response, html) => {
      if (error) throw new Error('Could not get locationUrl');

      let $ = cheerio.load(html);

      // filtering only products
      $('loc').each((i, e) => {
        if ($(e).text().includes('/p/')) {
          productsUrl.push($(e).text());
        }
      });

      locationsComplete++;

      if (locationsComplete == locationsLength) {
        console.log('Loaded all sitemaps, starting to parse products');
        console.log('Number of products: ' + productsUrl.length);

        productsUrl.forEach((productUrl) => {
          //console.log(productUrl);
        });
      }
    });
  });
});

// FUNCTIONS
function getSitemapsLocations($) {
  let locations = [];

  $('loc').each((i, e) => {
    locations.push($(e).text());
  });

  return locations;
}

function getProducts(locations) {

}
