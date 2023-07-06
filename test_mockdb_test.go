package sqlstmt

import (
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

// Prepare ...
func (m *mockDB) Prepare(query string) (*sql.Stmt, error) {
	if m.PrepareOk {
		return &sql.Stmt{}, nil
	}
	return nil, errors.New("mock error")
}

// Exec ...
func (m *mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
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

// QueryRow ...
func (m *mockDB) QueryRow(query string, args ...interface{}) *sql.Row {
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
