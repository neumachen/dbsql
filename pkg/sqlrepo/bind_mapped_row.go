package sqlrepo

import (
	"errors"

	"github.com/neumachen/dbsql"
	"github.com/neumachen/dbsql/internal"
)

func BindMappedRow(
	columnBinders dbsql.ColumnBinders,
	mappedRow dbsql.MappedRow,
) error {
	if internal.IsNil(columnBinders) {
		return errors.New("column binders is nil")
	}

	for i := range columnBinders {
		columnBinder := columnBinders[i]
		if internal.IsNilOrZeroValue(columnBinder) {
			return errors.New("column binder is nil or has a value")
		}

		if internal.IsNilOrZeroValue(columnBinder.Column()) {
			return errors.New("column binder has a zero value for column field")
		}

		if err := columnBinder.BindColumn(mappedRow); err != nil {
			return err
		}
	}

	return nil
}

// RowBinder ...
type RowBinder interface {
	ColumnBinders() dbsql.ColumnBinders
}

// BindMappedRows binds the mapped row to the column binders.
func BindMappedRows[T RowBinder](
	mappedRows dbsql.MappedRows,
) (
	[]T,
	error,
) {
	if len(mappedRows) < 1 {
		return nil, nil
	}

	rowsBinder := make([]T, len(mappedRows))
	for i := range mappedRows {
		mappedRow := mappedRows[i]
		for j := range rowsBinder {
			if err := BindMappedRow(
				rowsBinder[j].ColumnBinders(),
				mappedRow,
			); err != nil {
				return nil, err
			}
		}
	}

	return rowsBinder, nil
}
