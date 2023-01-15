package db

import (
	"database/sql"
	"log"
)

func MigrationTransactionHistory(db *sql.DB) *sql.DB {
	// Create table
	createTb := `CREATE TABLE IF NOT EXISTS transaction_history (
		transaction_id text,
		pocket_id integer,
		amount decimal,
		transaction_type text,
		description text,
		created_at timestamp default now(),
		PRIMARY KEY(transaction_id, pocket_id),
		CONSTRAINT fk_pocket_id FOREIGN KEY(pocket_id) REFERENCES pockets(id)
	);`
	_, err := db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}
	return db

}

func MigrationCloudPocket(db *sql.DB) *sql.DB {
	// Create table
	createTBPocket := `CREATE TABLE IF NOT EXISTS pockets (id SERIAL PRIMARY KEY, account_id INT, name TEXT, category TEXT, currency TEXT, balance FLOAT, CONSTRAINT fk_account_id FOREIGN KEY(account_id) REFERENCES accounts(id))`
	_, err := db.Exec(createTBPocket)
	if err != nil {
		log.Fatal("can't create table pockets", err)
	}
	return db

}

func MigrationAccount(db *sql.DB) *sql.DB {
	// Create table
	createTBAccount := `CREATE TABLE IF NOT EXISTS accounts (id SERIAL PRIMARY KEY, balance FLOAT)`

	_, err := db.Exec(createTBAccount)
	if err != nil {
		log.Fatal("can't create table accounts", err)
	}
	return db
}
