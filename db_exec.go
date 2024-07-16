package dbsql

import (
	"context"
	"database/sql"
)

// Exec executes the prepared SQL statement with the bound parameters.
// It uses the default background context.
//
// Parameters:
//   - dbPrepExec: An interface that can prepare and execute SQL statements.
//   - preparedStatement: The prepared statement to be executed.
//   - binderFuncs: Optional functions to bind parameter values.
//
// Returns:
//   - sql.Result: The result of the SQL execution.
//   - error: An error if the execution fails.
func Exec(
	dbPrepExec DBPreparerExecutor,
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
//
// Parameters:
//   - ctx: The context for the execution.
//   - dbPrepExec: An interface that can prepare and execute SQL statements.
//   - preparedStatement: The prepared statement to be executed.
//   - binderFuncs: Optional functions to bind parameter values.
//
// Returns:
//   - sql.Result: The result of the SQL execution.
//   - error: An error if the execution fails.
func ExecContext(
	ctx context.Context,
	dbPrepExec DBPreparerExecutor,
	preparedStatement PreparedStatement,
	binderFuncs ...BindParameterValueFunc,
) (
	sql.Result,
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

	return prepStmnt.ExecContext(
		ctx,
		preparedStatement.BoundParameterValues()...,
	)
}
