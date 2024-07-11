package sqldb

import (
	"fmt"

	"github.com/neumachen/sqldb/internal"
)

// Column is a string type representing a column name.
type Column string

// ToString returns the Column as a string.
func (c Column) ToString() string {
	return string(c)
}

// Columns is a slice of Column values.
type Columns []Column

// Count returns the number of columns.
func (c Columns) Count() int {
	return len(c)
}

// HasColumn checks if the Columns slice contains the given column.
func (c Columns) HasColumn(column string) bool {
	for i := range c {
		if c[i].ToString() == column {
			return true
		}
	}
	return false
}

// Columns returns a slice of strings representing the columns (field names) in the ColumnMapperMap.
// The order of the columns is determined by the order in which they were added to the map.
func (c ColumnMapperMap) Columns() Columns {
	count := len(c)
	if count < 1 {
		return nil
	}
	columns := make(Columns, count)
	for k := range c {
		columns[count-1] = k
		count--
	}
	return columns
}

// ColumnMapperFunc is a function that maps a column value to a field in a struct.
type ColumnMapperFunc func(column Column, row MappedRow) error

// ColumnMapperMap is a map of Column to ColumnMapperFunc, used to map row data to a struct.
type ColumnMapperMap map[Column]ColumnMapperFunc

// ColumnMapper is a generic type that represents a function responsible for setting a column value.
type ColumnMapper[T any] struct {
	// MapperFunc is a function that takes a value of type T and returns an error, if any.
	MapperFunc func(value T) error
}

// MapColumn is a higher-order function that takes a ColumnMapper and returns a ColumnMapperFunc.
// The returned ColumnMapperFunc can be used to set a column value in a row of type MappedRow,
// using the provided column and row values.
func MapColumn[T any](mapFunc func(value T) error) ColumnMapperFunc {
	return func(column Column, row MappedRow) error {
		value, ok := row[column]
		if internal.IsNilOrZeroValue(value) || !ok {
			return nil
		}
		typedValue, ok := value.(T)
		if !ok {
			// return error type assertion failed for the given
			// value using T.
			return fmt.Errorf(
				"column %s has a type of %T and does not match asserted type: %T",
				column.ToString(),
				value,
				*new(T),
			)
		}

		return mapFunc(typedValue)
	}
}
