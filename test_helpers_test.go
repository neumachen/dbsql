package sqlstmt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func areEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

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
