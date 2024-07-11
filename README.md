# SQLStmt

SQLStmt is a Go package that simplifies working with SQL statements containing named parameters. It provides a mechanism to convert SQL statements with named parameters to positional parameters for use with the `database/sql` package.

## Features

- [x] Converts SQL queries with named parameters to positional parameters.
- [x] Provides methods to set parameter values for the converted queries.
- [x] Seamless integration with the `sql/db` package for database interaction.
  - [x] Extend query functions
    - [x] [sql#DB.Exec](https://pkg.go.dev/database/sql#DB.Exec)
    - [x] [sql#DB.Query](https://pkg.go.dev/database/sql#DB.Query)
    - [x] [sql#DB.QueryRow](https://pkg.go.dev/database/sql#DB.QueryRow)
  - [ ] Support transaction methods.
  - [ ] Extend tests for the database helper functions.

## Installation

To use SQLStmt in your Go project, you need to have Go installed and set up. Then, you can use the following command to install the package:

```shell
go get github.com/neumachen/sqldb@latest
```

## NOTES

- Database functions reset the parameter values passed as argumetns OR using statement.SetParameters.
