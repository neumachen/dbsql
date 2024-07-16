package dbsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/url"
	"os"
	"testing"

	"github.com/neumachen/dbsql/internal"
	"github.com/stretchr/testify/require"
)

// PostgresDBCreds retrieves the PostgreSQL database credentials from the environment
// or uses a default connection string if not set.
func PostgresDBCreds(t *testing.T) *url.URL {
	v := os.Getenv("")
	if v == "" {
		v = "postgres://dbsql:dbsql@localhost:5432/dbsql_dev?sslmode=disable"
	}
	u, err := url.Parse(v)
	require.NoError(t, err)
	return u
}

// ConnectToDatabase creates a new SQL database connection using the test PostgreSQL credentials.
func ConnectToDatabase(t *testing.T) interface {
	DBPreparerExecutor
	DBCloser
} {
	sqlDatabase, err := sql.Open("postgres", PostgresDBCreds(t).String())
	require.NoError(t, err)
	return sqlDatabase
}

// CloseDB closes the provided SQL database connection, if it's not nil.
func CloseDB(
	t *testing.T,
	db interface {
		DBPreparerExecutor
		DBCloser
	},
) {
	if internal.IsNilOrZeroValue(db) {
		err := db.Close()
		require.NoError(t, err)
	}
}

// assertMappedCustomer checks that the mapped customer data matches the expected customer.
func assertMappedCustomer(t *testing.T, expectedData Customer, mappedRow MappedRow) {
	require.NotEmpty(t, mappedRow["customer_id"].(int64))
	require.Equal(t, expectedData.LastName, mappedRow["last_name"].(string))
	require.Equal(t, expectedData.FirstName, mappedRow["first_name"].(string))

	cInfo := ContactInfo{}
	err := json.Unmarshal(mappedRow["contact_info"].([]byte), &cInfo)
	require.NoError(t, err)
	require.NotEmpty(t, cInfo.EmailAddressID)
	require.Equal(t, expectedData.ContactInfo.EmailAddress, cInfo.EmailAddress)

	addr := Address{}
	err = json.Unmarshal(mappedRow["address"].([]byte), &addr)
	require.NoError(t, err)
	require.NotEmpty(t, addr.AddressID)
	require.Equal(t, expectedData.Address.StreetNumber, addr.StreetNumber)
	require.Equal(t, expectedData.Address.Route, addr.Route)
	require.Equal(t, expectedData.Address.Locality, addr.Locality)
	require.Equal(t, expectedData.Address.AdministrativeAreaLevel1, addr.AdministrativeAreaLevel1)
	require.Equal(t, expectedData.Address.PostalCode, addr.PostalCode)
	require.Equal(t, expectedData.Address.Latitude, addr.Latitude)
	require.Equal(t, expectedData.Address.Longitude, addr.Longitude)
}

func DeleteRecords(
	t *testing.T,
	db DBPreparerExecutor,
	expectedRowsAffected int,
	query string,
	binderFuncs ...BindParameterValueFunc,
) {
	preparedStatement, err := PrepareStatement(query)
	require.NoError(t, err)
	result, err := ExecContext(
		context.TODO(),
		db,
		preparedStatement,
		binderFuncs...,
	)
	require.NoError(t, err)
	require.NotNil(t, result)
	count, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(expectedRowsAffected), count)
}
