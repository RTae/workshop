package router

import (
	"database/sql"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/account"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/featflag"
	"github.com/kkgo-software-engineering/workshop/healthchk"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	pocket "github.com/kkgo-software-engineering/workshop/pockets"
	"github.com/kkgo-software-engineering/workshop/transaction"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func RegRoute(cfg config.Config, logger *zap.Logger, db *sql.DB) *echo.Echo {
	e := echo.New()
	e.Use(mlog.Middleware(logger))
	e.Use(middleware.BasicAuth(mw.Authenicate()))

	hHealthChk := healthchk.New(db)
	e.GET("/healthz", hHealthChk.Check)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	hAccount := account.New(cfg.FeatureFlag, db)
	e.POST("/accounts", hAccount.Create)

	hPocket := pocket.New(db)
	e.GET("/pockets", hPocket.Get)
	e.POST("/pockets", hPocket.Create)
	e.GET("/pockets/:id", hPocket.GetByID)

	hFeatFlag := featflag.New(cfg)
	e.GET("/features", hFeatFlag.List)

	hTransaction := transaction.New(cfg.FeatureFlag, db)
	e.GET("/pockets/:id/transactions", hTransaction.GetAll)
	e.POST("/pockets/:id/transfer", hTransaction.Create)

	return e
}
