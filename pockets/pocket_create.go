package pocket

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	cStmt = `INSERT INTO tbl_pockets (amount, name, "accountId", currency) VALUES ($1, $2, $3, $4) RETURNING id;`
)

func (h handler) Create(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	var p Pocket
	err := c.Bind(&p)
	if err != nil {
		logger.Error("Bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	if p.Currency == "" {
		p.Currency = "THB"
	}

	if p.Amount > 0.0 {
		p.Amount = 0.0
	}

	if len(p.Currency) != 3 {
		logger.Error("Bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Currency must be 3 characters"})
	}

	var lastInsertId uint
	err = h.db.QueryRowContext(ctx, cStmt, p.Amount, p.Name, p.AccountId, p.Currency).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	logger.Info("create successfully", zap.Uint("id", lastInsertId))
	p.ID = lastInsertId
	return c.JSON(http.StatusCreated, p)
}
