package dbsql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryRowContext(t *testing.T) {
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
				row, err := QueryRowContext(
					context.TODO(),
					&mockDB{
						PrepareOk: false,
					},
					preparedStatement,
					nil,
				)
				require.Error(t, err, desc)
				require.Nil(t, row)
			},
		},
		{
			desc: "Create a record and expect the result",
			assertion: func(t *testing.T, desc string) {
				db := ConnectToDatabase(t)
				defer CloseDB(t, db)

				fakeCustomer := NewFakeCustomerData(t)

				preparedStatement, err := PrepareStatement(createCustomerQuery)
				require.NoError(t, err, "failed to convert to positional params")

				row, err := QueryRowContext(
					context.TODO(),
					db,
					preparedStatement,
					fakeCustomer.ParameterValues(t)...,
				)
				require.NoError(t, err, desc)
				require.NotNil(t, row)

				columns := Columns{
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

				DeleteRecords(
					t,
					db,
					1,
					deleteCustomerQuery,
					BindParameterValue("customer_id", mappedRow["customer_id"].(int64)),
					BindParameterValue("email_address", fakeCustomer.ContactInfo.EmailAddress),
				)
			},
		},
		{
			desc: "Query a record without any parameters",
			assertion: func(t *testing.T, desc string) {
				db := ConnectToDatabase(t)
				defer CloseDB(t, db)

				preparedStatement, err := PrepareStatement(selectCustomerQuery)
				require.NoError(t, err, desc)
				row, err := QueryRowContext(
					context.TODO(),
					db,
					preparedStatement,
				)
				require.NoError(t, err, desc)
				require.NotNil(t, row)
			},
		},
		{
			desc: "Query a record with parameters",
			assertion: func(t *testing.T, desc string) {
				db := ConnectToDatabase(t)
				defer CloseDB(t, db)

				createdCustomer := CreateNewCustomerForTesting(t)

				preparedStatement, err := PrepareStatement(selectCustomerQuery)
				require.NoError(t, err, desc)
				row, err := QueryRowContext(
					context.TODO(),
					db,
					preparedStatement,
					BindParameterValue("customer_id", createdCustomer.CustomerID),
					BindParameterValue("email_address", createdCustomer.ContactInfo.EmailAddress),
				)
				require.NoError(t, err, desc)
				require.NotNil(t, row)
				mappedRow, err := MapRow(row, Columns{
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
