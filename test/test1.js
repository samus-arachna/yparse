'use strict';

let request = require('request');
let cheerio = require('cheerio');
let async = require('async');

const siteMapRootUrl = 'http://www.yves-rocher.ru/sitemap.xml';

request(siteMapRootUrl, (error, response, html) => {
  if (error) throw new Error('Could not get siteMapRootUrl');

  let $ = cheerio.load(html);

  const locations = getSitemapsLocations($);

  // for each sitemap.. make request
  let products = [];
  let locationsLength = locations.length;
  let locationsComplete = 0;

  locations.forEach((locationUrl) => {

    request(locationUrl, (error, response, html) => {
      if (error) throw new Error('Could not get locationUrl');

      let $ = cheerio.load(html);

      $('loc').each((i, e) => {
        products.push($(e).text());
      });

      locationsComplete++;
      if (locationsComplete == locationsLength) {
        console.log(products.length);
        console.log('Everything is ready!');
      }
    });

  });

  // get products from all
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
