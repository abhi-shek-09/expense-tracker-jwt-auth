package models

import "time"

type Expense struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	Category    string   `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}
