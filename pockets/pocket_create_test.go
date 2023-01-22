//go:build unit

package pocket

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreatePocket(t *testing.T) {
	tests := []struct {
		name       string
		sqlFn      func() (*sql.DB, error)
		reqBody    string
		wantStatus int
		wantBody   string
	}{
		{"create pocket succesfully",
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				row := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(cStmt).WithArgs(0.0, "test", 2, "THB").WillReturnRows(row)
				return db, err
			},
			`{"amount": 0.0, "name": "test", "accountId": 2, "currency": "THB"}`,
			http.StatusCreated,
			`{"id": 1, "amount": 0.0, "name": "test", "accountId": 2, "currency": "THB"}`,
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
			assert.NoError(t, h.Create(c))
			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.JSONEq(t, tc.wantBody, rec.Body.String())
		})
	}
}
