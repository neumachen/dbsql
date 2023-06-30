# SQLStmt

SQLStmt is a Go package that simplifies working with SQL queries containing named parameters. It provides a mechanism to convert SQL statements with named parameters to positional parameters for use with the `lib/pq` package, which is a popular PostgreSQL driver for Go.

## Features

- [ ] Converts SQL queries with named parameters to positional parameters.
- [ ] Provides methods to set parameter values for the converted queries.
- [ ] Seamless integration with the `sql/db` package for database interaction.
  - [ ] Extend query functions
    - [ ] [sql#DB.Exec](https://pkg.go.dev/database/sql#DB.Exec)
    - [ ] [sql#DB.Query](https://pkg.go.dev/database/sql#DB.Query)
    - [ ] [sql#DB.QueryRow](https://pkg.go.dev/database/sql#DB.QueryRow)
  - [ ] Support transaction methods.
  - [ ] Extend tests for the database helper functions.

## Installation

To use SQLStmt in your Go project, you need to have Go installed and set up. Then, you can use the following command to install the package:

```shell
go get github.com/neumachen/sqlstmt@latest
```

## NOTES

- Database functions reset the parameter values passed as argumetns OR using statement.SetParameters.
