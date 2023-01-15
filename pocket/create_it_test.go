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
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreatePocketIT(t *testing.T) {
	e := echo.New()

	cfg := config.New().All()
	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}
	cfgFlag := config.FeatureFlag{}

	hAccount := New(cfgFlag, sql)

	e.POST("/cloud-pockets", hPocket.CreatePocket)

	reqBody := `{
			"name": "Travel Fund",
			"currency": "THB",
			"balance": 999.0
		}`
	req := httptest.NewRequest(http.MethodPost, "/cloud-pockets", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := `{
		"id": 1,
		"name": "Travel Fund",
		"category": "Vacation",
		"currency": "THB",
		"balance": 999
	}`
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
