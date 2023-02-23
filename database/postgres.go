package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"

	"stories/constants"
)

type PostgresQuery struct {
	db       *sql.DB
	instance *Connection
	config   *DatabaseConfig
}

func (repo *PostgresQuery) Initialize(newConfig *DatabaseConfig) bool {
	if repo.config != nil {
		return false
	}
	repo.config = newConfig
	return true
}

func parsePostgresError(value *pq.Error) error {
	uid := generateUid()
	log.Printf("database error %s-%X %s\n%s", value.Code, uid, value.Message, value.Detail)
	log.Println(value.InternalQuery)
	if len(value.Hint) > 0 {
		log.Println(value.Hint)
	}
	switch value.Code {
	case ForeignKeyViolation:
		return fmt.Errorf(constants.SQLErrorForeignKeyViolation)
	case UniqueViolation:
		return fmt.Errorf(constants.SQLErrorDuplicateKey)
	case NotNullViolation:
		return fmt.Errorf(constants.SQLErrorNotNullViolation)
	}
	return fmt.Errorf(constants.SQLErrorPlaceholder, value.Code, uid)
}

func (repo *PostgresQuery) Connect() error {
	var err error
	dsn := fmt.Sprintf(DSN, repo.config.User, repo.config.Password, repo.config.Host, repo.config.Port, repo.config.Name)
	repo.db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	log.Printf("connected to postgresql server")
	// Set charset and encoding
	_, err = repo.db.Exec(SetEncodingQuery)
	if err != nil {
		return err
	}
	repo.instance = &Connection{db: repo.db}
	return nil
}

func (repo *PostgresQuery) GetDatabase() *Connection {
	ctx.do(func() {
		err := repo.Connect()
		if err != nil {
			panic(err)
		}
	})
	return repo.instance
}

func (repo *PostgresQuery) NewTransaction() (*Transaction, error) {
	var connection *Connection
	var db *sql.DB
	var tx *sql.Tx
	var err error
	connection = repo.GetDatabase()
	db = connection.db
	tx, err = db.Begin()
	if err != nil {
		return nil, err
	}
	return &Transaction{tx: tx}, nil
}

func (repo *PostgresQuery) NewQuery() IQuery {
	connection := repo.GetDatabase()
	return &Query{db: connection.db, instance: repo.instance, config: repo.config}
}
