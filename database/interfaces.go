package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"stories/constants"
	"stories/types"
)

const (
	DSN              = "user=%s password=%s host=%s port=%d dbname=%s sslmode=disable"
	SetEncodingQuery = "SET CLIENT_ENCODING TO 'UTF8'"
)

var (
	NegativeOffsetErr = errors.New(constants.NegativeOffsetInPagedQuery)
	PageEmptyErr      = errors.New(constants.EmptyPageInPagedQuery)
)

// IQuery interface will help using a `Query' and a `Transaction' seamlessly
type IQuery interface {
	// QueryOne queries just one record in a table
	QueryOne(query string, scanner Scanner, values ...interface{}) (interface{}, error)
	// QueryMany queries multiple rows in a table
	QueryMany(query string, scanner Scanner, values ...interface{}) ([]interface{}, error)
	// QueryManyPaged queries multiple rows in a table abd builds a page
	QueryManyPaged(query string, scanner Scanner, page int, limit int, values ...interface{}) (types.Page, error)
	// Exec executes arbitrary SQL statements, used for `INSERT/UPDATE/DELETE'
	Exec(query string, values ...interface{}) error
	// Commit commits a transaction
	Commit() error
	// Rollback rolls back a transaction
	Rollback()
	Connect() error
	Initialize(newConfig *DatabaseConfig) bool
	GetDatabase() *Connection
	SetSqlDB(db *sql.DB)
	NewQuery() IQuery
}

type RowInterface interface {
	Scan(value ...interface{}) error
}

type NativeQuery interface {
	Query(query string, values ...interface{}) (*sql.Rows, error)
	QueryRow(query string, values ...interface{}) *sql.Row
	Exec(query string, values ...interface{}) (sql.Result, error)
}

type Connection struct {
	db *sql.DB
}

type once struct {
	mutex *sync.Mutex
	done  uint32
}

var mutex sync.Mutex

var ctx = &once{mutex: &mutex, done: 0}

func (synced *once) do(fn func()) {
	if atomic.LoadUint32(&synced.done) == 1 {
		return
	}
	mutex := synced.mutex
	// get the lock when it's free
	mutex.Lock()
	defer mutex.Unlock()
	if synced.done == 0 {
		defer atomic.StoreUint32(&synced.done, 1)
		// Finally call the function
		fn()
	}
}

func closeRows(rows *sql.Rows) {
	var err error
	err = rows.Close()
	if err != nil {
		log.Println(err)
	}
}

func generateUid() int64 {
	now := time.Now()
	return now.Unix()
}

const (
	ForeignKeyViolation pq.ErrorCode = "23503"
	UniqueViolation     pq.ErrorCode = "23505"
	NotNullViolation    pq.ErrorCode = "23502"
)

func parseSqlError(err error) error {
	switch value := err.(type) {
	case *pq.Error:
		return parsePostgresError(value)
	default:
		return err
	}
}

func queryMany(native NativeQuery, query string, scanner Scanner, values ...interface{}) ([]interface{}, error) {
	var rows *sql.Rows
	var err error
	var result []interface{}
	rows, err = native.Query(query, values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, parseSqlError(err)
	}
	defer closeRows(rows)
	// Always return an array, even if empty
	result = make([]interface{}, 0)
	for rows.Next() {
		var value interface{}
		var err error
		value, err = scanner.ScanRow(rows)
		if err != nil {
			return nil, parseSqlError(err)
		}
		result = append(result, value)
	}
	return result, parseSqlError(rows.Err())
}

func queryManyPaged(native NativeQuery, query string, scanner Scanner, pageNumber int, limit int, values ...interface{}) (types.Page, error) {
	var rows *sql.Rows
	var err error
	var result []interface{}
	var slice []interface{}
	var previous types.Int
	var next types.Int
	var offset int
	var pageCount int
	var page types.Page
	if pageNumber < 1 {
		return types.Page{}, NegativeOffsetErr
	}
	offset = (pageNumber - 1) * limit
	rows, err = native.Query(query, values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Page{}, err
		}
		return types.Page{}, parseSqlError(err)
	}
	defer closeRows(rows)
	for rows.Next() {
		var value interface{}
		var err error
		value, err = scanner.ScanRow(rows)
		if err != nil {
			return types.EmptyPage, parseSqlError(err)
		}
		result = append(result, value)
	}
	if len(result) == 0 {
		if pageNumber == 1 {
			return types.EmptyPage, nil
		}
		return types.EmptyPage, PageEmptyErr
	}
	if pageNumber > 1 {
		previous = types.NewInt(pageNumber - 1)
	}
	pageCount = len(result) / limit
	if len(result)%limit > 0 {
		pageCount += 1
	}
	slice = result
	if len(slice) > offset {
		slice = slice[offset:]
	}
	if len(slice) > limit {
		slice = slice[:limit]
		next = types.NewInt(pageNumber + 1)
	}
	page = types.Page{
		PageCount:    pageCount,
		NextPage:     next,
		PreviousPage: previous,
		Data:         slice,
	}
	page.PageNumber = pageNumber
	page.PageSize = limit
	page.TotalCount = len(result)
	return page, parseSqlError(rows.Err())
}

func queryOne(native NativeQuery, query string, scanner Scanner, values ...interface{}) (interface{}, error) {
	var err error
	var rows *sql.Rows
	var value interface{}
	// Perform the queries
	rows, err = native.Query(query, values...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return value, err
		}
		return value, parseSqlError(err)
	}
	defer func() {
		_ = rows.Close()
	}()
	// We only care about 1 row
	if !rows.Next() {
		return value, sql.ErrNoRows
	}
	// Scan now (we ignore the return value because it's set int the scanner)
	value, err = scanner.ScanRow(rows)
	if rows.Next() {
		return value, fmt.Errorf("oops, found more than a single row in a 1 row queries")
	}
	return value, parseSqlError(err)
}

func exec(native NativeQuery, query string, values ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error
	result, err = native.Exec(query, values...)
	return result, parseSqlError(err)
}
