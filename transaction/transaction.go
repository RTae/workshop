package transaction

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
)

type Transaction struct {
	ID     int64     `json:"id"`
	From   int64     `json:"from"`
	To     int64     `json:"to"`
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
}

const (
	cStmt = "SELECT * FROM TBL_transactions WHERE pocketId = $1"
)

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

func (h handler) GetAll(c echo.Context) error {

	return c.JSON(http.StatusOK, []Transaction{
		{
			ID:     1,
			From:   12345,
			To:     67890,
			Amount: 50.0,
			Date:   time.Now(),
		},
	})
}
