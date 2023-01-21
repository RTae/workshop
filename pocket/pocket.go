package pocket

import (
	"database/sql"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Pocket struct {
	ID        uint    `json:"id"`
	Ammount   float64 `json:"ammount"`
	Name      string  `json:"name"`
	AccountId uint    `json:"accountId"`
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

const gStmt = "SELECT id, ammount, name, account_id FROM pockets;"

func (h handler) Get(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	var pockets []Pocket
	rows, err := h.db.QueryContext(ctx, gStmt)
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	defer rows.Close()
	for rows.Next() {
		var pocket Pocket
		err = rows.Scan(&pocket.ID, &pocket.Ammount, &pocket.Name, &pocket.AccountId)
		if err != nil {
			logger.Error("scan error", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
		pockets = append(pockets, pocket)
	}
	return c.JSON(http.StatusOK, pockets)
}
