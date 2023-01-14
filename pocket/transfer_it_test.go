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

func TestPocketTransfer(t *testing.T) {
	e := echo.New()

	cfg := config.New().All()
	sql, err := sql.Open("postgres", cfg.DBConnection)
	if err != nil {
		t.Error(err)
	}
	cfgFlag := config.FeatureFlag{}

	hPocket := New(cfgFlag, sql)

	e.POST("/cloud-pockets/transfer", hPocket.Transfer)

	// need to be post pocket
	hPocket.db.Exec(`INSERT INTO pockets(name, category, currency, balance) VALUES ('Travel Fund', 'Vacation', 'THB', 200), ('Savings', 'Emergency Fund', 'THB', 100);`)

	reqBody := `{
		"source_cloud_pocket_id": 1,
		"destination_cloud_pocket_id": 2,
		"amount": 50.00,
		"description": "Transfer from Travel fund to savings"
	}`
	req := httptest.NewRequest(http.MethodPost, "/cloud-pockets/transfer", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := `{"id": 1, "balance": 999.99}`
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
