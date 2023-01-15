package pocket

import (
	"database/sql"

	"github.com/kkgo-software-engineering/workshop/config"
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
	chkBalance = "SELECT balance FROM accounts WHERE id = $1"
	// chkBalancePocket = "SELECT balance FROM pockets WHERE account_id = $1"
	// setBalance = "UPDATE accounts SET balance = $1 WHERE id = $2 RETURNING id"
)
