package sqlstmt

import "fmt"

type (
	ColumnMapperFunc func(column Column, row MappedRow) error

	ColumnMapperMap map[Column]ColumnMapperFunc

	Column string

	Columns []Column
)

func (c Column) ToString() string {
	return string(c)
}

func (c Columns) Count() int {
	return len(c)
}

func (c Columns) HasColumn(column string) bool {
	for i := range c {
		if c[i].ToString() == column {
			return true
		}
	}
	return false
}

// Columns returns a slice of strings representing the columns (field names) in the ColumnMappers map.
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

// ColumnMapper is a generic type that represents a function responsible for setting a column value.
type ColumnMapper[T any] struct {
	// MapperFunc function takes a value of type T and returns an error, if any.
	MapperFunc func(value T) error
}

// MapColumn is a higher-order function that takes a ColumnMapper and returns a ColumnMapperFunc.
// The returned ColumnMapperFunc can be used to set a column value in a row of type MappedRow,
// using the provided column and row values.
func MapColumn[T any](mapFunc func(value T) error) ColumnMapperFunc {
	return func(column Column, row MappedRow) error {
		value, ok := row[column]
		if isNilOrZero(value) || !ok {
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
