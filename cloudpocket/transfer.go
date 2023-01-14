package cloudpocket

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
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

type handler struct {
	cfg config.FeatureFlag
	db  *sql.DB
}

func New(cfgFlag config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfgFlag, db}
}

func (h handler) Transfer(c echo.Context) error {
	var tfr TransferRequest
	err := c.Bind(&tfr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if tfr.Amount < 0 {
		return c.JSON(http.StatusBadRequest, "Amount must be greater than 0")
	}
	sp := Pocket{}
	dp := Pocket{}

	sp, _ = h.getPocketById(tfr.SourceCloudPocketId, c)
	dp, _ = h.getPocketById(tfr.DestinationCloudPocketId, c)

	newSpBal := sp.Balance - tfr.Amount
	if newSpBal < 0 {
		return c.JSON(http.StatusBadRequest, "Not enough money in source cloud pocket")
	}
	sp.Balance = newSpBal

	newDpBal := dp.Balance + tfr.Amount
	dp.Balance = newDpBal

	err = h.updateBalance(sp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = h.updateBalance(dp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	txn, err := h.logTransferTxn(sp, dp, tfr.Amount, tfr.Description)
	resp := TransferResponse{
		TransactionId:          txn,
		SourceCloudPocket:      sp,
		DestinationCloudPocket: dp,
		Status:                 "Success",
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
func (h handler) getPocketById(id int, c echo.Context) (Pocket, error) {
	pocket, err := h.selectPocketById(id)
	if err != nil {
		fmt.Printf("getPocketById error : %v", err)
		return pocket, c.JSON(http.StatusBadRequest, err.Error())
	}
	return pocket, err
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
