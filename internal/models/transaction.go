package models

import "time"

type Transaction struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	Name     string    `json:"name"`
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Type     string    `json:"type"`
}