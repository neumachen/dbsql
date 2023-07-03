package sqlstmt

import (
	"database/sql"
)

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
		rowData[col] = *(values[i].(*interface{}))
	}

	return rowData, nil
}
