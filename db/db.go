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
	createTb := `CREATE TABLE IF NOT EXISTS pockets (
		id SERIAL primary key ,
		name text,
		category text,
		currency text,
		balance float8
	);`
	_, err := db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}
	return db

}
