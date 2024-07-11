package sqlstmt

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestMapRow(t *testing.T) {
	t.Parallel()

	db := testConnectToDatabase(t)
	defer testCloseDB(t, db)

	stmt, err := PrepareStatement(
		insertTestingDataTypeQuery,
	)
	require.NoError(t, err)
	require.NotNil(t, stmt)

	testData := genFakeTestingDataType(t)

	result, execErr := stmt.Exec(
		db,
		testData.asParameters()...,
	)
	require.NoError(t, execErr)
	require.NotNil(t, result)
	rowsAffected, rowErr := result.RowsAffected()
	require.NoError(t, rowErr)
	require.Equal(t, int64(1), rowsAffected)

	stmt, err = PrepareStatement(selectTestingDataTypeQuery)
	require.NoError(t, err)
	require.NotNil(t, stmt)

	row, err := stmt.QueryRow(db,
		BindNamedParameterValue("uuids", pq.Array([]string{testData.UUID})),
	)
	require.NoError(t, err)
	require.NotNil(t, row)
	columns := []string{
		"testing_datatype_id",
		"testing_datatype_uuid",
		"word",
		"paragraph",
		"metadata",
		"created_at",
	}
	mappedRow, err := MapRow(row, columns)
	require.NoError(t, err)
	require.NotNil(t, mappedRow)
}
