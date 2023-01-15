package pocket

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TransferRequest struct {
	SourceCloudPocketId      int64   `json:"source_cloud_pocket_id"`
	DestinationCloudPocketId int64   `json:"destination_cloud_pocket_id"`
	Amount                   float64 `json:"amount"`
	Description              string  `json:"description"`
}
type TransferResponse struct {
	TransactionId          string `json:"transaction_id"`
	SourceCloudPocket      Pocket `json:"source_cloud_pocket"`
	DestinationCloudPocket Pocket `json:"destination_cloud_pocket"`
	Status                 string `json:"status"`
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

	newSpBal := SubBalance(sp.Balance, tfr.Amount)
	if newSpBal < 0 {
		logger.Error("Insufficient balance")
		return c.JSON(http.StatusBadRequest, Err{Message: "Insufficient balance"})
	}

	dp.Balance = AddBalance(dp.Balance, tfr.Amount)

	//print dp.Balance
	fmt.Println("dp.Balance:", dp.Balance)

	txn := TransferTxn{
		Src:         sp,
		Dest:        dp,
		Amount:      tfr.Amount,
		Description: tfr.Description,
	}

	txnId, err := transferBalanceAndLog(h.db, txn)
	if err != nil {
		logger.Error("Error while transferring balance", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	resp := TransferResponse{
		TransactionId:          txnId,
		SourceCloudPocket:      sp,
		DestinationCloudPocket: dp,
		Status:                 "Success",
	}

	return c.JSON(http.StatusOK, resp)
}
func AddBalance(amount1 float64, amount2 float64) float64 {
	bigTfr, _ := new(big.Float).SetPrec(200).SetString(fmt.Sprintf("%f", amount1))
	bigDpBalance, _ := new(big.Float).SetPrec(200).SetString(fmt.Sprintf("%f", amount2))
	bigNewDpBal := new(big.Float).Add(bigDpBalance, bigTfr)
	result, _ := bigNewDpBal.Float64()
	return result
}

func SubBalance(amount1 float64, amount2 float64) float64 {
	bigTfr, _ := new(big.Float).SetPrec(200).SetString(fmt.Sprintf("%f", amount1))
	bigDpBalance, _ := new(big.Float).SetPrec(200).SetString(fmt.Sprintf("%f", amount2))
	bigNewDpBal := new(big.Float).Sub(bigTfr, bigDpBalance)
	result, _ := bigNewDpBal.Float64()
	return result
}

type TransactionHistory struct {
	TransactionId   string  `json:"transaction_id"`
	CloudPocketId   int64   `json:"cloud_pocket_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	Description     string  `json:"description"`
	CreatedAt       string  `json:"created_at"`
}

type TransferTxn struct {
	Src         Pocket
	Dest        Pocket
	Amount      float64
	Description string
}
