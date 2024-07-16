package dbsql

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
