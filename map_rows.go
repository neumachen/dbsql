package sqldb

import (
	"database/sql"
)

// MappedRows represents a collection of mapped rows.
type MappedRows []MappedRow

// Count ...
func (m MappedRows) Count() int {
	return len(m)
}

// MapRows maps each column and value of the given sql.Rows to a MappedRows structure.
func MapRows(rows *sql.Rows) (MappedRows, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	result := make(MappedRows, 0) // Initialize an empty result

	for rows.Next() {
		// Prepare the value pointers for scanning
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create a map to store the column-value pairs for the current row
		rowMap := make(MappedRow)

		// Map each column to its corresponding value
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				// Convert []byte to []byte to prevent underlying array sharing
				rowMap[Column(col)] = append([]byte(nil), b...)
			} else {
				rowMap[Column(col)] = val
			}
		}

		// Append the mapped row to the result
		result = append(result, rowMap)
	}

	return result, nil
}
