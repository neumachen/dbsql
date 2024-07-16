package dbsql

import (
	"context"
	"database/sql"
)

// Query executes the prepared SQL statement as a query with the bound parameters.
func Query(
	dbPrepExec DBPreparerExecutor,
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
	dbPrepExec DBPreparerExecutor,
	preparedStatement PreparedStatement,
	binderFuncs ...BindParameterValueFunc,
) (
	*sql.Rows,
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

	return prepStmnt.QueryContext(
		ctx,
		preparedStatement.BoundParameterValues()...,
	)
}
