package database

import (
	"database/sql"
	"log"

	"stories/types"
)

// Transaction is used to run multiple queries in a single transaction
type Transaction struct {
	tx        *sql.Tx
	committed bool
}

// QueryManyPaged queries multiple rows in a table abd builds a page
func (transaction *Transaction) QueryManyPaged(query string, scanner Scanner, page int, limit int, values ...interface{}) (types.Page, error) {
	return queryManyPaged(transaction.tx, query, scanner, page, limit, values)
}

// QueryOne queries just one record in a table
func (transaction *Transaction) QueryOne(statement string, scanner Scanner, values ...interface{}) (interface{}, error) {
	return queryOne(transaction.tx, statement, scanner, values...)
}

// QueryMany queries multiple rows in a table
func (transaction *Transaction) QueryMany(statement string, scanner Scanner, values ...interface{}) ([]interface{}, error) {
	return queryMany(transaction.tx, statement, scanner, values...)
}

// Exec executes arbitrary SQL statements, used for `INSERT/UPDATE/DELETE'
func (transaction *Transaction) Exec(statement string, values ...interface{}) error {
	var err error
	_, err = exec(transaction.tx, statement, values...)
	return err
}

// Rollback rolls back a transaction
func (transaction *Transaction) Rollback() {
	var tx *sql.Tx
	var err error
	if transaction.committed {
		// Rolling back after having committed is a no-op because
		// we could do
		//
		// defer tx.Rollback()
		//
		return
	}
	tx = transaction.tx
	err = tx.Rollback()
	if err != nil {
		log.Println(err)
	}
}

// Commit commits a transaction
func (transaction *Transaction) Commit() error {
	var tx *sql.Tx
	var err error
	tx = transaction.tx
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}
	transaction.committed = true
	return nil
}
