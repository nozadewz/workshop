package pocket

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TransferRequest struct {
	SourceCloudPocketId      int     `json:"source_cloud_pocket_id"`
	DestinationCloudPocketId int     `json:"destination_cloud_pocket_id"`
	Amount                   float64 `json:"amount"`
	Description              string  `json:"description"`
}

type Pocket struct {
	Id       int     `json:"id"`
	Name     string  `json:"title"`
	Category string  `json:"amount"`
	Currency string  `json:"-"`
	Balance  float64 `json:"balance"`
}

type TransferResponse struct {
	TransactionId          string `json:"transaction_id"`
	SourceCloudPocket      Pocket `json:"source_cloud_pocket"`
	DestinationCloudPocket Pocket `json:"destination_cloud_pocket"`
	Status                 string `json:"status"`
}
type Err struct {
	Message string `json:"message"`
}

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

func (h handler) Transfer(c echo.Context) error {
	logger := mlog.L(c)
	var tfr TransferRequest
	err := c.Bind(&tfr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if tfr.Amount < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Amount must be greater than 0"})
	}

	var sp Pocket
	var dp Pocket

	sp, err = h.getPocketById(tfr.SourceCloudPocketId)
	if err != nil {
		logger.Error("Source pocket not found", zap.Error(err))
		return c.JSON(http.StatusBadRequest, Err{Message: "Source pocket not found"})
	}

	dp, err = h.getPocketById(tfr.DestinationCloudPocketId)
	if err != nil {
		logger.Error("Destination pocket not found", zap.Error(err))
		return c.JSON(http.StatusBadRequest, Err{Message: "Destination pocket not found"})
	}

	newSpBal := sp.Balance - tfr.Amount
	if newSpBal < 0 {
		logger.Error("Insufficient balance")
		return c.JSON(http.StatusBadRequest, Err{Message: "Insufficient balance"})
	}

	sp.Balance = newSpBal
	newDpBal := dp.Balance + tfr.Amount
	dp.Balance = newDpBal

	err = h.updateBalance(sp)
	if err != nil {
		logger.Error("updateSourceBalance Error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}
	err = h.updateBalance(dp)
	if err != nil {
		logger.Error("updateDestinationBalance Error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	txn, err := h.logTransferTxn(sp, dp, tfr.Amount, tfr.Description)
	if err != nil {
		logger.Error("logTransferTxn Error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, Err{Message: "Internal Server Error"})
	}

	resp := TransferResponse{
		TransactionId:          txn,
		SourceCloudPocket:      sp,
		DestinationCloudPocket: dp,
		Status:                 "Success",
	}

	return c.JSON(http.StatusOK, resp)
}

type TransactionHistory struct {
	TransactionId   string  `json:"transaction_id"`
	CloudPocketId   int     `json:"cloud_pocket_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	Description     string  `json:"description"`
	CreatedAt       string  `json:"created_at"`
}

func (h handler) logTransferTxn(src Pocket, dest Pocket, amount float64, desc string) (string, error) {
	uuid := uuid.New()
	spTxn := TransactionHistory{
		TransactionId:   uuid.String(),
		CloudPocketId:   src.Id,
		Amount:          amount,
		TransactionType: "debit",
		Description:     desc,
	}

	dpTxn := TransactionHistory{
		TransactionId:   uuid.String(),
		CloudPocketId:   dest.Id,
		Amount:          amount,
		TransactionType: "credit",
		Description:     desc,
	}

	err := h.InsertTransactionHistory(spTxn)
	if err != nil {
		return spTxn.TransactionId, err
	}
	err = h.InsertTransactionHistory(dpTxn)
	if err != nil {
		return spTxn.TransactionId, err
	}

	return spTxn.TransactionId, err
}
