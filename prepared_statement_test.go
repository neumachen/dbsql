package sqldb

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		assertion func(*testing.T, string)
	}{
		{
			name: "PrepareStatement returns an error when the Prepare operation fails",
			assertion: func(t *testing.T, desc string) {
				statement, err := PrepareStatement(insertTestingDataTypeQuery)
				require.NoError(t, err, desc)
				result, err := statement.Exec(&mockDB{
					PrepareOk: false,
				}, nil)
				require.Error(t, err, desc)
				require.Nil(t, result)
			},
		},
		{
			name: "Insert and delete records with the expected affected rows",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer testCloseDB(t, db)

				fake := faker.New()
				genUUID := uuid.New()

				setters := []BindNamedParameterValueFunc{
					BindNamedParameterValue("uuid", genUUID.String()),
					BindNamedParameterValue("word", fake.Lorem().Text(10)),
					BindNamedParameterValue("paragraph", fake.Lorem().Text(1000)),
					BindNamedParameterValue("metadata", []byte(`{"test": "foo"}`)),
					BindNamedParameterValue("created_at", time.Now().UTC()),
				}

				statement, err := PrepareStatement(insertTestingDataTypeQuery)
				require.NoError(t, err, desc)

				result, err := statement.Exec(db, setters...)
				require.NoError(t, err, desc)
				affectedRows, err := result.RowsAffected()
				require.NoError(t, err, desc)
				require.Equal(t, int64(1), affectedRows)

				statement, err = PrepareStatement(deleteTestingDataTypeQuery)
				require.NoError(t, err, desc)
				result, err = statement.Exec(
					db,
					BindNamedParameterValue("uuid", genUUID.String()),
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
		t.Run(test.name, func(t *testing.T) {
			test.assertion(t, test.name)
		})
	}
}

func TestQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		assertion func(*testing.T, string)
	}{
		{
			desc: "PrepareStatement returns an error when the Prepare operation fails",
			assertion: func(t *testing.T, desc string) {
				statement, err := PrepareStatement(createCustomerQuery)
				require.NoError(t, err, "failed to convert to positional params")
				rows, err := statement.Query(&mockDB{PrepareOk: false}, nil)
				require.Error(t, err, desc)
				require.Nil(t, rows)
			},
		},
		{
			desc: "Create a customer and map the returned rows",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer testCloseDB(t, db)

				fakeCustomer := genFakeCustomerData(t)

				statement, err := PrepareStatement(createCustomerQuery)
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
					BindNamedParameterValue("customer_id", customerID),
					BindNamedParameterValue("email_address", fakeCustomer.ContactInfo.EmailAddress),
				)
			},
		},
		{
			desc: "Query for records without any parameters",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer testCloseDB(t, db)

				statement, err := PrepareStatement(selectCustomerQuery)
				require.NoError(t, err, desc)
				rows, err := statement.Query(db, nil)
				require.NoError(t, err, desc)
				require.NotNil(t, rows)

				// Assert that the returned rows contain the expected data
				mappedRows, err := MapRows(rows)
				require.NoError(t, err, desc)
				require.NotNil(t, mappedRows)
				require.Greater(t, len(mappedRows), 0, "expected at least one customer record to be returned")

				// Assert the specific values in the first returned row
				firstRow := mappedRows[0]
				require.NotEmpty(t, firstRow["customer_id"], "customer_id should not be empty")
				require.NotEmpty(t, firstRow["last_name"], "last_name should not be empty")
				require.NotEmpty(t, firstRow["first_name"], "first_name should not be empty")
				require.NotEmpty(t, firstRow["contact_info"], "contact_info should not be empty")
				require.NotEmpty(t, firstRow["address"], "address should not be empty")
			},
		},
		{
			desc: "Query for records with parameters",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer testCloseDB(t, db)

				createdCustomer := createCustomerForTesting(t)

				statement, err := PrepareStatement(selectCustomerQuery)
				require.NoError(t, err, desc)
				rows, err := statement.Query(
					db,
					BindNamedParameterValue("customer_id", createdCustomer.CustomerID),
					BindNamedParameterValue("email_address", createdCustomer.ContactInfo.EmailAddress),
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
		t.Run(test.desc, func(t *testing.T) {
			test.assertion(t, test.desc)
		})
	}
}

func TestQueryRow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		assertion func(*testing.T, string)
	}{
		{
			desc: "PrepareStatement returns an error when the Prepare operation fails",
			assertion: func(t *testing.T, desc string) {
				statement, err := PrepareStatement(createCustomerQuery)
				require.NoError(t, err, "failed to convert to positional params")
				row, err := statement.QueryRow(&mockDB{PrepareOk: false}, nil)
				require.Error(t, err, desc)
				require.Nil(t, row)
			},
		},
		{
			desc: "Create a record and expect the result",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer testCloseDB(t, db)

				fakeCustomer := genFakeCustomerData(t)

				statement, err := PrepareStatement(createCustomerQuery)
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
					BindNamedParameterValue("customer_id", mappedRow["customer_id"].(int64)),
					BindNamedParameterValue("email_address", fakeCustomer.ContactInfo.EmailAddress),
				)
			},
		},
		{
			desc: "Query a record without any parameters",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer testCloseDB(t, db)

				statement, err := PrepareStatement(selectCustomerQuery)
				require.NoError(t, err, desc)
				row, err := statement.QueryRow(db)
				require.NoError(t, err, desc)
				require.NotNil(t, row)
			},
		},
		{
			desc: "Query a record with parameters",
			assertion: func(t *testing.T, desc string) {
				db := testConnectToDatabase(t)
				defer testCloseDB(t, db)

				createdCustomer := createCustomerForTesting(t)

				statement, err := PrepareStatement(selectCustomerQuery)
				require.NoError(t, err, desc)
				row, err := statement.QueryRow(
					db,
					BindNamedParameterValue("customer_id", createdCustomer.CustomerID),
					BindNamedParameterValue("email_address", createdCustomer.ContactInfo.EmailAddress),
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
		t.Run(test.desc, func(t *testing.T) {
			test.assertion(t, test.desc)
		})
	}
}
