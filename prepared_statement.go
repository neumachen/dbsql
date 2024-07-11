package sqlstmt

import (
	"database/sql"
	"errors"

	"github.com/neumachen/sqlstmt/internal"
)

// NamedParameterPositions is a struct that represents the positions of parameters in an SQL statement.
type NamedParameterPositions struct {
	//	parameterPositions A map that stores the positions of parameters, where the key is the parameter name and the value is a slice of integers representing the positions.
	parameterPositions map[string][]int
	//	totalPositions An integer representing the total number of parameter positions.
	totalPositions int
}

// getPositions is a method of the ParameterPositions struct.
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

// insert is a method of the ParameterPositions struct.
// It takes a parameter name and a position as input and inserts the position into the paramPositions map.
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
	// UpreparedStatement returns the original SQL statement before preparation.
	UnpreparedStatement() string
	// It returns the parsed query with positional parameters as a byte slice.
	Revised() string
	NamedParameterPositions() *NamedParameterPositions
	BoundNamedParameterValues() BoundNamedParameterValues
	BindNamedParameterValue(bindParameter string, bindValue any) error
	BindNamedParameterValues(binderFuncs ...BindNamedParameterValueFunc) error
	Exec(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (sql.Result, error)
	Query(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Rows, error)
	QueryRow(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Row, error)
}

// preparedStatement is a struct that handles the translation of named parameters to positional parameters for SQL statements.
type preparedStatement struct {
	boundNamedParamValues BoundNamedParameterValues
	namedParamPositions   *NamedParameterPositions
	revisedStatement      string
	originalStatement     string
}

// getTotalIndices it returns the total number of parameter positions in the statement.
func (p preparedStatement) getTotalIndices() int {
	if internal.IsNilOrZeroValue(p.namedParamPositions) {
		return 0
	}
	return p.namedParamPositions.totalPositions
}

// resetNamedParametersValues it resets the parameters field to an empty BoundNamedParameterValues slice.
func (p *preparedStatement) resetNamedParametersValues() {
	if count := p.getTotalIndices(); count > 0 {
		p.boundNamedParamValues = make(BoundNamedParameterValues, count)
	}
}

// setNamedParameterPosition sets the position for a named parameter in the parameterPositions field.
func (p *preparedStatement) setNamedParameterPosition(parameter string, position int) {
	if internal.IsNilOrZeroValue(p.namedParamPositions) {
		p.namedParamPositions = &NamedParameterPositions{}
	}

	p.namedParamPositions.insert(parameter, position)
}

func (p preparedStatement) UnpreparedStatement() string {
	return string(p.originalStatement)
}

// PreparedStatement returns the parsed query with positional parameters.
func (p preparedStatement) Revised() string {
	return p.revisedStatement
}

// BoundNamedParameterValues ...
func (p preparedStatement) BoundNamedParameterValues() BoundNamedParameterValues {
	if len(p.boundNamedParamValues) < 1 {
		return nil
	}

	return p.boundNamedParamValues
}

// NamedParameterPositions returns the parameter positions for the SQL statement.
func (p preparedStatement) NamedParameterPositions() *NamedParameterPositions {
	if internal.IsNilOrZeroValue(p.namedParamPositions) {
		return nil
	}

	if len(p.namedParamPositions.parameterPositions) < 1 {
		return nil
	}

	return p.namedParamPositions
}

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

// BindNamedParameterValue creates a SetParameterFunc that sets the value for a named parameter.
func BindNamedParameterValue(bindParameter string, bindValue any) BindNamedParameterValueFunc {
	return func(p PreparedStatement) error {
		return p.BindNamedParameterValue(bindParameter, bindValue)
	}
}

// BindNamedParameterValues sets the values for multiple named parameters using a list of SetParameterFunc functions.
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

func (p *preparedStatement) prepare(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Stmt, error) {
	if internal.IsNilOrZeroValue(sqlExecutor) {
		return nil, errors.New("sql executor (db) is nil")
	}

	prepared, err := sqlExecutor.Prepare(p.Revised())
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

func (p *preparedStatement) Exec(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (sql.Result, error) {
	defer func() {
		p.resetNamedParametersValues()
	}()

	prepStmnt, err := p.prepare(sqlExecutor, binderFuncs...)
	if err != nil {
		return nil, err
	}

	return prepStmnt.Exec(p.boundNamedParamValues...)
}

func (p *preparedStatement) Query(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Rows, error) {
	defer func() {
		p.resetNamedParametersValues()
	}()

	prepStmnt, err := p.prepare(sqlExecutor, binderFuncs...)
	if err != nil {
		return nil, err
	}

	return prepStmnt.Query(p.boundNamedParamValues...)
}

func (p *preparedStatement) QueryRow(sqlExecutor SQLExecutor, binderFuncs ...BindNamedParameterValueFunc) (*sql.Row, error) {
	defer func() {
		p.resetNamedParametersValues()
	}()

	prepStmnt, err := p.prepare(sqlExecutor, binderFuncs...)
	if err != nil {
		return nil, err
	}

	return prepStmnt.QueryRow(p.boundNamedParamValues...), nil
}

var _ PreparedStatement = (*preparedStatement)(nil)
