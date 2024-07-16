package dbsql

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
)

func TestExecContext(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		assertion func(*testing.T, string)
	}{
		{
			name: "PrepareStatement returns an error when the Prepare operation fails",
			assertion: func(t *testing.T, desc string) {
				preparedStatement, err := PrepareStatement(insertTestingDataTypeQuery)
				require.NoError(t, err, desc)
				result, err := ExecContext(
					context.Background(),
					&mockDB{
						PrepareOk: false,
					},
					preparedStatement,
					nil,
				)
				require.Error(t, err, desc)
				require.Nil(t, result)
			},
		},
		{
			name: "Insert and delete records with the expected affected rows",
			assertion: func(t *testing.T, desc string) {
				db := ConnectToDatabase(t)
				defer CloseDB(t, db)

				fake := faker.New()
				genUUID := uuid.New()

				setters := []BindParameterValueFunc{
					BindParameterValue("uuid", genUUID.String()),
					BindParameterValue("word", fake.Lorem().Text(10)),
					BindParameterValue("paragraph", fake.Lorem().Text(1000)),
					BindParameterValue("metadata", []byte(`{"test": "foo"}`)),
					BindParameterValue("created_at", time.Now().UTC()),
				}

				preparedStatement, err := PrepareStatement(insertTestingDataTypeQuery)
				require.NoError(t, err, desc)

				result, err := ExecContext(
					context.Background(),
					db,
					preparedStatement,
					setters...,
				)
				require.NoError(t, err, desc)
				affectedRows, err := result.RowsAffected()
				require.NoError(t, err, desc)
				require.Equal(t, int64(1), affectedRows)

				preparedStatement, err = PrepareStatement(deleteTestingDataTypeQuery)
				require.NoError(t, err, desc)
				result, err = ExecContext(
					context.Background(),
					db,
					preparedStatement,
					BindParameterValue("uuid", genUUID.String()),
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
