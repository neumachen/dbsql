package sqlstmt

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestFieldSetter_Columns(t *testing.T) {
	td := testingDataType{}

	require.NotEmpty(t, td.ColumnMapperMap().Columns())
}

func TestSetField(t *testing.T) {
	db := testConnectToDatabase(t)
	defer testCloseDB(t, db)

	stmt, err := ConvertNamedToPositionalParams(
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

	stmt, err = ConvertNamedToPositionalParams(selectTestingDataTypeQuery)
	require.NoError(t, err)
	require.NotNil(t, stmt)

	row, err := stmt.QueryRow(db,
		SetParameter("uuids", pq.Array([]string{testData.UUID})),
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

	td := testingDataType{}
	err = td.mapRow(mappedRow)
	require.NoError(t, err)

	require.NotEmpty(t, td.ID)
	require.Equal(t, testData.UUID, td.UUID)
	require.Equal(t, testData.Word, td.Word)
	require.Equal(t, testData.Paragraph, td.Paragraph)
}
