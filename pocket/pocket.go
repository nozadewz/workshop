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
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

const (
	cStmt = "INSERT INTO pockets (name,category,currency,balance) VALUES ($1,$2,$3,$4) RETURNING id;"
	//cBalanceLimit = 10000
)

// var (
// 	hErrBalanceLimitExceed = echo.NewHTTPError(http.StatusBadRequest,
// 		"create account balance exceed limitation")
// )

func (h handler) CreatePocket(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	var cp Pocket
	err := c.Bind(&cp)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	// if h.cfg.IsLimitMaxBalanceOnCreate && cp.Balance > cBalanceLimit {
	// 	logger.Error("account limit on account creating", zap.Error(hErrBalanceLimitExceed))
	// 	return hErrBalanceLimitExceed
	// }

	var lastInsertId int64
	if cp.Category == "" {
		cp.Category = "Vacation"
	}
	err = h.db.QueryRowContext(ctx, cStmt, cp.Name, cp.Category, cp.Currency, cp.Balance).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return err
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	cp.ID = lastInsertId
	return c.JSON(http.StatusCreated, cp)
}
