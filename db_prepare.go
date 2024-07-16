package dbsql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/neumachen/dbsql/internal"
)

func dbPrepare(
	ctx context.Context,
	dbPrep DBPreparer,
	preparedStatement PreparedStatement,
	binderFuncs ...BindParameterValueFunc,
) (
	*sql.Stmt,
	error,
) {
	if internal.IsNil(dbPrep) {
		return nil, errors.New("db connection is nil")
	}

	if internal.IsNil(preparedStatement) {
		return nil, errors.New("prepared statement is nil")
	}

	ctx = internal.InitIfNilContext(ctx)

	prepared, err := dbPrep.PrepareContext(ctx, preparedStatement.Revised())
	if err != nil {
		return nil, err
	}

	if bindValuesGiven(binderFuncs) {
		if err := preparedStatement.BindParameterValues(binderFuncs...); err != nil {
			return nil, err
		}
	}

	return prepared, nil
}

func bindValuesGiven(binderFuncs []BindParameterValueFunc) bool {
	return !internal.IsNilOrZeroValue(binderFuncs)
}
