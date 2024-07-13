package dbsql

import (
	"context"
	"database/sql"
	"errors"
)

// mockDB mocks out sqlx.m for the purpose of testing.
type mockDB struct {
	PingOk     bool
	ExecOk     bool
	PrepareOk  bool
	QueryOk    bool
	QueryRowOk bool
}

// func (m *mockDB) ExecContext(ctx context.Context)

// Prepare ...
func (m *mockDB) Prepare(query string) (*sql.Stmt, error) {
	if m.PrepareOk {
		return &sql.Stmt{}, nil
	}
	return nil, errors.New("mock error")
}

func (m *mockDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if m.PrepareOk {
		return &sql.Stmt{}, nil
	}
	return nil, errors.New("mock error")
}

// Exec ...
func (m *mockDB) Exec(query string, args ...any) (sql.Result, error) {
	if m.ExecOk {
		return &Result{}, nil
	}
	return nil, errors.New("mock error")
}

func (m *mockDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if m.ExecOk {
		return &Result{}, nil
	}
	return nil, errors.New("mock error")
}

// Query ...
func (m *mockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.QueryOk {
		return nil, nil
	}
	return nil, errors.New("mock error")
}

func (m *mockDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if m.QueryOk {
		return nil, nil
	}
	return nil, errors.New("mock error")
}

// QueryRow ...
func (m *mockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.QueryRowOk {
		return &sql.Row{}
	}
	return nil
}

func (m *mockDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if m.QueryRowOk {
		return &sql.Row{}
	}
	return nil
}

// Ping ...
func (m *mockDB) Ping() error {
	if m.PingOk {
		return nil
	}
	return errors.New("mock Ping error")
}

// Close ...
func (m *mockDB) Close() error {
	return nil
}

// Result is a mock of sql.Result.
type Result struct {
	LastInsertIDOk bool
}

// LastInsertId ...
func (r *Result) LastInsertId() (int64, error) {
	return 1, nil
}

// RowsAffected ...
func (r *Result) RowsAffected() (int64, error) {
	return 1, nil
}
