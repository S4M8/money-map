package database

import (
	"log"
	"time"

	"money-map/internal/models"
)

func GetDelta(month, year int) (models.Delta, error) {
	// Get the start and end of the selected month
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	// Get the total income for the selected month
	totalIncome, err := getMonthlyTotalForCategory("Income", "", startOfMonth, endOfMonth)
	if err != nil {
		return models.Delta{}, err
	}

	// Get the total expenses for each category for the selected month
	coreExpenses, err := getMonthlyTotalForCategory("Expense", "Core", startOfMonth, endOfMonth)
	if err != nil {
		return models.Delta{}, err
	}
	choiceExpenses, err := getMonthlyTotalForCategory("Expense", "Choice", startOfMonth, endOfMonth)
	if err != nil {
		return models.Delta{}, err
	}

	totalExpenses := coreExpenses + choiceExpenses
	remainingAmount := totalIncome - totalExpenses

	// Calculate the delta
	delta := models.Delta{
		TotalIncome:     totalIncome,
		TotalExpenses:   totalExpenses,
		RemainingAmount: remainingAmount,
		CoreExpenses:    coreExpenses,
		ChoiceExpenses:  choiceExpenses,
	}

	return delta, nil
}

func getMonthlyTotalForCategory(transactionType string, category string, startOfMonth time.Time, endOfMonth time.Time) (float64, error) {
	var total float64
	var query string
	var err error

	// Hardcode user_id to 1 for now
	userID := 1

	if category == "" {
		query = `
			SELECT COALESCE(SUM(amount), 0)
			FROM transactions
			WHERE type = $1 AND date >= $2 AND date < $3 AND user_id = $4;
		`
		err = DB.QueryRow(query, transactionType, startOfMonth, endOfMonth, userID).Scan(&total)
	} else {
		query = `
			SELECT COALESCE(SUM(amount), 0)
			FROM transactions
			WHERE type = $1 AND category = $2 AND date >= $3 AND date < $4 AND user_id = $5;
		`
		err = DB.QueryRow(query, transactionType, category, startOfMonth, endOfMonth, userID).Scan(&total)
	}

	if err != nil {
		log.Printf("Error getting monthly total for category %s: %v", category, err)
		return 0, err
	}
	return total, nil
}