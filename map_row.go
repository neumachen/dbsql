package dbsql

import (
	"database/sql"
)

// MappedRow represents a row of mapped column-value pairs.
type MappedRow map[Column]any

// Count returns the number of columns in the MappedRow.
func (m MappedRow) Count() int {
	return len(m)
}

// MapRow maps the columns and values of the given sql.Row to a MappedRow.
func MapRow(row *sql.Row, columns []string) (MappedRow, error) {
	values := make([]any, len(columns))
	for i := range values {
		values[i] = new(any)
	}

	err := row.Scan(values...)
	if err != nil {
		return nil, err
	}

	rowData := make(MappedRow)
	for i, col := range columns {
		rowData[Column(col)] = *(values[i].(*any))
	}

	return rowData, nil
}
