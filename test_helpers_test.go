package dbsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/neumachen/dbsql/internal"
	"github.com/stretchr/testify/require"
)

// testPgDBCreds retrieves the PostgreSQL database credentials from the environment
// or uses a default connection string if not set.
func testPgDBCreds(t *testing.T) *url.URL {
	v := os.Getenv("")
	if v == "" {
		v = "postgres://dbsql:dbsql@localhost:5432/dbsql_dev?sslmode=disable"
	}
	u, err := url.Parse(v)
	require.NoError(t, err)
	return u
}

// testConnectToDatabase creates a new SQL database connection using the test PostgreSQL credentials.
func testConnectToDatabase(t *testing.T) SQLExecutor {
	sqlDatabase, err := sql.Open("postgres", testPgDBCreds(t).String())
	require.NoError(t, err)
	return sqlDatabase
}

// testCloseDB closes the provided SQL database connection, if it's not nil.
func testCloseDB(t *testing.T, db SQLExecutor) {
	if internal.IsNilOrZeroValue(db) {
		err := db.Close()
		require.NoError(t, err)
	}
}

// areEqualJSON compares two JSON strings for equality.
func areEqualJSON(s1, s2 string) (bool, error) {
	var o1, o2 interface{}
	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error marshalling string 1: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error marshalling string 2: %s", err.Error())
	}
	return reflect.DeepEqual(o1, o2), nil
}

// assertMappedCustomer checks that the mapped customer data matches the expected customer.
func assertMappedCustomer(t *testing.T, expectedData customer, mappedRow MappedRow) {
	require.NotEmpty(t, mappedRow["customer_id"].(int64))
	require.Equal(t, expectedData.LastName, mappedRow["last_name"].(string))
	require.Equal(t, expectedData.FirstName, mappedRow["first_name"].(string))
	cInfo := contactInfo{}
	err := json.Unmarshal(mappedRow["contact_info"].([]byte), &cInfo)
	require.NoError(t, err)
	require.NotEmpty(t, cInfo.EmailAddressID)
	require.Equal(t, expectedData.ContactInfo.EmailAddress, cInfo.EmailAddress)
	addr := address{}
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

func deleteRecords(
	t *testing.T,
	db SQLExecutor,
	expectedRowsAffected int,
	query string,
	binderFuncs ...BindNamedParameterValueFunc,
) {
	stmnt, err := PrepareStatement(query)
	require.NoError(t, err)
	result, err := stmnt.Exec(db, binderFuncs...)
	require.NoError(t, err)
	require.NotNil(t, result)
	count, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(expectedRowsAffected), count)
}
