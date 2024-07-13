package dbsql

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/neumachen/randata"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

type testingDataType struct {
	ID        int64
	UUID      string
	Word      string
	Paragraph string
	Metadata  json.RawMessage
	CreatedAt time.Time
}

func (t *testingDataType) ColumnMapperMap() ColumnBinderMap {
	return ColumnBinderMap{
		"testing_datatype_id": BindColumnToField[int64](func(value int64) error {
			t.ID = value
			return nil
		},
		),
		"testing_datatype_uuid": BindColumnToField[[]uint8](func(value []uint8) error {
			t.UUID = string(value)
			return nil
		},
		),
		"word": BindColumnToField[string](func(value string) error {
			t.Word = value
			return nil
		},
		),
		"paragraph": BindColumnToField[string](func(value string) error {
			t.Paragraph = value
			return nil
		},
		),
		"metadata": BindColumnToField[[]byte](func(value []byte) error {
			t.Metadata = json.RawMessage(value)
			return nil
		},
		),
		"created_at": BindColumnToField[time.Time](func(value time.Time) error {
			t.CreatedAt = value
			return nil
		},
		),
	}
}

func (t *testingDataType) mapRow(row MappedRow) error {
	if row.Count() < 1 {
		return nil
	}

	for column, mapperFunc := range t.ColumnMapperMap() {
		if err := mapperFunc(column, row); err != nil {
			return err
		}
	}
	return nil
}

func (t testingDataType) asParameters() []BindParameterValueFunc {
	return []BindParameterValueFunc{
		BindParameterValue("uuid", t.UUID),
		BindParameterValue("word", t.Word),
		BindParameterValue("paragraph", t.Paragraph),
		BindParameterValue("metadata", t.Metadata),
		BindParameterValue("created_at", t.CreatedAt),
	}
}

func genFakeTestingDataType(t *testing.T) *testingDataType {
	genUUID := uuid.New()
	fake := faker.New()

	mapData := map[string]any{
		fake.Lorem().Text(1): fake.Lorem().Text(1),
	}

	b, err := json.Marshal(mapData)
	require.NoError(t, err)
	metadata := json.RawMessage(b)

	return &testingDataType{
		UUID:      genUUID.String(),
		Word:      fake.Lorem().Text(10),
		Paragraph: fake.Lorem().Text(100),
		Metadata:  metadata,
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

func (c customer) asParameters(t *testing.T) []BindParameterValueFunc {
	contactInfo, err := json.Marshal(c.ContactInfo)
	require.NoError(t, err)
	address, err := json.Marshal(c.Address)
	require.NoError(t, err)
	return []BindParameterValueFunc{
		BindParameterValue("last_name", c.LastName),
		BindParameterValue("first_name", c.FirstName),
		BindParameterValue("contact_info", contactInfo),
		BindParameterValue("address", address),
	}
}

func genFakeCustomerData(t *testing.T) customer {
	t.Helper()

	// NOTE: given we are running t.Parallel() on tests, to avoid race conditions when
	// creating a new customer with the email address uniquness, we need to ensure only
	// on routine is using that email address at a time before generating a new one.
	var r sync.Mutex
	r.Lock()
	defer r.Unlock()

	randomAddress, err := randata.USAddress()
	require.NoError(t, err)

	fake := faker.NewWithSeed(rand.NewSource(time.Now().UTC().UnixNano()))

	guid := xid.New()
	emailAddress := fmt.Sprintf("%s@%s", guid.String(), fake.Internet().Domain())

	c := customer{}
	c.LastName = fake.Person().LastName()
	c.FirstName = fake.Person().FirstName()
	c.ContactInfo = &contactInfo{
		EmailAddress: emailAddress,
	}
	c.Address = &address{
		Address: *randomAddress,
	}

	return c
}

func createCustomerForTesting(t *testing.T) *customer {
	t.Helper()
	db := testConnectToDatabase(t)
	defer testCloseDB(t, db)

	createCustomer := genFakeCustomerData(t)

	preparedStatement, err := PrepareStatement(createCustomerQuery)
	require.NoError(t, err, "failed to convert to positional params")

	rows, err := QueryContext(
		context.TODO(),
		db,
		preparedStatement,
		createCustomer.asParameters(t)...,
	)
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
