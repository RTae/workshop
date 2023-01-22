//go:build integration

package pocket

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreatePocketIT(t *testing.T) {
	e := echo.New()

	cfg := config.New().All()
	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}

	hPocket := New(sql)

	e.POST("/pockets", hPocket.Create)

	reqBody := `{"amount": 0.0, "name": "test", "accountId": 1, "currency": "THB"}`
	req := httptest.NewRequest(http.MethodPost, "/pockets", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := `{"id": 2, "amount": 0.0, "name": "test", "accountId": 1, "currency": "THB"}`
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
