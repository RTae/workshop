//go:build unit

package pocket

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetPockets(t *testing.T) {
	tests := []struct {
		name       string
		sqlFn      func() (*sql.DB, error)
		reqArg     string
		wantStatus int
		wantBody   string
	}{
		{"Get pockets successfully",
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				row := sqlmock.NewRows([]string{"ID", "Amount", "Name", "AccountId"}).
					AddRow(
						"1",
						79.00,
						"Test pocket",
						"1",
					)
				mock.ExpectQuery(gStmt).WillReturnRows(row)
				return db, err
			},
			"",
			http.StatusOK,
			`[{"id": 1, "amount": 79.00, "name": "Test pocket", "accountId": 1}]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
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
		sqlFn   func() (*sql.DB, error)
		reqArg  string
		wantErr string
	}{
		{"Should return internal service error if error happend",
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				mock.ExpectQuery(gStmt).
					WillReturnError(someErr)
				return db, err
			},
			"",
			`{"message":"some random error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
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

func TestGetPocketByID(t *testing.T) {
	tests := []struct {
		name       string
		sqlFn      func() (*sql.DB, error)
		reqArg     string
		wantStatus int
		wantBody   string
	}{
		{"Get pocket by id successfully",
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				row := sqlmock.NewRows([]string{"ID", "Amount", "Name", "AccountId"}).
					AddRow(
						"1",
						79.00,
						"Test pocket",
						"1",
					)
				mock.ExpectQuery(gbiStmt).WithArgs("1").WillReturnRows(row)
				return db, err
			},
			"1",
			http.StatusOK,
			`{"id": 1, "amount": 79.00, "name": "Test pocket", "accountId": 1}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.reqArg)

			db, err := tc.sqlFn()
			h := New(db)
			// Assertions
			assert.NoError(t, err)
			assert.NoError(t, h.GetByID(c))
			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.JSONEq(t, tc.wantBody, rec.Body.String())
		})
	}
}

func TestGetPocketByID_Error(t *testing.T) {
	someErr := errors.New("some random error")
	tests := []struct {
		name    string
		sqlFn   func() (*sql.DB, error)
		reqArg  string
		wantErr string
	}{
		{"Should return unprocessable entity error if pocket id is not integer",
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				mock.ExpectQuery(gbiStmt).
					WithArgs("dwdwdw").
					WillReturnError(errors.New("invalid input syntax"))
				return db, err
			},
			"dwdwdw",
			`{"message":"Param id must be integer"}`,
		},
		{"Should return not found error if the request pocket is not exist",
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				mock.ExpectQuery(gbiStmt).
					WithArgs("1").
					WillReturnError(errors.New("No rows in result set"))
				return db, err
			},
			"1",
			`{"message":"Record not found"}`,
		},
		{"Should return internal error if can not query pocket",
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				mock.ExpectQuery(gbiStmt).
					WithArgs("1").
					WillReturnError(someErr)
				return db, err
			},
			"1",
			`{"message":"some random error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.reqArg)

			db, _ := tc.sqlFn()
			h := New(db)

			h.GetByID(c)
			// Assertions
			assert.JSONEq(t, rec.Body.String(), tc.wantErr)
		})
	}
}
