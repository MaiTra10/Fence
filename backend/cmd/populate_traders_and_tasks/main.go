package main

import (
	"log"

	"github.com/MaiTra10/Fence/backend/cmd/populate_traders_and_tasks/internal"
)

func main() {

	traders := internal.GetTradersAndTasks()

	if err := internal.FillTaskRelatedQuests(traders); err != nil {
		log.Fatal("Something went wrong while filling in task related quests:", err)
	}

	if err := internal.Populate(traders); err != nil {
		log.Fatal("Something went wrong while populating PSQL:", err)
	}

}
