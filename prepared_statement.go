package dbsql

import (
	"github.com/neumachen/dbsql/internal"
)

// ParameterPositions is a struct that represents the positions of parameters in an SQL statement.
// The parameterPositions field is a map that stores the positions of parameters, where the key is the parameter name
// and the value is a slice of integers representing the positions.
// The totalPositions field is an integer representing the total number of parameter positions.
type ParameterPositions struct {
	parameterPositions map[string][]int
	totalPositions     int
}

// getPositions is a method of the NamedParameterPositions struct.
// It takes a parameter name as input and returns a slice of integers representing the positions of the parameter.
func (p ParameterPositions) getPositions(parameter string) []int {
	if internal.IsNilOrZeroValue(p.parameterPositions) {
		return nil
	}

	v, ok := p.parameterPositions[parameter]
	if !ok {
		return nil
	}

	return v
}

// insert is a method of the NamedParameterPositions struct.
// It takes a parameter name and a position as input and inserts the position into the parameterPositions map.
func (p *ParameterPositions) insert(parameter string, position int) {
	if internal.IsNilOrZeroValue(p.parameterPositions) {
		p.parameterPositions = make(map[string][]int)
	}

	v, ok := p.parameterPositions[parameter]
	if !ok {
		v = make([]int, 0, 1)
	}

	p.parameterPositions[parameter] = append(v, position)
	p.totalPositions++
}

// BoundParameterValues is a type alias for a slice of any (an empty interface).
// It represents the positional parameters in an SQL statement.
type BoundParameterValues []any

// PreparedStatement is an interface that represents an SQL query with named parameters.
// It defines several methods for manipulating and executing the query.
type PreparedStatement interface {
	// UnpreparedStatement returns the original SQL statement before preparation.
	UnpreparedStatement() string
	// Revised returns the parsed query with positional parameters.
	Revised() string
	ResetParametersValues()
	// ParameterPositions returns the parameter positions for the SQL statement.
	ParameterPositions() *ParameterPositions
	// BoundNamedParameterValues returns the bound named parameter values.
	BoundParameterValues() BoundParameterValues
	// BindParameterValue binds a value to a named parameter in the SQL statement.
	BindParameterValue(bindParameter string, bindValue any) error
	// BindParameterValues sets the values for multiple named parameters using a list of BindNamedParameterValueFunc functions.
	BindParameterValues(binderFuncs ...BindParameterValueFunc) error
}

// preparedStatement is a struct that handles the translation of named parameters to positional parameters for SQL statements.
type preparedStatement struct {
	boundNamedParamValues BoundParameterValues
	namedParamPositions   *ParameterPositions
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

// ResetParametersValues resets the boundNamedParamValues field to an empty BoundParameterValues slice.
func (p *preparedStatement) ResetParametersValues() {
	if count := p.getTotalIndices(); count > 0 {
		p.boundNamedParamValues = make(BoundParameterValues, count)
	}
}

// UnpreparedStatement returns the original SQL statement before preparation.
func (p preparedStatement) UnpreparedStatement() string {
	return p.originalStatement
}

// Revised returns the parsed query with positional parameters.
func (p preparedStatement) Revised() string {
	return p.revisedStatement
}

// BoundNamedParameterValues returns the bound named parameter values.
func (p preparedStatement) BoundParameterValues() BoundParameterValues {
	if len(p.boundNamedParamValues) < 1 {
		return nil
	}

	return p.boundNamedParamValues
}

// ParameterPositions returns the named parameter positions for the SQL statement.
func (p preparedStatement) ParameterPositions() *ParameterPositions {
	if internal.IsNilOrZeroValue(p.namedParamPositions) {
		return nil
	}

	if len(p.namedParamPositions.parameterPositions) < 1 {
		return nil
	}

	return p.namedParamPositions
}

// BindParameterValue binds a value to a named parameter in the SQL statement.
func (p *preparedStatement) BindParameterValue(parameterName string, bindValue any) error {
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

// BindParameterValueFunc is a function that sets the value for a named parameter in the query.
type BindParameterValueFunc func(p PreparedStatement) error

// BindParameterValueFuncs ...
type BindParameterValueFuncs []BindParameterValueFunc

// NewBindParameterValueFuncs ...
func NewBindParameterValueFuncs(bindParameterValueFuncs ...BindParameterValueFunc) BindParameterValueFuncs {
	count := len(bindParameterValueFuncs)
	if count == 0 {
		return nil
	}
	newBindParameterValueFuncs := make(BindParameterValueFuncs, count)
	copy(newBindParameterValueFuncs, bindParameterValueFuncs)
	return newBindParameterValueFuncs
}

// BindParameterValue creates a BindNamedParameterValueFunc that sets the value for a named parameter.
func BindParameterValue(bindParameter string, bindValue any) BindParameterValueFunc {
	return func(p PreparedStatement) error {
		return p.BindParameterValue(bindParameter, bindValue)
	}
}

// BindParameterValues sets the values for multiple named parameters using a list of BindNamedParameterValueFunc functions.
func (p *preparedStatement) BindParameterValues(binderFuncs ...BindParameterValueFunc) error {
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

var _ PreparedStatement = (*preparedStatement)(nil)
