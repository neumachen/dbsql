package dbsql

import (
	"context"
	"database/sql"
)

// Query executes the prepared SQL statement as a query with the bound parameters.
func Query(
	dbPrepExec DbPreparerExecutor,
	preparedStatement PreparedStatement,
	binderFuncs ...BindParameterValueFunc,
) (
	*sql.Rows,
	error,
) {
	return QueryContext(
		context.Background(),
		dbPrepExec,
		preparedStatement,
		binderFuncs...,
	)
}

// QueryContext executes the prepared SQL statement as a query with the bound parameters in the provided context.
func QueryContext(
	ctx context.Context,
	dbPrepExec DbPreparerExecutor,
	preparedStatement PreparedStatement,
	binderFuncs ...BindParameterValueFunc,
) (
	*sql.Rows,
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

	return prepStmnt.QueryContext(
		ctx,
		preparedStatement.BoundParameterValues()...,
	)
}
