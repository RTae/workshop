package account

import (
	"database/sql"

	"github.com/kkgo-software-engineering/workshop/config"
)

type Account struct {
	ID      int64   `json:"id"`
	Balance float64 `json:"balance"`
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}
