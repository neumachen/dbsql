package dbsql

import (
	"context"
	"database/sql"
)

// Exec executes the prepared SQL statement with the bound parameters.
func Exec(
	dbPrepExec DbPreparerExecutor,
	preparedStatement PreparedStatement,
	binderFuncs ...BindParameterValueFunc,
) (
	sql.Result,
	error,
) {
	return ExecContext(
		context.Background(),
		dbPrepExec,
		preparedStatement,
		binderFuncs...,
	)
}

// ExecContext executes the prepared SQL statement with the bound parameters in the provided context.
func ExecContext(
	ctx context.Context,
	dbPrepExec DbPreparerExecutor,
	preparedStatement PreparedStatement,
	binderFuncs ...BindParameterValueFunc,
) (
	sql.Result,
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

	return prepStmnt.ExecContext(
		ctx,
		preparedStatement.BoundParameterValues()...,
	)
}
