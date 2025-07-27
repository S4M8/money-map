package database

import (
	"fmt"
	"money-map/internal/models"
)

func GetIncome(month, year int) ([]models.Transaction, error) {
	query := `
		SELECT id, date, name, amount, category
		FROM transactions
		WHERE type = 'Income' AND EXTRACT(MONTH FROM date) = $1 AND EXTRACT(YEAR FROM date) = $2
	`
	rows, err := DB.Query(query, month, year)
	if err != nil {
		return nil, fmt.Errorf("could not get income: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Name, &t.Amount, &t.Category); err != nil {
			return nil, fmt.Errorf("could not scan income transaction: %w", err)
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func GetExpenses(month, year int) ([]models.Transaction, error) {
	query := `
		SELECT id, date, name, amount, category
		FROM transactions
		WHERE type = 'Expense' AND EXTRACT(MONTH FROM date) = $1 AND EXTRACT(YEAR FROM date) = $2
	`
	rows, err := DB.Query(query, month, year)
	if err != nil {
		return nil, fmt.Errorf("could not get expenses: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Name, &t.Amount, &t.Category); err != nil {
			return nil, fmt.Errorf("could not scan expense transaction: %w", err)
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func AddIncome(t models.Transaction) error {
	_, err := DB.Exec("INSERT INTO transactions (date, name, amount, category, type) VALUES ($1, $2, $3, 'Income', 'Income')", t.Date, t.Name, t.Amount)
	return err
}

func AddExpense(t models.Transaction) error {
	_, err := DB.Exec("INSERT INTO transactions (date, name, amount, category, type) VALUES ($1, $2, $3, $4, 'Expense')", t.Date, t.Name, t.Amount, t.Category)
	return err
}

func DeleteIncome(id int) error {
	_, err := DB.Exec("DELETE FROM transactions WHERE id = $1 AND type = 'Income'", id)
	return err
}

func DeleteExpense(id int) error {
	_, err := DB.Exec("DELETE FROM transactions WHERE id = $1 AND type = 'Expense'", id)
	return err
}