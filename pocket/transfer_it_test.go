//go:build integration

package pocket

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/db"
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
	db.MigrationCloudPocket(sql)
	db.MigrationTransactionHistory(sql)
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

	byteBody, err := ioutil.ReadAll(rec.Body)
	assert.NoError(t, err)

	var tfrs TransferResponse
	json.Unmarshal(byteBody, &tfrs)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, tfrs.TransactionId)
	//assert source cloud pocket
	assert.Equal(t, 1, tfrs.SourceCloudPocket.Id)
	assert.Equal(t, "Travel Fund", tfrs.SourceCloudPocket.Name)
	assert.Equal(t, "Vacation", tfrs.SourceCloudPocket.Category)
	assert.Equal(t, 150.00, tfrs.SourceCloudPocket.Balance)
	//assert destination cloud pocket
	assert.Equal(t, 2, tfrs.DestinationCloudPocket.Id)
	assert.Equal(t, "Savings", tfrs.DestinationCloudPocket.Name)
	assert.Equal(t, "Emergency Fund", tfrs.DestinationCloudPocket.Category)
	assert.Equal(t, 150.00, tfrs.DestinationCloudPocket.Balance)
}
