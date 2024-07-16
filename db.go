package dbsql

import (
	"context"
	"database/sql"
)

// NOTE: Benefits of Using Multiple Interfaces in Go
// 1. It promotes the Interface Segregation Principle (ISP), allowing clients to depend only on the methods they need.
// 2. It increases modularity and flexibility, making it easier to swap implementations or mock specific behaviors in tests.
// 3. It improves code readability and maintainability by clearly defining the responsibilities of each interface.
// 4. It allows for easier composition of interfaces, as demonstrated by the DBPreparerExecutor and DB interfaces.
// 5. It provides better abstraction, hiding unnecessary details from clients that only need specific functionality.

// DBPreparer defines an interface for creating prepared statements. It mirrors the core functionality of
// database/sql.DB, providing both context-aware and non-context methods for preparing SQL statements.
type DBPreparer interface {
	// Prepare creates a prepared statement for later queries or executions.
	// It returns a sql.Stmt and an error, if any.
	Prepare(query string) (*sql.Stmt, error)

	// PrepareContext creates a prepared statement for later queries or executions.
	// It accepts a context.Context for cancellation and timeout control.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// DBExecutor defines an interface for executing SQL operations.
// It mirrors the core functionality of database/sql.DB, providing both
// context-aware and non-context methods for database operations.
type DBExecutor interface {
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

// DBPreparerExecutor defines an interface for preparing and executing SQL operations.
// It combines the Preparer and Executor interfaces, providing methods for
// creating prepared statements and executing SQL operations.
type DBPreparerExecutor interface {
	// Preparer provides methods for creating prepared statements.
	DBPreparer
	// Executor provides methods for executing SQL operations.
	DBExecutor
}

type DBCloser interface {
	// Close closes the database, releasing any open resources.
	Close() error
}

type DBPinger interface {
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	Ping() error
}

// DB ...
type DB interface {
	DBExecutor
	DBPreparer
	DBCloser
	DBPinger
}
