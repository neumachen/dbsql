package sqlstmt

import (
	"database/sql"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
)

func deleteRecords(
	t *testing.T,
	db Database,
	expectedRowsAffected int,
	query []byte,
	setterFuncs ...SetParameterFunc,
) {
	stmnt, err := ConvertNamedToPositionalParams(query)
	require.NoError(t, err)
	result, err := stmnt.Exec(db, setterFuncs...)
	require.NoError(t, err)
	require.NotNil(t, result)
	count, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(expectedRowsAffected), count)
}

func testPgDBCreds(t *testing.T) *url.URL {
	v := os.Getenv("MAPSQL_DB")
	if v == "" {
		v = "postgres://sqlstmt:sqlstmt@localhost:5432/sqlstmt_dev?sslmode=disable"
	}
	u, err := url.Parse(v)
	require.NoError(t, err)
	return u
}

func testConnectToDatabase(t *testing.T) Database {
	sqlDatabase, err := sql.Open("postgres", testPgDBCreds(t).String())
	require.NoError(t, err)
	return sqlDatabase
}

func TestExec(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc      string
		assertion func(*testing.T, string)
	}{
		{
			desc: "db.Prepare returned an error",
			assertion: func(t *testing.T, desc string) {
				statement, err := ConvertNamedToPositionalParams(insertTestingDataTypeQuery)
				require.NoError(t, err, desc)
				result, err := statement.Exec(&mockDB{
					PrepareOk: false,
				}, nil)
				require.Error(t, err, desc)
				require.Nil(t, result)
			},
		},
		{
			desc: "insert and delete then affected rows",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer db.Close()

				fake := faker.New()
				genUUID := uuid.New()

				setters := []SetParameterFunc{
					SetParameter("uuid", genUUID.String()),
					SetParameter("word", fake.Lorem().Text(10)),
					SetParameter("paragraph", fake.Lorem().Text(1000)),
					SetParameter("metadata", []byte(`{"test": "foo"}`)),
					SetParameter("created_at", time.Now().UTC()),
				}

				statement, err := ConvertNamedToPositionalParams(insertTestingDataTypeQuery)
				require.NoError(t, err, desc)

				result, err := statement.Exec(db, setters...)
				require.NoError(t, err, desc)
				affectedRows, err := result.RowsAffected()
				require.NoError(t, err, desc)
				require.Equal(t, int64(1), affectedRows)

				statement, err = ConvertNamedToPositionalParams(deleteTestingDataTypeQuery)
				require.NoError(t, err, desc)
				result, err = statement.Exec(
					db,
					SetParameter("uuid", genUUID.String()),
				)
				require.NoError(t, err, desc)
				require.NotNil(t, result, desc)
				affectedRows, err = result.RowsAffected()
				require.NoError(t, err, desc)
				require.Equal(t, int64(1), affectedRows)
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}

func TestQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		assertion func(*testing.T, string)
	}{
		{
			desc: "db.Prepare returned an error",
			assertion: func(t *testing.T, desc string) {
				statement, err := ConvertNamedToPositionalParams(createCustomerQuery)
				require.NoError(t, err, "failed to convert to positional params")
				rows, err := statement.Query(&mockDB{PrepareOk: false}, nil)
				require.Error(t, err, desc)
				require.Nil(t, rows)
			},
		},
		{
			desc: "create a customer and map rows",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer db.Close()

				fakeCustomer := genFakeCustomerData(t)

				statement, err := ConvertNamedToPositionalParams(createCustomerQuery)
				require.NoError(t, err, "failed to convert to positional params")

				rows, err := statement.Query(db, fakeCustomer.asParameters(t)...)
				require.NoError(t, err, desc)
				require.NotNil(t, rows)
				mappedRows, err := MapRows(rows)
				require.NoError(t, err, desc)
				require.NotNil(t, mappedRows)
				require.Equal(t, 1, mappedRows.Count())

				customerID, ok := mappedRows[0]["customer_id"]
				require.True(t, ok)
				require.NotEmpty(t, customerID)

				deleteRecords(
					t,
					db,
					1,
					deleteCustomerQuery,
					SetParameter("customer_id", customerID),
					SetParameter("email_address", fakeCustomer.ContactInfo.EmailAddress),
				)
			},
		},
		{
			desc: "query for record without parameters",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer db.Close()

				statement, err := ConvertNamedToPositionalParams(selectCustomerQuery)
				require.NoError(t, err, desc)
				rows, err := statement.Query(db, nil)
				require.NoError(t, err, desc)
				require.NotNil(t, rows)
				// NOTE: assert the returned values with more certainty.
				// right now, we are only checking if values are returned
				// by ensuring that there is a record being returned despite
				// the query executed without parameters.
				mappedRows, err := MapRows(rows)
				require.NoError(t, err, desc)
				require.NotNil(t, mappedRows)
				require.NotEmpty(t, mappedRows)
			},
		},
		{
			desc: "query with parameters",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer db.Close()

				createdCustomer := createCustomerForTesting(t)

				statement, err := ConvertNamedToPositionalParams(selectCustomerQuery)
				require.NoError(t, err, desc)
				rows, err := statement.Query(
					db,
					SetParameter("customer_id", createdCustomer.CustomerID),
					SetParameter("email_address", createdCustomer.ContactInfo.EmailAddress),
				)
				require.NoError(t, err, desc)
				require.NotNil(t, rows)
				mappedRows, err := MapRows(rows)
				require.NoError(t, err, desc)
				require.NotNil(t, mappedRows)
				require.NotEmpty(t, mappedRows)
				mappedRow := mappedRows[0]
				assertMappedCustomer(t, *createdCustomer, mappedRow)
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}

