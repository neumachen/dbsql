package dbsql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		assertion func(*testing.T, string)
	}{
		{
			desc: "PrepareStatement returns an error when the Prepare operation fails",
			assertion: func(t *testing.T, desc string) {
				preparedStatement, err := PrepareStatement(createCustomerQuery)
				require.NoError(t, err, "failed to convert to positional params")
				rows, err := QueryContext(
					context.TODO(),
					&mockDB{PrepareOk: false},
					preparedStatement,
					nil,
				)
				require.Error(t, err, desc)
				require.Nil(t, rows)
			},
		},
		{
			desc: "Create a customer and map the returned rows",
			assertion: func(t *testing.T, desc string) {
				db := ConnectToDatabase(t)
				defer CloseDB(t, db)

				fakeCustomer := NewFakeCustomerData(t)

				preparedStatement, err := PrepareStatement(createCustomerQuery)
				require.NoError(t, err, "failed to convert to positional params")

				rows, err := QueryContext(
					context.TODO(),
					db,
					preparedStatement,
					fakeCustomer.ParameterValues(t)...,
				)
				require.NoError(t, err, desc)
				require.NotNil(t, rows)
				mappedRows, err := MapRows(rows)
				require.NoError(t, err, desc)
				require.NotNil(t, mappedRows)
				require.Equal(t, 1, len(mappedRows))

				customerID, ok := mappedRows[0]["customer_id"]
				require.True(t, ok)
				require.NotEmpty(t, customerID)

				DeleteRecords(
					t,
					db,
					1,
					deleteCustomerQuery,
					BindParameterValue("customer_id", customerID),
					BindParameterValue("email_address", fakeCustomer.ContactInfo.EmailAddress),
				)
			},
		},
		{
			desc: "Query for records without any parameters",
			assertion: func(t *testing.T, desc string) {
				db := ConnectToDatabase(t)
				defer CloseDB(t, db)

				preparedStatement, err := PrepareStatement(selectCustomerQuery)
				require.NoError(t, err, desc)
				rows, err := QueryContext(
					context.TODO(),
					db,
					preparedStatement,
					nil,
				)
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
				db := ConnectToDatabase(t)
				defer CloseDB(t, db)

				createdCustomer := CreateNewCustomerForTesting(t)

				preparedStatement, err := PrepareStatement(selectCustomerQuery)
				require.NoError(t, err, desc)
				rows, err := QueryContext(
					context.TODO(),
					db,
					preparedStatement,
					BindParameterValue("customer_id", createdCustomer.CustomerID),
					BindParameterValue("email_address", createdCustomer.ContactInfo.EmailAddress),
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
