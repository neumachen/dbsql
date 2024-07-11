package sqldb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neumachen/sqldb/internal"
)

// NamedParameterPositions is a struct that represents the positions of parameters in an SQL statement.
// The parameterPositions field is a map that stores the positions of parameters, where the key is the parameter name
// and the value is a slice of integers representing the positions.
// The totalPositions field is an integer representing the total number of parameter positions.
type NamedParameterPositions struct {
	parameterPositions map[string][]int
	totalPositions     int
}

// getPositions is a method of the NamedParameterPositions struct.
// It takes a parameter name as input and returns a slice of integers representing the positions of the parameter.
func (n NamedParameterPositions) getPositions(parameter string) []int {
	if internal.IsNilOrZeroValue(n.parameterPositions) {
		return nil
	}

	v, ok := n.parameterPositions[parameter]
	if !ok {
		return nil
	}

	return v
}

// insert is a method of the NamedParameterPositions struct.
// It takes a parameter name and a position as input and inserts the position into the parameterPositions map.
func (n *NamedParameterPositions) insert(parameter string, position int) {
	if internal.IsNilOrZeroValue(n.parameterPositions) {
		n.parameterPositions = make(map[string][]int)
	}

	v, ok := n.parameterPositions[parameter]
	if !ok {
		v = make([]int, 0, 1)
	}

	n.parameterPositions[parameter] = append(v, position)
	n.totalPositions++
}

// BoundNamedParameterValues is a type alias for a slice of any (an empty interface).
// It represents the positional parameters in an SQL statement.
type BoundNamedParameterValues []any

