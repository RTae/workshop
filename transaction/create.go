package transaction

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	pocket "github.com/kkgo-software-engineering/workshop/pockets"
)

const (
	cStmt = `INSERT INTO TBL_Transactions ("fromPocketId", "toPocketId", amount) VALUES ($1, $2, $3) 
	RETURNING id, "fromPocketId", "toPocketId", amount, date`
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
	row := h.db.QueryRowContext(ctx, `INSERT INTO TBL_Transactions ("fromPocketId", "toPocketId", amount) VALUES ($1, $2, $3) 
	RETURNING id, "fromPocketId", "toPocketId", amount, date`, uint(from), tn.To, tn.Amount)
	err = row.Scan(&tn.ID, &tn.From, &tn.To, &tn.Amount, &txDate)

	tn.Date = txDate.Format(time.RFC3339)

	if err != nil {
		logger.Error("Cannot insert transactions", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Err{Message: fmt.Sprint("Cannot insert transactions ", err.Error())})
	}

	// Get From Pocket By Id
	fp, err := getPocketById(c, h.db, tn.From)
	if err != nil {
		return err
	}

	// Get To Pocket By Id
	tp, err := getPocketById(c, h.db, tn.To)
	if err != nil {
		return err
	}

	// Discount from Pocket Amount
	// fp.Amount = fp.Amount - tn.Amount
	_, err = updateAmountFromPocketById(c, h.db, fp, tn.Amount)
	if err != nil {
		return err
	}

	// // // Increase To Pocket Amount
	// tp.Amount = tp.Amount + tn.Amount
	_, err = updateAmountToPocketById(c, h.db, tp, tn.Amount)
	if err != nil {
		fmt.Println("444")
		return err
	}

	return c.JSON(http.StatusCreated, tn)
}

func updateAmountFromPocketById(c echo.Context, db *sql.DB, pocket *pocket.Pocket, amount float64) (*pocket.Pocket, error) {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	row := db.QueryRowContext(ctx, "UPDATE TBL_Pockets SET amount = (amount - $2) WHERE id = $1 RETURNING id, amount;", pocket.ID, amount)
	err := row.Scan(&pocket.ID, &pocket.Amount)
	if err != nil {
		logger.Error("cannot update pocket", zap.Error(err))
		return nil, c.JSON(http.StatusInternalServerError, Err{Message: "cannot update pocket"})
	}

	return pocket, nil
}

func updateAmountToPocketById(c echo.Context, db *sql.DB, pocket *pocket.Pocket, amount float64) (*pocket.Pocket, error) {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	row := db.QueryRowContext(ctx, "UPDATE TBL_Pockets SET amount = (amount + $2) WHERE id = $1 RETURNING id, amount;", pocket.ID, amount)
	err := row.Scan(&pocket.ID, &pocket.Amount)
	if err != nil {
		logger.Error("cannot update pocket", zap.Error(err))
		return nil, c.JSON(http.StatusInternalServerError, Err{Message: "cannot update pocket"})
	}

	return pocket, nil
}

func getPocketById(c echo.Context, db *sql.DB, id uint) (*pocket.Pocket, error) {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	var p pocket.Pocket
	err := db.QueryRowContext(ctx, `SELECT id, amount, "name", "accountId" FROM tbl_pockets WHERE id = $1;`, id).Scan(&p.ID, &p.Amount, &p.Name, &p.AccountId)
	if err != nil {
		match, errMatch := regexp.MatchString("invalid input syntax", err.Error())
		if match {
			logger.Error("Pocket id must be integer", zap.Error(err))
			return nil, c.JSON(
				http.StatusUnprocessableEntity,
				Err{Message: "Pocket id must be integer"},
			)
		}
		if errMatch != nil {
			logger.Error("Match fail", zap.Error(err))
			return nil, c.JSON(
				http.StatusInternalServerError,
				Err{Message: err.Error()},
			)
		}
		match, errMatch = regexp.MatchString("no rows in result set", err.Error())
		if match {
			logger.Error("Pocket not found", zap.Error(err))
			return nil, c.JSON(
				http.StatusNotFound,
				Err{Message: "Pocket not found"},
			)
		}
		if errMatch != nil {
			logger.Error("Match fail", zap.Error(err))
			return nil, c.JSON(
				http.StatusInternalServerError,
				Err{Message: err.Error()},
			)
		}
		logger.Error("Internal error", zap.Error(err))
		return nil, c.JSON(
			http.StatusInternalServerError,
			Err{Message: err.Error()},
		)
	}
	return &p, nil
}
