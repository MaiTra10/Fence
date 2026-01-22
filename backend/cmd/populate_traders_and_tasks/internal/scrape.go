package internal

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func GetTradersAndTasks() []Trader {

	// Define EFT Wiki URL
	eftWikiTasksUrl := "https://escapefromtarkov.fandom.com/wiki/Quests"

	// Define the CSS selector used to obtain the div with Traders and Tasks information
	tradersAndTasksCssSelector := "div.wds-tabber.dealer-tabber"

	// Initialize the var to store HTML of the div which contains Traders and Tasks
	var traderAndTaskHtml string

	// Scrape the EFT Wiki for the surrounding div
	scrape(eftWikiTasksUrl, tradersAndTasksCssSelector, &traderAndTaskHtml)

	// Create GoQuery document for use in extraction
	doc, err := CreateDoc(traderAndTaskHtml)
	if err != nil {
		log.Fatal("Failed while creating document from html:", err)
	}

	// Extract traders to their respective structs list
	traders, err := ExtractTradersFromHTML(doc)
	if err != nil {
		log.Fatal("Failed while extracting traders:", err)
	}

	// Sanity check to ensure traders is a non-empty list
	if len(traders) < 1 {
		log.Fatal("'traders' list is empty")
	}

	// Extract tasks into Trader.Task component of struct
	ExractTasksFromHTML(doc, &traders)

	return traders

}

func FillTaskRelatedQuests(traders []Trader) error {

	// Iterate though tasks within traders
	for traderIndex := range traders {
		for taskIndex := range traders[traderIndex].Tasks {

			task := &traders[traderIndex].Tasks[taskIndex]

			prereq, otherChoices, err := scrapeTaskRelatedQuests(task.WikiURL)
			if err != nil {
				return fmt.Errorf("failed to get prereq or other choices strings: %w", err)
			}

			fmt.Println(prereq)
			fmt.Println(otherChoices)

		}
	}

	return nil

}

func scrape(url string, cssSelector string, htmlElement *string) {

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

func scrapeTaskRelatedQuests(url string) (string, string, error) {

	start := time.Now()
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (compatible; Colly)"))

	var previous string
	var otherChoices string

	c.OnHTML("table.va-infobox-group", func(e *colly.HTMLElement) {

		// Check if this table is "Related quests"
		header := e.DOM.Find("th.va-infobox-header").First().Text()
		if strings.TrimSpace(header) != "Related quests" {
			return // skip tables that aren't related quests
		}

		// Now find <td> elements inside this table
		e.DOM.Find("td").Each(func(_ int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())

			if strings.HasPrefix(text, "Previous:") {
				var tasks []string
				s.Find("a").Each(func(_ int, a *goquery.Selection) {
					tasks = append(tasks, strings.TrimSpace(a.Text()))
				})
				previous = strings.Join(tasks, ", ")
			} else if strings.HasPrefix(text, "Other choices:") {
				var choices []string
				s.Find("a").Each(func(_ int, a *goquery.Selection) {
					choices = append(choices, strings.TrimSpace(a.Text()))
				})
				otherChoices = strings.Join(choices, ", ")
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Fatal("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(url)
	if err != nil {
		return "", "", fmt.Errorf("failed to visit task: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Println("Scrape completed in:", elapsed)
	return previous, otherChoices, nil
}
