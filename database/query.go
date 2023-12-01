package database

import (
	"database/sql"

	"stories/types"
)

// Query is used to run a single queries without transactions
type Query struct {
	db       *sql.DB
	instance *Connection
	config   *DatabaseConfig
}

func (repo *Query) SetSqlDB(db *sql.DB) {
	repo.db = db
}

func (repo *Query) Initialize(newConfig *DatabaseConfig) bool {
	if repo.config != nil {
		return false
	}
	repo.config = newConfig
	return true
}

func (*Query) Connect() error {
	return nil
}

func (repo *Query) GetDatabase() *Connection {
	ctx.do(func() {
		err := repo.Connect()
		if err != nil {
			panic(err)
		}
	})
	return repo.instance
}

// QueryMany queries multiple rows in a table
func (query *Query) QueryMany(statement string, scanner Scanner, values ...interface{}) ([]interface{}, error) {
	return queryMany(query.db, statement, scanner, values...)
}

// QueryManyPaged queries multiple rows in a table abd builds a page
func (query *Query) QueryManyPaged(statement string, scanner Scanner, page int, limit int, values ...interface{}) (types.Page, error) {
	return queryManyPaged(query.db, statement, scanner, page, limit, values...)
}

// QueryOne queries just one record in a table
func (query *Query) QueryOne(statement string, scanner Scanner, values ...interface{}) (interface{}, error) {
	return queryOne(query.db, statement, scanner, values...)
}

// Exec executes arbitrary SQL statements, used for `INSERT/UPDATE/DELETE'
func (query *Query) Exec(statement string, values ...interface{}) error {
	var err error
	_, err = exec(query.db, statement, values...)
	return err
}

// Commit commits a transaction
func (*Query) Commit() error {
	// This is intentionally no-op as single queries are committed automatically
	return nil
}

// Rollback rolls back a transaction
func (*Query) Rollback() {
	// This is intentionally no-op as single queries are committed automatically
}
