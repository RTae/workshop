//go:build unit

package pocket

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetPockets(t *testing.T) {
	tests := []struct {
		name       string
		cfgFlag    config.FeatureFlag
		sqlFn      func() (*sql.DB, error)
		reqBody    string
		wantStatus int
		wantBody   string
	}{
		{"Get pockets successfully",
			config.FeatureFlag{},
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				row := sqlmock.NewRows([]string{"ID", "Amount", "Name", "AccountId"}).AddRow(
					"1",
					79.00,
					"Test pocket",
					"1",
				)
				mock.ExpectQuery(gStmt).WillReturnRows(row)
				return db, err
			},
			``,
			http.StatusOK,
			`[{"id": 1, "amount": 79.00, "name": "Test pocket", "accountId": 1}]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			db, err := tc.sqlFn()
			h := New(db)
			// Assertions
			assert.NoError(t, err)
			assert.NoError(t, h.Get(c))
			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.JSONEq(t, tc.wantBody, rec.Body.String())
		})
	}
}

func TestGetPockets_Error(t *testing.T) {
	someErr := errors.New("some random error")
	tests := []struct {
		name    string
		cfgFlag config.FeatureFlag
		sqlFn   func() (*sql.DB, error)
		reqBody string
		wantErr string
	}{
		{"Get pocket failed",
			config.FeatureFlag{},
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				mock.ExpectQuery(gStmt).WillReturnError(someErr)
				return db, err
			},
			`{"balance": 1000.0}`,
			`{"message":"some random error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			db, _ := tc.sqlFn()
			h := New(db)

			h.Get(c)
			// Assertions
			assert.JSONEq(t, rec.Body.String(), tc.wantErr)
		})
	}
}
