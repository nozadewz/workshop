package pocket

import (
	"database/sql"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Pocket struct {
	ID         int64   `json:"id"`
	Account_ID int64   `json:"account_id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	Currency   string  `json:"currency"`
	Balance    float64 `json:"initial_balance"`
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

const (
	cStmt      = "INSERT INTO pockets (name,account_id,category,currency,balance) VALUES ($1,$2,$3,$4,$5) RETURNING id;"
	chkMoney   = "SELECT balance FROM accounts WHERE id = $1"
	setBalance = "UPDATE accounts SET balance = $1 WHERE id = $2 RETURNING id"
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
	row := h.db.QueryRow(chkMoney, cp.Account_ID)
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

	sb := checkBalance - cp.Balance
	var acc_id int64
	row = h.db.QueryRow(setBalance, sb, cp.Account_ID)
	if err := row.Scan(&acc_id); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	cp.ID = lastInsertId
	return c.JSON(http.StatusCreated, cp)
}

type Err struct {
	Message string `json:"message"`
}
