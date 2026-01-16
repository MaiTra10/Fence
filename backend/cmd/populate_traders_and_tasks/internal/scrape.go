package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)

func Scrape(url string, cssSelector string, htmlElement *string) {

	// Start timer to keep track of time taken to scrape
	start := time.Now()

	// Start Colly collector instance to be used for scraping the EFT Wiki
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Colly)"),
	)

	// Capture the elemet defined by the CSS selector
	c.OnHTML(cssSelector, func(e *colly.HTMLElement) {

		html, err := e.DOM.Html()
		if err != nil {
			log.Fatal("Error getting 'div.wds-tabber.dealer-tabber' HTML:", err)
		}

		*htmlElement = html

	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Fatal("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Visit the site
	err := c.Visit(url)
	if err != nil {
		log.Fatal("Failed to visit", url, "\nError:", err)
	}

	// Log elapsed time to complete scrape
	elapsed := time.Since(start)
	fmt.Println("Scrape completed in:", elapsed)
	fmt.Println("Length of stored HTML:", len(*htmlElement))

}
