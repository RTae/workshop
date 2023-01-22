//go:build integration

package pocket

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetPocketIT(t *testing.T) {
	e := echo.New()

	cfg := config.New().All()
	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}

	hPocket := New(sql)
	e.GET("/pockets", hPocket.Get)

	req := httptest.NewRequest(http.MethodGet, "/pockets", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := `[{"id": 2, "amount": 200.00, "name": "test_pocket", "accountId": 2, "currency": "THB"}]`
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}

func TestGetPocketByIDIT(t *testing.T) {
	e := echo.New()

	cfg := config.New().All()
	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}

	hPocket := New(sql)
	e.GET("/pockets/:id", hPocket.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/pockets/2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := `{"id": 2, "amount": 200.00, "name": "test_pocket", "accountId": 2, "currency": "THB"}`
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
