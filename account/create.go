package account

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	cStmt         = "INSERT INTO tbl_accounts (balance) VALUES ($1) RETURNING id;"
	cBalanceLimit = 10000
)

var (
	hErrBalanceLimitExceed = echo.NewHTTPError(http.StatusBadRequest,
		"create account balance exceed limitation")
)

func (h handler) Create(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	var ac Account
	err := c.Bind(&ac)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	if h.cfg.IsLimitMaxBalanceOnCreate && ac.Balance > cBalanceLimit {
		logger.Error("account limit on account creating", zap.Error(hErrBalanceLimitExceed))
		return hErrBalanceLimitExceed
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, ac.Balance).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	ac.ID = lastInsertId
	return c.JSON(http.StatusCreated, ac)
}
