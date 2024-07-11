package sqldb

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestMapRows(t *testing.T) {
	t.Parallel()

	db := testConnectToDatabase(t)
	defer testCloseDB(t, db)

	stmt, err := PrepareStatement(
		insertTestingDataTypeQuery,
	)
	require.NoError(t, err)
	require.NotNil(t, stmt)

	count := 3
	uuids := make([]string, count)

	expectedData := make(map[string]testingDataType)
	for i := 0; i < count; i++ {
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
		expectedData[testData.UUID] = *testData
		uuids[i] = testData.UUID
	}

	stmt, err = PrepareStatement(selectTestingDataTypeQuery)
	require.NoError(t, err)
	require.NotNil(t, stmt)

	rows, err := stmt.Query(db,
		BindNamedParameterValue("uuids", pq.Array(uuids)),
	)
	require.NoError(t, err)
	require.NotNil(t, rows)

	mappedRows, err := MapRows(rows)
	require.NoError(t, err)
	require.NotNil(t, mappedRows)
	require.Equal(t, count, mappedRows.Count())

	for i := range mappedRows {
		mappedRow := mappedRows[i]
		require.NotEmpty(t, mappedRow["testing_datatype_id"].(int64))
		value, ok := mappedRow["testing_datatype_uuid"]
		require.True(t, ok)
		require.NotEmpty(t, value)
		idStr := string(value.([]uint8))

		data, ok := expectedData[idStr]
		require.True(t, ok)
		require.NotNil(t, data)

		require.Equal(t, data.UUID, idStr)
		require.Equal(t, data.Word, mappedRow["word"].(string))
		require.Equal(t, data.Paragraph, mappedRow["paragraph"].(string))
		require.Equal(t, data.CreatedAt, mappedRow["created_at"].(time.Time))
		mData := json.RawMessage(mappedRow["metadata"].([]byte))
		isEqual, err := areEqualJSON(string(data.Metadata), string(mData))
		require.NoError(t, err)
		require.True(t, isEqual)
	}

	stmt, err = PrepareStatement(deleteTestingDataTypeQuery)
	require.NoError(t, err)
	require.NotNil(t, stmt)

	result, err := stmt.Exec(db,
		BindNamedParameterValue("uuids", pq.Array(uuids)),
	)
	require.NoError(t, err)
	require.NotNil(t, result)
	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(count), rowsAffected)
}
