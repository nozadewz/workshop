package pocket

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTransfer(t *testing.T) {

	t.Run("Transfer Success", func(t *testing.T) {
		e := echo.New()
		body := bytes.NewBufferString(`{
			"source_cloud_pocket_id": 1,
			"destination_cloud_pocket_id": 2,
			"amount": 50.00,
			"description":"Transfer from Travel fund to savings"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/cloud-pockets/transfer", body)
		req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		col := []string{"id", "name", "category", "currency", "balance"}
		mock.ExpectPrepare("select (.+) from cloud_pockets").
			ExpectQuery().WithArgs(1).
			WillReturnRows(sqlmock.NewRows(col).AddRow(1, "apocket", "A", "THB", 100.0))
		mock.ExpectPrepare("select (.+) from cloud_pockets").
			ExpectQuery().WithArgs(2).
			WillReturnRows(sqlmock.NewRows(col).AddRow(2, "bpocket", "B", "THB", 50.0))
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE cloud_pockets SET (.+)").WithArgs(50.0, 1).WillReturnResult(driver.RowsAffected(1))
		mock.ExpectExec("UPDATE cloud_pockets SET (.+)").WithArgs(100.0, 2).WillReturnResult(driver.RowsAffected(1))
		mock.ExpectExec("INSERT INTO transaction_history (.+)").WillReturnResult(driver.RowsAffected(1))
		mock.ExpectExec("INSERT INTO transaction_history (.+)").WillReturnResult(driver.RowsAffected(1))
		mock.ExpectCommit()

		h := handler{db: db}
		err = h.Transfer(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			res := TransferResponse{}
			json.Unmarshal(rec.Body.Bytes(), &res)
			assert.Equal(t, "Success", res.Status)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
	t.Run("Transfer but update balance in DB failed return HTTP Internal error", func(t *testing.T) {
		e := echo.New()
		body := bytes.NewBufferString(`{
			"source_cloud_pocket_id": 1,
			"destination_cloud_pocket_id": 2,
			"amount": 50.00,
			"description":"Transfer from Travel fund to savings"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/cloud-pockets/transfer", body)
		req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		col := []string{"id", "name", "category", "currency", "balance"}
		mock.ExpectPrepare("select (.+) from cloud_pockets").
			ExpectQuery().WithArgs(1).
			WillReturnRows(sqlmock.NewRows(col).AddRow(1, "apocket", "A", "THB", 100.0))
		mock.ExpectPrepare("select (.+) from cloud_pockets").
			ExpectQuery().WithArgs(2).
			WillReturnRows(sqlmock.NewRows(col).AddRow(2, "bpocket", "B", "THB", 50.0))
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE cloud_pockets SET (.+)").WithArgs(50.0, 1).WillReturnResult(driver.RowsAffected(1))
		mock.ExpectExec("UPDATE cloud_pockets SET (.+)").WithArgs(100.0, 2).WillReturnError(driver.ErrBadConn)
		mock.ExpectRollback()

		h := handler{db: db}
		err = h.Transfer(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Transfer but insert txn history in DB failed return HTTP Internal error", func(t *testing.T) {
		e := echo.New()
		body := bytes.NewBufferString(`{
			"source_cloud_pocket_id": 1,
			"destination_cloud_pocket_id": 2,
			"amount": 50.00,
			"description":"Transfer from Travel fund to savings"
		}`)
		req := httptest.NewRequest(http.MethodPost, "/cloud-pockets/transfer", body)
		req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		col := []string{"id", "name", "category", "currency", "balance"}
		mock.ExpectPrepare("select (.+) from cloud_pockets").
			ExpectQuery().WithArgs(1).
			WillReturnRows(sqlmock.NewRows(col).AddRow(1, "apocket", "A", "THB", 100.0))
		mock.ExpectPrepare("select (.+) from cloud_pockets").
			ExpectQuery().WithArgs(2).
			WillReturnRows(sqlmock.NewRows(col).AddRow(2, "bpocket", "B", "THB", 50.0))
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE cloud_pockets SET (.+)").WithArgs(50.0, 1).WillReturnResult(driver.RowsAffected(1))
		mock.ExpectExec("UPDATE cloud_pockets SET (.+)").WithArgs(100.0, 2).WillReturnResult(driver.RowsAffected(1))
		mock.ExpectExec("INSERT INTO transaction_history (.+)").WillReturnError(assert.AnError)
		mock.ExpectRollback()

		h := handler{db: db}
		err = h.Transfer(c)

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}
