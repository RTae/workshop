package transaction

import (
	"fmt"
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

func TestCreateTransaction(t *testing.T) {
	now := time.Now()
	t.Run("create transaction succesfully", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/pockets/1/transfer", strings.NewReader(`{"to": 2, "amount": 100.0}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/pockets/:id/transfer")
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		row := sqlmock.NewRows([]string{"id", "fromPocketId", "toPocketId", "amount", "date"}).
			AddRow(1, 1, 2, 100.0, now)
		mock.ExpectQuery(cStmt).WithArgs(1, 2, 100.0).WillReturnRows(row)

		// Act
		h := New(config.FeatureFlag{}, db)

		wantBody := fmt.Sprintf(`{"id": 1, "from": 1, "to": 2, "amount": 100.0, "date": "%v"}`, now.Format(time.RFC3339))

		// Assertions
		assert.NoError(t, err)
		if assert.NoError(t, h.Create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.JSONEq(t, wantBody, rec.Body.String())
		}
	})
}

func TestCreateTransaction_Error(t *testing.T) {
	now := time.Now()
	t.Run("Test Fail invalid from", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/pockets//transfers", strings.NewReader(`{"to": 2, "amount": 100.0}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		row := sqlmock.NewRows([]string{"id", "fromPocketId", "toPocketId", "amount", "date"}).
			AddRow(1, 1, 2, 100.0, now)
		mock.ExpectQuery(cStmt).WithArgs(1, 2, 100.0).WillReturnRows(row)

		h := New(config.FeatureFlag{}, db)
		c := e.NewContext(req, rec)
		c.SetPath("/pockets/:id/transfer")
		c.SetParamNames("id")
		c.SetParamValues("")

		// Act
		err = h.Create(c)

		// Assertion
		assert.NoError(t, err)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Test Fail invalid to", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/pockets/1/transfers", strings.NewReader(`{"to": 0, "amount": 100.0}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		row := sqlmock.NewRows([]string{"id", "fromPocketId", "toPocketId", "amount", "date"}).
			AddRow(1, 1, 0, 100.0, now)
		mock.ExpectQuery(cStmt).WithArgs(1, 2, 100.0).WillReturnRows(row)

		h := New(config.FeatureFlag{}, db)
		c := e.NewContext(req, rec)
		c.SetPath("/pockets/:id/transfer")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Act
		err = h.Create(c)

		// Assertion
		assert.NoError(t, err)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Test Fail invalid amount", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/pockets/1/transfers", strings.NewReader(`{"to": 2, "amount": 0.0}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		row := sqlmock.NewRows([]string{"id", "fromPocketId", "toPocketId", "amount", "date"}).
			AddRow(1, 1, 0, 100.0, now)
		mock.ExpectQuery(cStmt).WithArgs(1, 2, 100.0).WillReturnRows(row)

		h := New(config.FeatureFlag{}, db)
		c := e.NewContext(req, rec)
		c.SetPath("/pockets/:id/transfer")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Act
		err = h.Create(c)

		// Assertion
		assert.NoError(t, err)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}
