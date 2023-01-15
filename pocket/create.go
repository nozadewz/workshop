package pocket

import (
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h handler) CreatePocket(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	var cp Pocket
	err := c.Bind(&cp)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	if cp.Category == "" {
		cp.Category = "Vacation"
	}

	if cp.Account_ID == 0 {
		cp.Account_ID = 1
	}

	var checkBalance float64
	row := h.db.QueryRow(chkBalance, cp.Account_ID)
	err = row.Scan(&checkBalance)
	if err != nil {
		logger.Error("row scan error:", zap.Error(err))
	}

	if checkBalance < cp.Balance {
		logger.Error("bad request not enough money in balance")
		return c.JSON(http.StatusBadRequest, "bad request not enough money in balance")
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, cp.Name, cp.Account_ID, cp.Category, cp.Currency, cp.Balance).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err)
	}

	// sb := checkBalance - cp.Balance
	// var acc_id int64
	// row = h.db.QueryRow(setBalance, sb, cp.Account_ID)
	// if err := row.Scan(&acc_id); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	// }

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	cp.ID = lastInsertId
	return c.JSON(http.StatusCreated, cp)
}

type Err struct {
	Message string `json:"message"`
}
