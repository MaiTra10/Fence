package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CreateDoc(html string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func ExtractTradersFromHTML(doc *goquery.Document) ([]Trader, error) {

	var traders []Trader

	// Source: ChatGPT
	doc.Find("ul.wds-tabs.wds-tabs__wrapper > li").Each(func(i int, s *goquery.Selection) {
		// Get the trader name from the span[title] attribute
		name, _ := s.Find("span[title]").Attr("title")

		// Get the trader image
		imgSel := s.Find("img")
		img, exists := imgSel.Attr("data-src")
		if !exists {
			img, _ = imgSel.Attr("src")
		}

		// Append to slice
		traders = append(traders, Trader{
			Name:     strings.TrimSpace(name),
			ImageURL: strings.TrimSpace(img),
		})
	})

	return traders, nil

}

func ExractTasksFromHTML(doc *goquery.Document, traders *[]Trader) {

	doc.Find("div.wds-tab__content").Each(func(i int, content *goquery.Selection) {
		if i >= len(*traders) {
			return
		}

		var tasks []Task

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

			var requiredForKappa bool
			requiredForKappaText := row.Find("td").Eq(4).Find("font").Text()
			// Trim whitespace just in case
			requiredForKappaText = strings.TrimSpace(requiredForKappaText)
			// Set bool
			requiredForKappa = requiredForKappaText == "Yes"

			tasks = append(tasks, Task{
				Name:             name,
				WikiURL:          "https://escapefromtarkov.fandom.com" + href,
				Objectives:       objectives,
				Rewards:          rewards,
				RequiredForKappa: requiredForKappa,
			})
		})

		(*traders)[i].Tasks = tasks
	})

}

func ExtractPrereqTasks(doc *goquery.Document) ([]string, error) {
	return extractRelatedTasksByLabel(doc, "Previous:")
}

func ExtractOtherChoices(doc *goquery.Document) ([]string, error) {
	return extractRelatedTasksByLabel(doc, "Other choices:")
}

// Helper function for task related quests
// Source: ChatGPT
func extractRelatedTasksByLabel(doc *goquery.Document, label string) ([]string, error) {
	var tasks []string

	// Find the <td> that contains the label text
	selection := doc.Find("td.va-infobox-content").FilterFunction(
		func(_ int, s *goquery.Selection) bool {
			text := s.Text()

			// Replace non-breaking spaces with regular space
			text = strings.ReplaceAll(text, "\u00a0", " ")

			// Trim surrounding spaces and newlines
			text = strings.TrimSpace(text)

			// Make sure the label is at the start
			return strings.HasPrefix(text, label)
		},
	)

	fmt.Println("Selection:", selection)

	if selection.Length() == 0 {
		return nil, errors.New("label not found: " + label)
	}

	// Extract all anchor text values
	selection.Find("a").Each(func(_ int, s *goquery.Selection) {
		task := strings.TrimSpace(s.Text())
		if task != "" {
			tasks = append(tasks, task)
		}
	})

	// If no <a> tags exist → "-" case → empty slice
	return tasks, nil
}
