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

func TestGetById(t *testing.T) {

	//Arrange
	ex := Pocket{
		ID:         1,
		Account_ID: 1,
		Name:       "Travel Fund",
		Category:   "Vacation",
		Currency:   "THB",
		Balance:    50,
	}
	cfgFlag := config.FeatureFlag

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", strings.NewReader(tc.reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	newsMockRows := sqlmock.NewRows([]string{"id", "account_id", "name", "category", "currency", "balance"}).
		AddRow(ex.ID, ex.Account_ID, ex.Name, ex.Category, ex.Currency, ex.Balance)

	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT (.+) FROM pockets").WillReturnRows(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := New(cfgFlag, db)

	//Act
	err = h.GetExpensesHandler(c)

	actual := Expenses{}
	err = util.ConvertToStruct(rec, &actual)
	if err != nil {
		t.Errorf("Test Failed because: %v", err)
	}

	//Assert
	assert.NoError(t, err)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, ex, actual)
	}
}
