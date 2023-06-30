package sqlstmt

import (
	"database/sql"
	"errors"
)

// ParameterPositions is a struct that represents the positions of parameters in an SQL statement.
type ParameterPositions struct {
	//	paramPositions A map that stores the positions of parameters, where the key is the parameter name and the value is a slice of integers representing the positions.
	paramPositions map[string][]int
	//	totalPositions An integer representing the total number of parameter positions.
	totalPositions int
}

// getPositions is a method of the ParameterPositions struct.
// It takes a parameter name as input and returns a slice of integers representing the positions of the parameter.
func (p ParameterPositions) getPositions(parameter string) []int {
	if isNilOrZero(p.paramPositions) {
		return nil
	}

	v, ok := p.paramPositions[parameter]
	if !ok {
		return nil
	}

	return v
}

// insert is a method of the ParameterPositions struct.
// It takes a parameter name and a position as input and inserts the position into the paramPositions map.
func (p *ParameterPositions) insert(parameter string, position int) {
	if isNilOrZero(p.paramPositions) {
		p.paramPositions = make(map[string][]int)
	}

	v, ok := p.paramPositions[parameter]
	if !ok {
		v = make([]int, 0, 1)
	}

	p.paramPositions[parameter] = append(v, position)
	p.totalPositions++
}

// PositionalParameters is a type alias for a slice of any (an empty interface).
// It represents the positional parameters in an SQL statement.
type PositionalParameters []any

// Statement is an interface that represents an SQL query with named parameters.
// It defines several methods for manipulating and executing the query.
type Statement interface {
	// GetQueryWithPositionals is a method of the Statement interface implemented by the statement struct.
	// It returns the parsed query with positional parameters as a byte slice.
	GetQueryWithPositionals() []byte
	// GetParameterPositions returns the parsed positional parameters.
	GetParameterPositions() *ParameterPositions
	// GetPositionalParameters returns the positional parameters.
	GetPositionalParameters() PositionalParameters
	// SetParameter sets the value for a named parameter in the query.
	SetParameter(name string, value any) error
	// SetParameters sets the values for multiple named parameters using a list of SetParameterFunc functions.
	SetParameters(funcs ...SetParameterFunc) error
	// Exec it executes the SQL statement using the provided database connection (db) and returns the result.
	// It takes a variadic number of SetParameterFunc functions as input to set the parameter values.
	// It resets the parameters for the given Statement regardless of outcoe of the execution.
	Exec(db Database, setters ...SetParameterFunc) (sql.Result, error)
	// Query it executes the SQL statement as a query using the provided database connection (db) and returns the resulting rows.
	// It takes a variadic number of SetParameterFunc functions as input to set the parameter values.
	// It resets the parameters for the given Statement regardless of outcoe of the execution.
	Query(db Database, setters ...SetParameterFunc) (*sql.Rows, error)
	// QueryRow it executes the SQL statement as a query and returns a single row using the provided database connection (db).
	// It takes a variadic number of SetParameterFunc functions as input to set the parameter values.
	// It resets the parameters for the given Statement regardless of outcoe of the execution.
	QueryRow(db Database, setters ...SetParameterFunc) (*sql.Row, error)
}

// statement is a struct that handles the translation of named parameters to positional parameters for SQL statements.
type statement struct {
	// parameters A PositionalParameters slice representing the positional parameters.
	parameters PositionalParameters
	// parameterPositions A pointer to a ParameterPositions struct representing the parameter positions.
	parameterPositions *ParameterPositions
	// revisedQuery A byte slice representing the revised query with positional parameters.
	revisedQuery []byte
}

var _ Statement = (*statement)(nil)

// getTotalIndices it returns the total number of parameter positions in the statement.
func (s statement) getTotalIndices() int {
	if isNilOrZero(s.parameterPositions) {
		return 0
	}
	return s.parameterPositions.totalPositions
}

// resetParameters it resets the parameters field to an empty PositionalParameters slice.
func (s *statement) resetParameters() {
	if count := s.getTotalIndices(); count > 0 {
		s.parameters = make(PositionalParameters, count)
	}
}

// setPosition sets the position for a named parameter in the parameterPositions field.
func (s *statement) setPosition(parameter string, position int) {
	if isNilOrZero(s.parameterPositions) {
		s.parameterPositions = &ParameterPositions{}
	}

	s.parameterPositions.insert(parameter, position)
}

// GetQueryWithPositionals returns the parsed query with positional parameters.
func (s statement) GetQueryWithPositionals() []byte {
	return s.revisedQuery
}

// GetPositionalParameters ...
func (s statement) GetPositionalParameters() PositionalParameters {
	if len(s.parameters) < 1 {
		return nil
	}

	return s.parameters
}

// GetParameterPositions returns the parameter positions for the SQL statement.
func (s statement) GetParameterPositions() *ParameterPositions {
	if isNilOrZero(s.parameterPositions) {
		return nil
	}

	if len(s.parameterPositions.paramPositions) < 1 {
		return nil
	}

	return s.parameterPositions
}

// SetParameter sets the value for a named parameter in the query.
func (s *statement) SetParameter(name string, value any) error {
	if isNilOrZero(s.parameterPositions) {
		return nil
	}

	if positions := s.parameterPositions.getPositions(name); len(positions) > 0 {
		for i := range positions {
			s.parameters[positions[i]] = value
		}
	}

	return nil
}

// SetParameterFunc is a function that sets the value for a named parameter in the query.
type SetParameterFunc func(p Statement) error

// SetParameter creates a SetParameterFunc that sets the value for a named parameter.
func SetParameter(parameter string, value any) SetParameterFunc {
	return func(p Statement) error {
		return p.SetParameter(parameter, value)
	}
}

// SetParameters sets the values for multiple named parameters using a list of SetParameterFunc functions.
func (s *statement) SetParameters(funcs ...SetParameterFunc) error {
	for i := range funcs {
		if isNilOrZero(funcs[i]) {
			continue
		}
		if err := funcs[i](s); err != nil {
			return err
		}
	}

	return nil
}

func parametersGiven(setterParameters []SetParameterFunc) bool {
	return !isNilOrZero(setterParameters)
}

func (s *statement) prepare(db Database, setters ...SetParameterFunc) (*sql.Stmt, error) {
	if isNilOrZero(db) {
		return nil, errors.New("database is nil")
	}

	prepStmnt, err := db.Prepare(string(s.GetQueryWithPositionals()))
	if err != nil {
		return nil, err
	}

	if parametersGiven(setters) {
		if err := s.SetParameters(setters...); err != nil {
			return nil, err
		}
	}

	return prepStmnt, nil
}

func (s *statement) Exec(db Database, setters ...SetParameterFunc) (sql.Result, error) {
	defer func() {
		s.resetParameters()
	}()

	prepStmnt, err := s.prepare(db, setters...)
	if err != nil {
		return nil, err
	}

	return prepStmnt.Exec(s.parameters...)
}

func (s *statement) Query(db Database, setters ...SetParameterFunc) (*sql.Rows, error) {
	defer func() {
		s.resetParameters()
	}()

	prepStmnt, err := s.prepare(db, setters...)
	if err != nil {
		return nil, err
	}

	return prepStmnt.Query(s.parameters...)
}

func (s *statement) QueryRow(db Database, setters ...SetParameterFunc) (*sql.Row, error) {
	defer func() {
		s.resetParameters()
	}()

	prepStmnt, err := s.prepare(db, setters...)
	if err != nil {
		return nil, err
	}

	return prepStmnt.QueryRow(s.parameters...), nil
}
