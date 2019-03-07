package main
import (
	"database/sql"
	_"github.com/mattn/go-sqlite3"
)
type Transaction interface {
	Exec(query string, args ...interface{})(sql.Result,error)
	Prepare(query string)(*sql.Stmt, error)
	Query(query string, args ...interface{})(*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
type TxFn func(Transaction)error
