package dbsql

import (
	"context"
	"database/sql"
)

// DbPreparer defines an interface for creating prepared statements. It mirrors the core functionality of
// database/sql.DB, providing both context-aware and non-context methods for preparing SQL statements.
type DbPreparer interface {
	// Prepare creates a prepared statement for later queries or executions.
	// It returns a sql.Stmt and an error, if any.
	Prepare(query string) (*sql.Stmt, error)

	// PrepareContext creates a prepared statement for later queries or executions.
	// It accepts a context.Context for cancellation and timeout control.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// DbExecutor defines an interface for executing SQL operations.
// It mirrors the core functionality of database/sql.DB, providing both
// context-aware and non-context methods for database operations.
type DbExecutor interface {
	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	// It returns a sql.Result summarizing the effect of the statement.
	Exec(query string, args ...any) (sql.Result, error)

	// ExecContext executes a query without returning any rows.
	// It accepts a context.Context for cancellation and timeout control.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	// Query executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	Query(query string, args ...any) (*sql.Rows, error)

	// QueryContext executes a query that returns rows, typically a SELECT.
	// It accepts a context.Context for cancellation and timeout control.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	QueryRow(query string, args ...any) *sql.Row

	// QueryRowContext executes a query that is expected to return at most one row.
	// It accepts a context.Context for cancellation and timeout control.
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// DbPreparerExecutor defines an interface for preparing and executing SQL operations.
// It combines the Preparer and Executor interfaces, providing methods for
// creating prepared statements and executing SQL operations.
type DbPreparerExecutor interface {
	// Preparer provides methods for creating prepared statements.
	DbPreparer
	// Executor provides methods for executing SQL operations.
	DbExecutor
}

type DbCloser interface {
	// Close closes the database, releasing any open resources.
	Close() error
}

type DbPinger interface {
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	Ping() error
}

// Db ...
type Db interface {
	DbExecutor
	DbPreparer
	DbCloser
	DbPinger
}
