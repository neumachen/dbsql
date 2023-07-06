package sqlstmt

import (
	"database/sql"
)

// MappedRow represents a row of mapped column-value pairs.
type MappedRow map[Column]any

func (m MappedRow) Count() int {
	return len(m)
}

func MapRow(row *sql.Row, columns []string) (MappedRow, error) {
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err := row.Scan(values...)
	if err != nil {
		return nil, err
	}

	rowData := make(MappedRow)
	for i, col := range columns {
		rowData[Column(col)] = *(values[i].(*interface{}))
	}

	return rowData, nil
}
