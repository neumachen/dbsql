# dbsql

The `dbsql` package provides a set of utilities and abstractions for working with SQL databases in Go. It aims to simplify the process of managing prepared statements, mapping query results to struct fields, and handling named parameters.


## Features

- **Prepared Statements**: The `PreparedStatement` interface and `preparedStatement` struct simplify the creation, execution, and management of prepared SQL statements with named parameters.
- **Column Mapping**: The `ColumnMapperFunc` and `ColumnMapper` types allow you to define custom mapping functions to map SQL query results to struct fields.
- **Row Mapping**: The `MapRow` and `MapRows` functions provide a convenient way to map SQL query results to `MappedRow` and `MappedRows` data structures.

## Installation

To use the `dbsql` package, you can install it using the following command:

```
go get github.com/your-username/dbsql
```

## Usage

### Prepared Statements

Here's an example of how to use the `PreparedStatement` interface:

```go
import "github.com/your-username/dbsql"

// Prepare a statement with named parameters
stmt, err := dbsql.PrepareStatement("SELECT * FROM users WHERE name = @name AND age > @age")
if err != nil {
    // Handle error
}

// Bind the named parameters
if err := stmt.BindNamedParameterValue("name", "John"); err != nil {
    // Handle error
}
if err := stmt.BindNamedParameterValue("age", 30); err != nil {
    // Handle error
}

// Execute the prepared statement
rows, err := stmt.Query(db)
if err != nil {
    // Handle error
}
```

### Column Mapping

To map SQL query results to struct fields, you can use the `ColumnMapperFunc` and `ColumnMapper` types:

```go
type User struct {
    ID        int64
    Name      string
    Age       int
}

// Define the column mapping functions
columnMappers := dbsql.ColumnMapperMap{
    "id":   dbsql.MapColumn[int64](func(v int64) error { u.ID = v; return nil }),
    "name": dbsql.MapColumn[string](func(v string) error { u.Name = v; return nil }),
    "age":  dbsql.MapColumn[int](func(v int) error { u.Age = v; return nil }),
}

// Use the column mappers to map the query result to a User struct
row, err := stmt.QueryRow(db)
user := User{}
if err := user.mapRow(row, columnMappers); err != nil {
    // Handle error
}
```

### Row Mapping

You can use the `MapRow` and `MapRows` functions to map SQL query results to `MappedRow` and `MappedRows` data structures:

```go
rows, err := stmt.Query(db)
if err != nil {
    // Handle error
}

mappedRows, err := dbsql.MapRows(rows)
if err != nil {
    // Handle error
}

for _, row := range mappedRows {
    // Access the column-value pairs in the mapped row
    name := row["name"].(string)
    age := row["age"].(int)
    // ...
}
```

## Testing

The `dbsql` package includes a comprehensive test suite to ensure the reliability of its features. You can run the tests using the following command:

```
go test -v ./...
```

The tests cover a range of scenarios, including error handling, data mapping, and prepared statement execution.

## Contributing

If you'd like to contribute to the `dbsql` package, please feel free to submit a pull request or open an issue. Contributions are welcome and appreciated!

## License

The `dbsql` package is licensed under the [MIT License](LICENSE).
