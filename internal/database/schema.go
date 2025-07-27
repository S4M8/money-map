package database

import "log"

func CreateTables() {
	createTransactionsTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL DEFAULT 1,
		date TIMESTAMPTZ NOT NULL,
		name TEXT,
		amount NUMERIC(10, 2) NOT NULL,
		category TEXT NOT NULL,
		type TEXT NOT NULL
	);
	`
	_, err := DB.Exec(createTransactionsTable)
	if err != nil {
		log.Fatal("failed to create transactions table:", err)
	}

	createFundsTable := `
	CREATE TABLE IF NOT EXISTS funds (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL DEFAULT 1,
		emergency_fund NUMERIC(10, 2) NOT NULL,
		education_fund NUMERIC(10, 2) NOT NULL,
		investments NUMERIC(10, 2) NOT NULL,
		other NUMERIC(10, 2) NOT NULL
	);
	`
	_, err = DB.Exec(createFundsTable)
	if err != nil {
		log.Fatal("failed to create funds table:", err)
	}

	// Insert a default row into funds if it's empty
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM funds").Scan(&count)
	if err != nil {
		log.Fatal("failed to check funds table:", err)
	}

	if count == 0 {
		insertFunds := `
		INSERT INTO funds (user_id, emergency_fund, education_fund, investments, other)
		VALUES (1, 0, 0, 0, 0);
		`
		_, err = DB.Exec(insertFunds)
		if err != nil {
			log.Fatal("failed to insert default funds:", err)
		}
	}
}
