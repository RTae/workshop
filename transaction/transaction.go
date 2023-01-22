package transaction

import (
	"database/sql"

	"github.com/kkgo-software-engineering/workshop/config"
)

type Tranaction struct {
	ID     uint    `json:"id"`
	From   uint    `json:"from"`
	To     uint    `json:"to"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
}

type Err struct {
	Message string `json:"message"`
}

// POST  /pockets/:id/transfers

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}
