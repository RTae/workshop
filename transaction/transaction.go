package transaction

import (
	"database/sql"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	gStmt = `SELECT * FROM tbl_transactions WHERE "fromPocketId" = $1 OR "toPocketId" = $2`
)

type Transaction struct {
	ID     uint    `json:"id"`
	From   uint    `json:"from"`
	To     uint    `json:"to"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
}

type Err struct {
	Message string `json:"message"`
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

func (h handler) GetAll(c echo.Context) error {
	logger := mlog.L(c)

	pocketId := c.Param("id")
	ctx := c.Request().Context()

	rows, err := h.db.QueryContext(ctx, gStmt, pocketId, pocketId)
	if err != nil {
		logger.Error("query transactions error", zap.Error(err))
		return err
	}

	var txns []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.From, &t.To, &t.Amount, &t.Date)
		if err != nil {

			return c.JSON(http.StatusInternalServerError, zap.Error(err))
		}
		txns = append(txns, t)
	}

	return c.JSON(http.StatusOK, txns)
}