// PreparedStatement is an interface that represents an SQL query with named parameters.
// It defines several methods for manipulating and executing the query.
type PreparedStatement interface {
	// UnpreparedStatement returns the original SQL statement before preparation.
	UnpreparedStatement() string
	// Revised returns the parsed query with positional parameters.
	Revised() string
	// NamedParameterPositions returns the parameter positions for the SQL statement.
	NamedParameterPositions() *NamedParameterPositions
	// BoundNamedParameterValues returns the bound named parameter values.
	BoundNamedParameterValues() BoundNamedParameterValues
	// BindNamedParameterValue binds a value to a named parameter in the SQL statement.
	BindNamedParameterValue(bindParameter string, bindValue any) error
	// BindNamedParameterValues sets the values for multiple named parameters using a list of BindNamedParameterValueFunc functions.
	BindNamedParameterValues(binderFuncs ...BindNamedParameterValueFunc) error
	// Exec executes the prepared SQL statement with the bound parameters.
	Exec(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (sql.Result, error)
	// ExecContext executes the prepared SQL statement with the bound parameters in the provided context.
	ExecContext(ctx context.Context, sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (sql.Result, error)
	// Query executes the prepared SQL statement as a query with the bound parameters.
	Query(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Rows, error)
	// QueryContext executes the prepared SQL statement as a query with the bound parameters in the provided context.
	QueryContext(ctx context.Context, sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Rows, error)
	// QueryRow executes the prepared SQL statement as a single-row query with the bound parameters.
	QueryRow(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Row, error)
	// QueryRowContext executes the prepared SQL statement as a single-row query with the bound parameters in the provided context.
	QueryRowContext(ctx context.Context, sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Row, error)
}

// preparedStatement is a struct that handles the translation of named parameters to positional parameters for SQL statements.
type preparedStatement struct {
	boundNamedParamValues BoundNamedParameterValues
	namedParamPositions   *NamedParameterPositions
	revisedStatement      string
	originalStatement     string
}

// getTotalIndices returns the total number of parameter positions in the statement.
func (p preparedStatement) getTotalIndices() int {
	if internal.IsNilOrZeroValue(p.namedParamPositions) {
		return 0
	}
	return p.namedParamPositions.totalPositions
}

// resetNamedParametersValues resets the boundNamedParamValues field to an empty BoundNamedParameterValues slice.
func (p *preparedStatement) resetNamedParametersValues() {
	if count := p.getTotalIndices(); count > 0 {
		p.boundNamedParamValues = make(BoundNamedParameterValues, count)
	}
}

// UnpreparedStatement returns the original SQL statement before preparation.
func (p preparedStatement) UnpreparedStatement() string {
	return string(p.originalStatement)
}

// Revised returns the parsed query with positional parameters.
func (p preparedStatement) Revised() string {
	return p.revisedStatement
}

// BoundNamedParameterValues returns the bound named parameter values.
func (p preparedStatement) BoundNamedParameterValues() BoundNamedParameterValues {
	if len(p.boundNamedParamValues) < 1 {
		return nil
	}

	return p.boundNamedParamValues
}

// NamedParameterPositions returns the named parameter positions for the SQL statement.
func (p preparedStatement) NamedParameterPositions() *NamedParameterPositions {
	if internal.IsNilOrZeroValue(p.namedParamPositions) {
		return nil
	}

	if len(p.namedParamPositions.parameterPositions) < 1 {
		return nil
	}

	return p.namedParamPositions
}

// BindNamedParameterValue binds a value to a named parameter in the SQL statement.
func (p *preparedStatement) BindNamedParameterValue(parameterName string, bindValue any) error {
	if internal.IsNilOrZeroValue(p.namedParamPositions) {
		return nil
	}

	if positions := p.namedParamPositions.getPositions(parameterName); len(positions) > 0 {
		for i := range positions {
			p.boundNamedParamValues[positions[i]] = bindValue
		}
	}

	return nil
}

// BindNamedParameterValueFunc is a function that sets the value for a named parameter in the query.
type BindNamedParameterValueFunc func(p PreparedStatement) error

// BindNamedParameterValue creates a BindNamedParameterValueFunc that sets the value for a named parameter.
func BindNamedParameterValue(bindParameter string, bindValue any) BindNamedParameterValueFunc {
	return func(p PreparedStatement) error {
		return p.BindNamedParameterValue(bindParameter, bindValue)
	}
}

// BindNamedParameterValues sets the values for multiple named parameters using a list of BindNamedParameterValueFunc functions.
func (p *preparedStatement) BindNamedParameterValues(binderFuncs ...BindNamedParameterValueFunc) error {
	for i := range binderFuncs {
		if internal.IsNilOrZeroValue(binderFuncs[i]) {
			continue
		}
		if err := binderFuncs[i](p); err != nil {
			return err
		}
	}

	return nil
}

func bindValuesGiven(binderFuncs []BindNamedParameterValueFunc) bool {
	return !internal.IsNilOrZeroValue(binderFuncs)
}

// prepare prepares the SQL statement with the provided context, SQL executor, and binder functions.
func (p *preparedStatement) prepare(ctx context.Context, sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Stmt, error) {
	if internal.IsNilOrZeroValue(sqlExecutor) {
		return nil, errors.New("sql executor (db) is nil")
	}

	if internal.IsNil(ctx) {
		ctx = context.Background()
	}

	prepared, err := sqlExecutor.PrepareContext(ctx, p.Revised())
	if err != nil {
		return nil, err
	}

	if bindValuesGiven(binderFuncs) {
		if err := p.BindNamedParameterValues(binderFuncs...); err != nil {
			return nil, err
		}
	}

	return prepared, nil
}

// Exec executes the prepared SQL statement with the bound parameters.
func (p *preparedStatement) Exec(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (sql.Result, error) {
	return p.ExecContext(context.Background(), sqlExecutor, binderFuncs...)
}

// ExecContext executes the prepared SQL statement with the bound parameters in the provided context.
func (p *preparedStatement) ExecContext(ctx context.Context, sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (sql.Result, error) {
	defer func() {
		p.resetNamedParametersValues()
	}()

	prepStmnt, err := p.prepare(ctx, sqlExecutor, binderFuncs...)
	if err != nil {
		return nil, err
	}

	return prepStmnt.Exec(p.boundNamedParamValues...)
}

// Query executes the prepared SQL statement as a query with the bound parameters.
func (p *preparedStatement) Query(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Rows, error) {
	return p.QueryContext(context.Background(), sqlExecutor, binderFuncs...)
}

// QueryContext executes the prepared SQL statement as a query with the bound parameters in the provided context.
func (p *preparedStatement) QueryContext(ctx context.Context, sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Rows, error) {
	defer func() {
		p.resetNamedParametersValues()
	}()

	prepStmnt, err := p.prepare(ctx, sqlExecutor, binderFuncs...)
	if err != nil {
		return nil, err
	}

	return prepStmnt.Query(p.boundNamedParamValues...)
}

// QueryRow executes the prepared SQL statement as a single-row query with the bound parameters.
func (p *preparedStatement) QueryRow(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Row, error) {
	return p.QueryRowContext(context.Background(), sqlExecutor, binderFuncs...)
}

// QueryRowContext executes the prepared SQL statement as a single-row query with the bound parameters in the provided context.
func (p *preparedStatement) QueryRowContext(ctx context.Context, sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Row, error) {
	defer func() {
		p.resetNamedParametersValues()
	}()

	prepStmnt, err := p.prepare(ctx, sqlExecutor, binderFuncs...)
	if err != nil {
		return nil, err
	}

	return prepStmnt.QueryRow(p.boundNamedParamValues...), nil
}

var _ PreparedStatement = (*preparedStatement)(nil)
