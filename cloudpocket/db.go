package cloudpocket

import (
	"database/sql"

	"github.com/google/uuid"
)

func (h handler) selectPocketById(id int) (Pocket, error) {
	p := Pocket{}
	stmt, err := h.db.Prepare("select * from cloud_pockets  where id=$1")
	if err != nil {
		return p, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&p.Id, &p.Name, &p.Category, &p.Currency, &p.Balance)
	return p, err
}

func transferBalanceAndLog(db *sql.DB, t TransferTxn) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	if err != nil {
		return "", err
	}
	err = updateBalance(tx, t.Src)
	if err != nil {
		return "", err
	}
	err = updateBalance(tx, t.Dest)
	if err != nil {
		return "", err
	}
	txnId, err := logTransferTxn(tx, t)
	if err != nil {
		return "", err
	}
	if err = tx.Commit(); err != nil {
		return "", err
	}
	return txnId, nil
}

func logTransferTxn(tx *sql.Tx, t TransferTxn) (string, error) {
	uuid := uuid.New()
	spTxn := TransactionHistory{
		TransactionId:   uuid.String(),
		CloudPocketId:   t.Src.Id,
		Amount:          t.Amount,
		TransactionType: "debit",
		Description:     t.Description,
	}

	dpTxn := TransactionHistory{
		TransactionId:   uuid.String(),
		CloudPocketId:   t.Dest.Id,
		Amount:          t.Amount,
		TransactionType: "credit",
		Description:     t.Description,
	}

	err := InsertTransactionHistory(tx, spTxn)
	if err != nil {
		return "", err
	}
	err = InsertTransactionHistory(tx, dpTxn)
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

func updateBalance(tx *sql.Tx, p Pocket) error {
	sqlStatement := `UPDATE cloud_pockets SET balance=$1 WHERE id=$2`
	_, err := tx.Exec(sqlStatement, p.Balance, p.Id)
	return err
}

func InsertTransactionHistory(tx *sql.Tx, transaction TransactionHistory) error {
	sqlStatement := `INSERT INTO transaction_history (transaction_id, pocket_id, amount, transaction_type, description) VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(sqlStatement, transaction.TransactionId, transaction.CloudPocketId, transaction.Amount, transaction.TransactionType, transaction.Description)
	return err
}
