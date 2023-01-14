package pocket

import 
(
	"database/sql"
	"github.com/kkgo-software-engineering/workshop/config"
)
type Pocket struct {
	ID     int      `json:"id"`
	Balance  float64   `json:"balance"`
	Currency string `json:"currency"`
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

type Err struct {
	Message string `json:"message"`
}