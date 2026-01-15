package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Trader struct {
	Name     string
	ImageURL string
	Quests   []Quest
}

type Quest struct {
	Name          string
	WikiURL       string
	Objectives    []string
	Rewards       []string
	RequiredKappa bool
}

func main() {
	start := time.Now()

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Colly)"),
	)

	// Capture the target div
	c.OnHTML("div.wds-tabber.dealer-tabber", func(e *colly.HTMLElement) {
		html, err := e.DOM.Html()
		if err != nil {
			log.Fatal("Failed to extract HTML:", err)
		}

		// Write to txt file in same directory
		err = os.WriteFile("dealer-tabber.txt", []byte(html), 0644)
		if err != nil {
			log.Fatal("Failed to write file:", err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatalf("Request failed (%d): %v", r.StatusCode, err)
	})

	// Start scraping
	err := c.Visit("https://escapefromtarkov.fandom.com/wiki/Quests")
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	fmt.Println("Scrape completed in:", elapsed)

	// Query Quests

	file, err := os.Open("dealer-tabber.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	var traders []Trader

	doc.Find("ul.wds-tabs.wds-tabs__wrapper > li").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Find("span[title]").Attr("title")

		imgSel := s.Find("img")
		img, exists := imgSel.Attr("data-src")
		if !exists {
			img, _ = imgSel.Attr("src")
		}

		traders = append(traders, Trader{
			Name:     strings.TrimSpace(name),
			ImageURL: strings.TrimSpace(img),
		})
	})

	// 2️⃣ Extract quest containers (same order!)
	doc.Find("div.wds-tab__content").Each(func(i int, content *goquery.Selection) {
		if i >= len(traders) {
			return
		}

		var quests []Quest

		content.Find("table tr").Each(func(_ int, row *goquery.Selection) {
			link := row.Find("td").Eq(1).Find("a")
			if link.Length() == 0 {
				return
			}

			name := strings.TrimSpace(link.Text())
			href, _ := link.Attr("href")

			var objectives []string
			row.Find("td").Eq(2).Find("li").Each(func(_ int, li *goquery.Selection) {
				objectives = append(objectives, strings.TrimSpace(li.Text()))
			})

			var rewards []string
			row.Find("td").Eq(3).Find("li").Each(func(_ int, li *goquery.Selection) {
				rewards = append(rewards, strings.TrimSpace(li.Text()))
			})

			kappaText := strings.ToLower(row.Text())
			requiredKappa := strings.Contains(kappaText, "kappa")

			quests = append(quests, Quest{
				Name:          name,
				WikiURL:       "https://escapefromtarkov.fandom.com" + href,
				Objectives:    objectives,
				Rewards:       rewards,
				RequiredKappa: requiredKappa,
			})
		})

		traders[i].Quests = quests
	})

	// Debug output

	fmt.Println(traders[2].Quests[9])

}
