package dbsql

import (
	"context"
	"database/sql"
)

// QueryRow executes the prepared SQL statement as a query with the bound parameters.
func QueryRow(
	dbPrepExec DbPreparerExecutor,
	preparedStatement PreparedStatement,
	binderFuncs ...BindParameterValueFunc,
) (
	*sql.Row,
	error,
) {
	return QueryRowContext(
		context.Background(),
		dbPrepExec,
		preparedStatement,
		binderFuncs...,
	)
}

// QueryRowContext executes the prepared SQL statement as a query with the bound parameters in the provided context.
func QueryRowContext(
	ctx context.Context,
	dbPrepExec DbPreparerExecutor,
	preparedStatement PreparedStatement,
	binderFuncs ...BindParameterValueFunc,
) (
	*sql.Row,
	error,
) {
	// TODO: Should this be resetting the bound values for the given prepared statement, or should it be left to the
	// caller? An argument can be made that the internal prepare method sets the bound values in this block. But does
	// that also mean that when a gdiven prepared statement that already has bound values will be reset?
	defer func() {
		preparedStatement.ResetParametersValues()
	}()

	prepStmnt, err := dbPrepare(
		ctx,
		dbPrepExec,
		preparedStatement,
		binderFuncs...,
	)
	if err != nil {
		return nil, err
	}

	return prepStmnt.QueryRowContext(
		ctx,
		preparedStatement.BoundParameterValues()...,
	), nil
}
