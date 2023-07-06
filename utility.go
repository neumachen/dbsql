package sqlstmt

import (
	"database/sql"
	"reflect"
)

// Database represents an interface for a SQL database connection.
type Database interface {
	// Query executes a query that returns rows.
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow executes a query that is expected to return at most one row.
	QueryRow(query string, args ...any) *sql.Row

	// Exec executes a query without returning any rows.
	Exec(query string, args ...interface{}) (sql.Result, error)

	// Prepare creates a prepared statement for later queries or executions.
	Prepare(query string) (*sql.Stmt, error)

	// Close closes the database connection.
	Close() error
}

func isNilOrZero(v any) bool {
	if v == nil {
		return true
	}

	valueOf := reflect.ValueOf(v)

	switch valueOf.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return valueOf.IsNil() || valueOf.IsZero()
	default:
		return false
	}
}
