//go:build unit

package api

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

func TestGetAll(t *testing.T) {

	tests := []struct {
		name       string
		sqlFn      func() (*sql.DB, error)
		reqBody    string
		wantStatus int
		wantBody   string
	}{
		{"get all pockets succesfully",
			func() (*sql.DB, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, err
				}
				rows := mock.NewRows([]string{"id", "name", "category", "currency", "balance"}).
					AddRow(12345, "Travel Fund", "Vacation", "THB", 100).
					AddRow(67890, "Savings", "Emergency Fund", "THB", 200)

				mock.ExpectQuery("SELECT").WillReturnRows(rows)
				return db, err
			},
			"",
			http.StatusOK,
			`{
				"id": "12345",
				"name": "Travel Fund",
				"category": "Vacation",
				"currency": "THB",
				"balance": 100
			},
			{
				"id": "67890",
				"name": "Savings",
				"category": "Emergency Fund",
				"currency": "THB",
				"balance": 200
			}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/cloud-pockets", strings.NewReader(tc.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			db, err := tc.sqlFn()
			h := New(db)
			// Assertions
			assert.NoError(t, err)
			if assert.NoError(t, h.GetAllPockets(c)) {
				assert.Equal(t, tc.wantStatus, rec.Code)
				assert.JSONEq(t, tc.wantBody, rec.Body.String())
			}
		})
	}
}

// func seedPocket(t *testing.T) Pocket {
// 	var c Pocket
// 	body := bytes.NewBufferString(`{
// 		"id": "12345",
// 		"name": "Travel Fund",
// 		"category": "Vacation",
// 		"currency": "THB",
// 		"balance": 100
// 	},
// 	{
// 		"id": "67890",
// 		"name": "Savings",
// 		"category": "Emergency Fund",
// 		"currency": "THB",
// 		"balance": 200
// 	}`)

// 	req, _ := http.NewRequest(http.MethodPost, "/cloud-pockets", body)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

// 	client := http.Client{}
// 	res, err := client.Do(req).Decode(&c)

// 	if err != nil {
// 		t.Fatal("can't create uomer:", err)
// 	}
// 	return c
// }
