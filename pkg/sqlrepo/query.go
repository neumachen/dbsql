package sqlrepo

import (
	"context"

	"github.com/neumachen/dbsql"
)

// QueryContext executes a query that returns rows, typically a SELECT.
// The binderFuncs are for binding parameter values to the prepared statement.
//
// It uses generics to allow for type-safe row binding. The generic type T must
// implement the RowBinder interface.
//
// This function handles the full lifecycle of the query execution:
// 1. Binding parameters (if any)
// 2. Executing the query
// 3. Mapping the rows to a slice of type T
// 4. Proper resource cleanup
//
// The implementation uses deferred closure to ensure that rows are always closed,
// preventing resource leaks.
//
// The generic approach allows for type-safe queries without the need for
// reflection, improving performance and compile-time type checking.
func QueryContext[T RowBinder](
	ctx context.Context,
	dbPrepExec dbsql.DB,
	preparedStatement dbsql.PreparedStatement,
	binderFuncs dbsql.BindParameterValueFuncs,
) ([]T, error) {
	if len(binderFuncs) > 0 {
		if err := preparedStatement.BindParameterValues(binderFuncs...); err != nil {
			return nil, err
		}
	}

	rows, err := dbsql.QueryContext(
		ctx,
		dbPrepExec,
		preparedStatement,
		binderFuncs...,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	mappedRows, err := dbsql.MapRows(rows)
	if err != nil {
		return nil, err
	}

	return BindMappedRows[T](mappedRows)
}
