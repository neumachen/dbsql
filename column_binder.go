package dbsql

import (
	"fmt"

	"github.com/neumachen/dbsql/internal"
)

// Column represents a database column name.
//
// It is a type alias for string, providing type safety and clarity
// when working with column names in database operations.
type Column string

// String returns the Column as a string.
//
// This method allows Column to satisfy the Stringer interface,
// making it easy to use in string contexts.
//
// Returns:
//   - string: The column name as a string.
//
// Example:
//
//	col := Column("user_id")
//	fmt.Println(col.String()) // Output: user_id
func (c Column) String() string {
	return string(c)
}

// Columns is a slice of Column values.
//
// It represents a collection of database columns, typically used
// to describe the structure of a database table or query result.
type Columns []Column

// Count returns the number of columns in the Columns slice.
//
// Returns:
//   - int: The number of columns.
//
// Example:
//
//	cols := Columns{Column("id"), Column("name"), Column("email")}
//	fmt.Println(cols.Count()) // Output: 3
func (c Columns) Count() int {
	return len(c)
}

// HasColumn checks if the Columns slice contains the given column.
//
// Parameters:
//   - column: The Column to search for.
//
// Returns:
//   - bool: true if the column is found, false otherwise.
//
// Example:
//
//	cols := Columns{Column("id"), Column("name"), Column("email")}
//	fmt.Println(cols.HasColumn(Column("name")))  // Output: true
//	fmt.Println(cols.HasColumn(Column("phone"))) // Output: false
func (c Columns) HasColumn(column Column) bool {
	found := false
	for i := range c {
		if found = c[i] == column; found {
			return found
		}
	}
	return found
}

// ColumnBinderFunc is a function type that defines how to bind a column value from a mapped row to a struct field.
//
// A ColumnBinderFunc takes two parameters:
//   - column: The Column object representing the database column.
//   - mappedRow: The MappedRow object containing the row data.
//
// It returns an error if the binding process fails.
//
// This function type is typically used in row mapping operations to customize how
// column values are bound to struct fields. It allows for type-specific binding logic,
// data transformations, and custom error handling during the binding process.
//
// Parameters:
//   - column: A Column object representing the database column being bound.
//   - mappedRow: A MappedRow object containing the data for the current row.
//
// Returns:
//   - error: An error if the binding process fails, or nil if successful.
//
// Example:
//
//	type User struct {
//	    ID   int
//	    Name string
//	}
//
//	// Define a custom ColumnBinderFunc for the "id" column
//	idBinder := func(column Column, mappedRow MappedRow) error {
//	    value, found := mappedRow.Get(column)
//	    if !found {
//	        return fmt.Errorf("column %s not found", column)
//	    }
//	    id, ok := value.(int)
//	    if !ok {
//	        return fmt.Errorf("expected int for column %s, got %T", column, value)
//	    }
//	    user.ID = id
//	    return nil
//	}
//
//	// Use the custom binder
//	binderMap := ColumnBinderMap{
//	    Column("id"): idBinder,
//	}
//
// Note: The actual implementation of MappedRow is assumed to be available in the package.
type ColumnBinderFunc func(column Column, mappedRow MappedRow) error

// ColumnBinderMap is a map of Column to ColumnBinderFunc, used to map row data to a struct.
//
// It provides a flexible way to define custom binding logic for each column
// when mapping database rows to Go structs.
type ColumnBinderMap map[Column]ColumnBinderFunc

// Count returns the number of columns in the ColumnBinderMap.
//
// Returns:
//   - int: The number of columns in the map.
//
// Example:
//
//	binderMap := ColumnBinderMap{
//	    Column("id"):    idBinder,
//	    Column("name"):  nameBinder,
//	    Column("email"): emailBinder,
//	}
//	fmt.Println(binderMap.Count()) // Output: 3
func (c ColumnBinderMap) Count() int {
	return len(c)
}

// Columns returns a slice of Column representing the columns (field names) in the ColumnBinderMap.
//
// The order of the columns is determined by the order in which they were added to the map.
//
// Returns:
//   - Columns: A slice of Column objects representing the columns in the map.
//
// Example:
//
//	binderMap := ColumnBinderMap{
//	    Column("id"):    idBinder,
//	    Column("name"):  nameBinder,
//	    Column("email"): emailBinder,
//	}
//	cols := binderMap.Columns()
//	for _, col := range cols {
//	    fmt.Println(col)
//	}
//	// Output (order may vary):
//	// id
//	// name
//	// email
func (c ColumnBinderMap) Columns() Columns {
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

// BindColumnToField creates a ColumnBinderFunc that binds a column value to a field of type T.
//
// This function is a generic helper that simplifies the process of binding column values
// to struct fields, especially when working with custom types or when additional processing
// is needed during the binding.
//
// Type Parameters:
//   - T: The type of the value to be bound. This should match the type of the column value.
//
// Parameters:
//   - bindFunc: A function that takes a value of type T and returns an error.
//     This function is responsible for setting the value on the target struct field.
//
// Returns:
//   - ColumnBinderFunc: A function that conforms to the ColumnBinderFunc type and can be
//     used in row mapping operations.
//
// The returned ColumnBinderFunc does the following:
//  1. Retrieves the value for the given column from the MappedRow.
//  2. If the value is nil, zero, or not found, it returns nil (no error).
//  3. Attempts to assert the value to type T.
//  4. If the type assertion fails, it returns an error with details about the mismatch.
//  5. If successful, it calls the provided bindFunc with the typed value.
//
// Example usage:
//
//	type MyStruct struct {
//	  ID int
//	}
//
//	binder := BindColumnToField(func(value int) error {
//	  s.ID = value
//	  return nil
//	})
//
// Error Handling:
//   - Returns an error if the type assertion fails, providing details about the
//     expected and actual types.
//   - Propagates any error returned by the bindFunc.
//
// Note: This function uses generics and requires Go 1.18 or later.
func BindColumnToField[T any](bindFunc func(value T) error) ColumnBinderFunc {
	return func(column Column, mappedRow MappedRow) error {
		value, found := mappedRow.Get(column)
		if internal.IsNilOrZeroValue(value) || !found {
			return nil
		}
		typedValue, ok := value.(T)
		if !ok {
			// return error type assertion failed for the given
			// value using T.
			return fmt.Errorf(
				"column %s has a type of %T and does not match asserted type: %T",
				column.String(),
				value,
				*new(T),
			)
		}

		return bindFunc(typedValue)
	}
}
