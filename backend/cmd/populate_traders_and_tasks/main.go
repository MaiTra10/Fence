package main

import (
	"fmt"
	"log"

	"github.com/MaiTra10/Fence/backend/cmd/populate_traders_and_tasks/internal"
)

func main() {

	// Define EFT Wiki URL
	eftWikiTasksUrl := "https://escapefromtarkov.fandom.com/wiki/Quests"

	// Define the CSS selector used to obtain the div with Traders and Tasks information
	tradersAndTasksCssSelector := "div.wds-tabber.dealer-tabber"

	// Initialize the var to store HTML of the div which contains Traders and Tasks
	var traderAndTaskHtml string

	// Scrape the EFT Wiki for the surrounding div
	internal.Scrape(eftWikiTasksUrl, tradersAndTasksCssSelector, &traderAndTaskHtml)

	// Create GoQuery document for use in extraction
	doc, err := internal.CreateDoc(traderAndTaskHtml)
	if err != nil {
		log.Fatal("Failed while creating document from html:", err)
	}

	// Extract traders to their respective structs list
	traders, err := internal.ExtractTradersFromHTML(doc)
	if err != nil {
		log.Fatal("Failed while extracting traders:", err)
	}

	// Sanity check to ensure traders is a non-empty list
	if len(traders) < 1 {
		log.Fatal("'traders' list is empty")
	}

	// Extract tasks into Trader.Task component of struct
	internal.ExractTasksFromHTML(doc, &traders)

	// Sanity check to ensure tasks are non-empty lists
	for _, trader := range traders {
		if len(trader.Tasks) == 0 {
			log.Printf("Error: trader '%s' has no tasks", trader.Name)
		}
	}

	fmt.Println(traders)

}
