package pocket

import (
	"database/sql"
)

type Pocket struct {
	ID        uint    `json:"id"`
	Amount    float64 `json:"amount"`
	Name      string  `json:"name"`
	AccountId uint    `json:"accountId"`
	Currency  string  `json:"currency"`
}

type handler struct {
	db *sql.DB
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func New(db *sql.DB) *handler {
	return &handler{db}
}
