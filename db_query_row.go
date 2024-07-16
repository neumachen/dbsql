package dbsql

import (
	"context"
	"database/sql"
)

// QueryRow executes the prepared SQL statement as a query with the bound parameters.
func QueryRow(
	dbPrepExec DBPreparerExecutor,
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
	dbPrepExec DBPreparerExecutor,
	preparedStatement PreparedStatement,
	binderFuncs ...BindParameterValueFunc,
) (
	*sql.Row,
	error,
) {
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
