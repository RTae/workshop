package transaction

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	cStmt = `INSERT INTO TBL_Transactions (fromPocketId, toPocketId, amount) VALUES ($1, $2, $3) 
	RETURNING id, fromPocketId, toPocketId, amount, date`
)

func (h handler) Create(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	from, err := strconv.Atoi(c.Param("id"))
	if err != nil || from <= 0 {
		logger.Error("invalid pocket id", zap.Error(err))
		return c.JSON(http.StatusBadRequest, Err{Message: "bad request body"})
	}

	var tn Transaction
	err = c.Bind(&tn)

	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, Err{Message: "bad request body"})
	}

	if tn.Amount <= 0 {
		logger.Error("amount must more than 0", zap.Error(err))
		return c.JSON(http.StatusBadRequest, Err{Message: "bad request body"})
	}

	if tn.To <= 0 {
		logger.Error("invalid pocket id", zap.Error(err))
		return c.JSON(http.StatusBadRequest, Err{Message: "bad request body"})
	}

	var txDate time.Time
	row := h.db.QueryRowContext(ctx, `INSERT INTO TBL_Transactions (fromPocketId, toPocketId, amount) VALUES ($1, $2, $3) 
		RETURNING id, fromPocketId, toPocketId, amount, date`, uint(from), tn.To, tn.Amount)
	err = row.Scan(&tn.ID, &tn.From, &tn.To, &tn.Amount, &txDate)

	tn.Date = txDate.Format(time.RFC3339)

	if err != nil {
		logger.Error("Cannot insert transactions", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Err{Message: fmt.Sprint("Cannot insert transactions ", err.Error())})
	}

	return c.JSON(http.StatusCreated, tn)
}
