package internal

import (
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
