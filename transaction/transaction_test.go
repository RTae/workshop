package transaction

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetTransactions(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		cfgFlag    config.FeatureFlag
		sqlFn      func() (*sql.DB, error)
		reqBody    string
		wantStatus int
		wantBody   []Transaction
	}{
		{
			"get all transaction for pocket id: 1",
			config.FeatureFlag{},
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}

				row := sqlmock.NewRows([]string{"id", "from", "to", "amount", "date"}).AddRow(1, 12345, 67890, 50.0, now)
				mock.ExpectQuery(gStmt).WithArgs("12345", "12345").WillReturnRows(row)
				return db, err
			},
			``,
			http.StatusOK,
			[]Transaction{
				{
					ID:     1,
					From:   12345,
					To:     67890,
					Amount: 50.0,
					Date:   now.Format(time.RFC3339),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetPath("/pockets/:id/transactions")
			c.SetParamNames("id")
			c.SetParamValues("12345")

			db, err := tc.sqlFn()
			h := New(tc.cfgFlag, db)
			// Assertions
			assert.NoError(t, err)
			if assert.NoError(t, h.GetAll(c)) {
				assert.Equal(t, tc.wantStatus, rec.Code)
				// assert.JSONEq(t, tc.wantBody, rec.Body.String())
			}
		})
	}

}