func TestQueryRow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		assertion func(*testing.T, string)
	}{
		{
			desc: "db.Prepare returned an error",
			assertion: func(t *testing.T, desc string) {
				statement, err := ConvertNamedToPositionalParams(createCustomerQuery)
				require.NoError(t, err, "failed to convert to positional params")
				row, err := statement.QueryRow(&mockDB{PrepareOk: false}, nil)
				require.Error(t, err, desc)
				require.Nil(t, row)
			},
		},
		{
			desc: "create a record and expect result",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer db.Close()

				fakeCustomer := genFakeCustomerData(t)

				statement, err := ConvertNamedToPositionalParams(createCustomerQuery)
				require.NoError(t, err, "failed to convert to positional params")

				row, err := statement.QueryRow(db, fakeCustomer.asParameters(t)...)
				require.NoError(t, err, desc)
				require.NotNil(t, row)

				columns := []string{
					"customer_id",
					"last_name",
					"first_name",
					"contact_info",
					"address",
				}

				mappedRow, err := MapRow(row, columns)
				require.NoError(t, err)
				require.NotNil(t, mappedRow)

				assertMappedCustomer(t, fakeCustomer, mappedRow)

				deleteRecords(
					t,
					db,
					1,
					deleteCustomerQuery,
					SetParameter("customer_id", mappedRow["customer_id"].(int64)),
					SetParameter("email_address", fakeCustomer.ContactInfo.EmailAddress),
				)
			},
		},
		{
			desc: "query record without parameters",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer db.Close()

				statement, err := ConvertNamedToPositionalParams(selectCustomerQuery)
				require.NoError(t, err)
				row, err := statement.QueryRow(db)
				require.NoError(t, err)
				require.NotNil(t, row)
			},
		},
		{
			desc: "query with parameters",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer db.Close()

				createdCustomer := createCustomerForTesting(t)

				statement, err := ConvertNamedToPositionalParams(selectCustomerQuery)
				require.NoError(t, err, desc)
				row, err := statement.QueryRow(
					db,
					SetParameter("customer_id", createdCustomer.CustomerID),
					SetParameter("email_address", createdCustomer.ContactInfo.EmailAddress),
				)
				require.NoError(t, err, desc)
				require.NotNil(t, row)
				mappedRow, err := MapRow(row, []string{
					"customer_id",
					"last_name",
					"first_name",
					"contact_info",
					"address",
				})
				require.NoError(t, err, desc)
				require.NotNil(t, mappedRow)
				require.NotEmpty(t, mappedRow)
				assertMappedCustomer(t, *createdCustomer, mappedRow)
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}
