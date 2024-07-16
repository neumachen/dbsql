package dbsql

import (
	"context"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestMapRow(t *testing.T) {
	t.Parallel()

	db := ConnectToDatabase(t)
	defer CloseDB(t, db)

	preparedStatement, err := PrepareStatement(
		insertTestingDataTypeQuery,
	)
	require.NoError(t, err)
	require.NotNil(t, preparedStatement)

	testData := NewFakeTestingDataType(t)

	result, execErr := ExecContext(
		context.Background(),
		db,
		preparedStatement,
		testData.ParameterValues()...,
	)
	require.NoError(t, execErr)
	require.NotNil(t, result)
	rowsAffected, rowErr := result.RowsAffected()
	require.NoError(t, rowErr)
	require.Equal(t, int64(1), rowsAffected)

	preparedStatement, err = PrepareStatement(selectTestingDataTypeQuery)
	require.NoError(t, err)
	require.NotNil(t, preparedStatement)

	row, err := QueryRowContext(
		context.TODO(),
		db,
		preparedStatement,
		BindParameterValue("uuids", pq.Array([]string{testData.UUID})),
	)
	require.NoError(t, err)
	require.NotNil(t, row)
	expectedColumns := Columns{
		"testing_datatype_id",
		"testing_datatype_uuid",
		"word",
		"paragraph",
		"metadata",
		"created_at",
	}
	mappedRow, err := MapRow(row, expectedColumns)
	require.NoError(t, err)
	require.NotNil(t, mappedRow)

	for _, expectedColumn := range expectedColumns {
		require.True(t, mappedRow.HasColumn(expectedColumn))
	}
}
