package pocket

func (h handler) getPocketById(id int) (Pocket, error) {
	p := Pocket{}
	stmt, err := h.db.Prepare("select * from pockets  where id=$1")
	if err != nil {
		return p, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&p.Id, &p.Name, &p.Category, &p.Currency, &p.Balance)
	return p, err
}

func (h handler) updateBalance(p Pocket) error {

	sqlStatement := `UPDATE pockets SET balance=$1 WHERE id=$2`
	_, err := h.db.Exec(sqlStatement, p.Balance, p.Id)
	return err
}

func (h handler) InsertTransactionHistory(transaction TransactionHistory) error {
	sqlStatement := `INSERT INTO transaction_history (transaction_id, pocket_id, amount, transaction_type, description) VALUES ($1, $2, $3, $4, $5)`
	_, err := h.db.Exec(sqlStatement, transaction.TransactionId, transaction.CloudPocketId, transaction.Amount, transaction.TransactionType, transaction.Description)
	return err
}
