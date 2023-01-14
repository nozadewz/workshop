package cloud_pockets

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/labstack/echo/v4"
)

type TransferRequest struct {
	SourceCloudPocketId      string  `json:"source_cloud_pocket_id"`
	DestinationCloudPocketId string  `json:"destination_cloud_pocket_id"`
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
	var transferRequest TransferRequest
	err := c.Bind(&transferRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if transferRequest.Amount < 0 {
		return c.JSON(http.StatusBadRequest, "Amount must be greater than 0")
	}
	srcPocket := Pocket{}
	destPocket := Pocket{}

	srcPocketIdRequest, _ := strconv.Atoi(transferRequest.SourceCloudPocketId)
	destPocketIdRequest, _ := strconv.Atoi(transferRequest.DestinationCloudPocketId)

	srcPocket, _ = h.getPocketById(srcPocketIdRequest, c)
	destPocket, _ = h.getPocketById(destPocketIdRequest, c)

	balance_srcPocket := srcPocket.Balance - transferRequest.Amount
	if balance_srcPocket < 0 {
		return c.JSON(http.StatusBadRequest, "Not enough money in source cloud pocket")
	}
	srcPocket.Balance = balance_srcPocket

	balance_destPocket := destPocket.Balance + transferRequest.Amount
	destPocket.Balance = balance_destPocket

	err = h.updateBalance(srcPocket)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = h.updateBalance(destPocket)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	txn, err := h.logTransferTxn(srcPocket, destPocket, transferRequest.Amount, transferRequest.Description)
	resp := TransferResponse{
		TransactionId:          txn,
		SourceCloudPocket:      srcPocket,
		DestinationCloudPocket: destPocket,
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
	src_txn := TransactionHistory{
		TransactionId:   uuid.String(),
		CloudPocketId:   src.Id,
		Amount:          amount,
		TransactionType: "debit",
		Description:     desc,
	}

	dest_txn := TransactionHistory{
		TransactionId:   uuid.String(),
		CloudPocketId:   dest.Id,
		Amount:          amount,
		TransactionType: "credit",
		Description:     desc,
	}

	err := h.InsertTransactionHistory(src_txn)
	if err != nil {
		return src_txn.TransactionId, err
	}
	err = h.InsertTransactionHistory(dest_txn)
	if err != nil {
		return src_txn.TransactionId, err
	}

	return src_txn.TransactionId, err
}
