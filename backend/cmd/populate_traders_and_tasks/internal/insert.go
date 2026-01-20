package internal

import (
	"context"
	"fmt"
	"strings"

	generic "github.com/MaiTra10/Fence/backend/generic/db"
)

func Populate(traders []Trader) error {

	// Connect to PostgreSQL
	ctx := context.Background()
	conn, err := generic.PSQLConnect()
	if err != nil {
		return fmt.Errorf("failed to connect to PSQL: %w", err)
	}
	defer conn.Close(ctx)

	// Begin transaction
	transaction, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer transaction.Rollback(ctx) // If any of the queries fail, undo all changes

	// Define trader insertion query
	insertTraderQuery := `
		INSERT INTO public.traders (name, image_url)
		VALUES ($1, $2)
		ON CONFLICT (name)
		DO UPDATE SET image_url = EXCLUDED.image_url
		RETURNING id;
	`
	// Define task insertion query
	insertTaskQuery := `
		INSERT INTO public.tasks (
			trader_id,
			name,
			wiki_url,
			objectives,
			rewards,
			required_for_kappa
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (trader_id, name)
		DO UPDATE SET
			wiki_url = EXCLUDED.wiki_url,
			objectives = EXCLUDED.objectives,
			rewards = EXCLUDED.rewards,
			required_for_kappa = EXCLUDED.required_for_kappa;
	`

	// Iterate through each trader to insert the trader and their tasks
	for _, trader := range traders {
		// Sanity check for if any traders have no tasks
		if len(trader.Tasks) == 0 {
			return fmt.Errorf("trader '%s' has no tasks", trader.Name)
		}
		fmt.Printf("%s: %d Tasks - ", trader.Name, len(trader.Tasks))

		// Add trader into transaction
		var traderID int
		err := transaction.QueryRow(
			ctx,
			insertTraderQuery,
			trader.Name,
			trader.ImageURL,
		).Scan(&traderID)
		if err != nil {
			return fmt.Errorf("failed inserting trader %s: %w", trader.Name, err)
		}
		fmt.Printf("Inserted trader (id=%d) - ", traderID)

		// Iterate through tasks for the trader
		for _, task := range trader.Tasks {
			// Join objectives and rewards since these are stored in Task struct as slice of string
			objectives := strings.Join(task.Objectives, ", ")
			rewards := strings.Join(task.Rewards, ", ")
			// Add trader's tasks
			_, err := transaction.Exec(
				ctx,
				insertTaskQuery,
				traderID,
				task.Name,
				task.WikiURL,
				objectives,
				rewards,
				task.RequiredForKappa,
			)
			if err != nil {
				return fmt.Errorf(
					"failed inserting task '%s' for trader '%s': %w",
					task.Name,
					trader.Name,
					err,
				)
			}
		}
		fmt.Print("All tasks added to table\n")
	}

	// Commit only after all inserts succeed
	if err := transaction.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
