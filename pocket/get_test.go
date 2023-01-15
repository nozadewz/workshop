package pocket

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetExpenseByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "account_id", "name", "category", "currency", "initial_balance"}).
		AddRow(1, 1, "Travel", "Vacation", "THB", 1000)
	mock.ExpectQuery("SELECT (.+) FROM pockets WHERE id = (.+)").WillReturnRows(rows)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/cloud-pockets/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := New(config.FeatureFlag{}, db)

	expected := `{"id":1,"account_id":1,"name":"Travel","category":"Vacation","currency":"THB","initial_balance":1000}`

	if assert.NoError(t, h.GetPocketBalanceById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}
