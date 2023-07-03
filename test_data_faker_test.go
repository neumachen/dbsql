package sqlstmt

import (
	"encoding/json"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/neumachen/randata"
	"github.com/stretchr/testify/require"
)

type testingDataType struct {
	UUID      string
	Word      string
	Paragraph string
	Metadata  *json.RawMessage
	CreatedAt time.Time
}

func (t testingDataType) asParameters() []SetParameterFunc {
	return []SetParameterFunc{
		SetParameter("uuid", t.UUID),
		SetParameter("word", t.Word),
		SetParameter("paragraph", t.Paragraph),
		SetParameter("metadata", t.Metadata),
		SetParameter("created_at", t.CreatedAt),
	}
}

func genFakeTestingDataType(t *testing.T) *testingDataType {
	genUUID := uuid.New()
	fake := faker.New()

	mapData := map[string]interface{}{
		fake.Lorem().Text(1): fake.Lorem().Text(1),
	}

	b, err := json.Marshal(mapData)
	require.NoError(t, err)
	metadata := json.RawMessage(b)

	return &testingDataType{
		UUID:      genUUID.String(),
		Word:      fake.Lorem().Text(10),
		Paragraph: fake.Lorem().Text(100),
		Metadata:  &metadata,
		CreatedAt: time.Now().UTC(),
	}
}

type address struct {
	randata.Address
	AddressID int `json:"address_id,omitempty"`
}

type contactInfo struct {
	EmailAddressID int    `json:"email_address_id,omitempty"`
	EmailAddress   string `json:"email_address"`
}

type customer struct {
	CustomerID  int64        `json:"customer_id,omitempty"`
	LastName    string       `json:"last_name,omitempty"`
	FirstName   string       `json:"first_name,omitempty"`
	ContactInfo *contactInfo `json:"contact_info,omitempty"`
	Address     *address     `json:"address,omitempty"`
}

func (c customer) asParameters(t *testing.T) []SetParameterFunc {
	contactInfo, err := json.Marshal(c.ContactInfo)
	require.NoError(t, err)
	address, err := json.Marshal(c.Address)
	require.NoError(t, err)
	return []SetParameterFunc{
		SetParameter("last_name", c.LastName),
		SetParameter("first_name", c.FirstName),
		SetParameter("contact_info", contactInfo),
		SetParameter("address", address),
	}
}

func genFakeCustomerData(t *testing.T) customer {
	var r sync.Mutex
	r.Lock()
	defer r.Unlock()

	randomAddress, err := randata.USAddress()
	require.NoError(t, err)

	fake := faker.NewWithSeed(rand.NewSource(time.Now().UTC().UnixNano()))

	c := customer{}
	c.LastName = fake.Person().LastName()
	c.FirstName = fake.Person().FirstName()
	c.ContactInfo = &contactInfo{
		EmailAddress: fake.Internet().Email(),
	}
	c.Address = &address{
		Address: *randomAddress,
	}

	return c
}

func createCustomerForTesting(t *testing.T) *customer {
	t.Helper()
	db := testConnectToDatabase(t)
	defer db.Close()

	createCustomer := genFakeCustomerData(t)

	statement, err := ConvertNamedToPositionalParams(createCustomerQuery)
	require.NoError(t, err, "failed to convert to positional params")

	rows, err := statement.Query(db, createCustomer.asParameters(t)...)
	require.NoError(t, err, "failed to query creating customer")
	require.NotNil(t, rows)
	mappedRows, err := MapRows(rows)
	require.NoError(t, err, "failed to map rows")
	require.NotNil(t, mappedRows)
	require.Equal(t, 1, mappedRows.Count())

	addr := address{}
	err = json.Unmarshal(mappedRows[0]["address"].([]uint8), &addr)
	require.NoError(t, err)

	createdCustomer := &customer{
		CustomerID: mappedRows[0]["customer_id"].(int64),
		LastName:   mappedRows[0]["last_name"].(string),
		FirstName:  mappedRows[0]["first_name"].(string),
		Address:    &addr,
	}

	b, err := json.Marshal(createCustomer)
	require.NoError(t, err)

	err = json.Unmarshal(b, &createdCustomer)
	require.NoError(t, err)

	return createdCustomer
}
