package sqlstmt

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type ConvertPositionalParamsTest struct {
	Name                       string
	Input                      string
	ExpectedStatement          string
	ExpectedParameterPositions *ParameterPositions
}

func TestConvertPositionalParams(t *testing.T) {
	convertPositionalParamsTests := []ConvertPositionalParamsTest{
		{
			Input:                      "SELECT * FROM table WHERE col1 = 1",
			ExpectedStatement:          "SELECT * FROM table WHERE col1 = 1",
			ExpectedParameterPositions: nil,
			Name:                       "No Parameter",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 = :name",
			ExpectedStatement: "SELECT * FROM table WHERE col1 = $1",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"name": {0},
				},
				totalPositions: 1,
			},
			Name: "Single Parameter",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 = :name AND col2 = :occupation",
			ExpectedStatement: "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"name":       {0},
					"occupation": {1},
				},
				totalPositions: 2,
			},
			Name: "Two Parameters",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 = :name AND col2 = :name",
			ExpectedStatement: "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"name": {0, 1},
				},
				totalPositions: 2,
			},
			Name: "Repeated Named Parameter",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 IN (:something, :else)",
			ExpectedStatement: "SELECT * FROM table WHERE col1 IN ($1, $2)",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"something": {0},
					"else":      {1},
				},
				totalPositions: 2,
			},
			Name: "Parameters In Parenthesis",
		},
		{
			Input:                      "SELECT * FROM table WHERE col1 = ':literal' AND col2 LIKE ':literal'",
			ExpectedStatement:          "SELECT * FROM table WHERE col1 = ':literal' AND col2 LIKE ':literal'",
			ExpectedParameterPositions: nil,
			Name:                       "Escaped Parameters",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 = ':literal' AND col2 = :literal AND col3 LIKE ':literal'",
			ExpectedStatement: "SELECT * FROM table WHERE col1 = ':literal' AND col2 = $1 AND col3 LIKE ':literal'",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"literal": {0},
				},
				totalPositions: 1,
			},
			Name: "Escaping and non-escaping parameters",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 = :foo AND col2 IN (SELECT id FROM tabl2 WHERE col10 = :bar)",
			ExpectedStatement: "SELECT * FROM table WHERE col1 = $1 AND col2 IN (SELECT id FROM tabl2 WHERE col10 = $2)",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"foo": {0},
					"bar": {1},
				},
				totalPositions: 2,
			},
			Name: "Parameters In Subclause",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 = :1234567890 AND col2 = :0987654321",
			ExpectedStatement: "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"1234567890": {0},
					"0987654321": {1},
				},
				totalPositions: 2,
			},
			Name: "Numeric Parameters",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 = :ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			ExpectedStatement: "SELECT * FROM table WHERE col1 = $1",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"ABCDEFGHIJKLMNOPQRSTUVWXYZ": {0},
				},
				totalPositions: 1,
			},
			Name: "Upcased Parameter",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 = :abc123ABC098",
			ExpectedStatement: "SELECT * FROM table WHERE col1 = $1",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"abc123ABC098": {0},
				},
				totalPositions: 1,
			},
			Name: "Multicased alphanumeric parameter",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 LIKE %:t%",
			ExpectedStatement: "SELECT * FROM table WHERE col1 LIKE %$1%",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"t": {0},
				},
				totalPositions: 1,
			},
			Name: "Pattern Matching",
		},
		{
			Input:             "ST_GeomFromText('POINT(' || :long :lat || ',4326)'",
			ExpectedStatement: "ST_GeomFromText('POINT(' || $1 $2 || ',4326)'",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"long": {0},
					"lat":  {1},
				},
				totalPositions: 2,
			},
			Name: "Concated parameters",
		},
		{
			Input:             "SELECT * FROM table WHERE col1 = :first_name AND col2 = :last_name",
			ExpectedStatement: "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameterPositions: &ParameterPositions{
				paramPositions: map[string][]int{
					"first_name": {0},
					"last_name":  {1},
				},
				totalPositions: 2,
			},
			Name: "Snake Case",
		},
	}

	for _, convertTest := range convertPositionalParamsTests {
		stmnt, err := ConvertNamedToPositionalParams([]byte(convertTest.Input))
		require.NoError(t, err)
		require.Equal(t, convertTest.ExpectedStatement, string(stmnt.GetQueryWithPositionals()), convertTest.Name)
		actualParamPositions := stmnt.GetParameterPositions()
		require.Equal(t, convertTest.ExpectedParameterPositions, actualParamPositions, convertTest.Name)
		// for k, v := range convertTest.ExpectedParameterPositions {
		// 	positions := actualParamPositions.getPositions(k)
		// 	require.Equal(t, v, positions)
		// }
	}
}

type ParameterParsingTest struct {
	Name               string
	Query              string
	Parameters         []TestQueryParameter
	ExpectedParameters []interface{}
}

type TestQueryParameter struct {
	Name  string
	Value interface{}
}

func TestSettingParameters(t *testing.T) {
	queryVariableTests := []ParameterParsingTest{
		{
			Name:  "Single String Parameter",
			Query: "SELECT * FROM table WHERE col1 = :foo",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "bar",
				},
			},
			ExpectedParameters: []interface{}{
				"bar",
			},
		},
		{
			Name:  "Two String Parameter",
			Query: "SELECT * FROM table WHERE col1 = :foo AND col2 = :foo2",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "bar",
				},
				{
					Name:  "foo2",
					Value: "bart",
				},
			},
			ExpectedParameters: []interface{}{
				"bar",
				"bart",
			},
		},
		{
			Name:  "Repeated Parameters",
			Query: "SELECT * FROM table WHERE col1 = :foo AND col2 = :foo",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "bar",
				},
			},
			ExpectedParameters: []interface{}{
				"bar",
				"bar",
			},
		},
		{
			Name:  "Type Parameters",
			Query: "SELECT * FROM table WHERE col1 = :str AND col2 = :int AND col3 = :pi",
			Parameters: []TestQueryParameter{
				{
					Name:  "str",
					Value: "foo",
				},
				{
					Name:  "int",
					Value: 1,
				},
				{
					Name:  "pi",
					Value: 3.14,
				},
			},
			ExpectedParameters: []interface{}{
				"foo",
				1,
				3.14,
			},
		},
		{
			Name:  "Ordered Parameters",
			Query: "SELECT * FROM table WHERE col1 = :foo AND col2 = :bar AND col3 = :foo AND col4 = :foo AND col5 = :bar",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "something",
				},
				{
					Name:  "bar",
					Value: "else",
				},
			},
			ExpectedParameters: []interface{}{
				"something", "else", "something", "something", "else",
			},
		},
		{
			Name:  "Case Sensitive",
			Query: "SELECT * FROM table WHERE col1 = :foo AND col2 = :FOO",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "baz",
				},
				{
					Name:  "FOO",
					Value: "quux",
				},
			},
			ExpectedParameters: []interface{}{
				"baz", "quux",
			},
		},
		{
			Name:  "Nil Parameter",
			Query: "SELECT * FROM table WHERE col1 = :foo",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: pq.Array([]string{}),
				},
			},
			ExpectedParameters: []interface{}{
				pq.Array([]string{}),
			},
		},
		{
			Name:  "Casted Type Parameter",
			Query: "SELECT * FROM table WHERE col1 = :foo",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "'testing'::varchar",
				},
			},
			ExpectedParameters: []interface{}{
				"'testing'::varchar",
			},
		},
	}

	for _, variableTest := range queryVariableTests {
		parameterFuncs := make([]SetParameterFunc, 0)
		stmt, err := ConvertNamedToPositionalParams([]byte(variableTest.Query))
		require.NoError(t, err)

		for _, queryVariable := range variableTest.Parameters {
			err = stmt.SetParameter(queryVariable.Name, queryVariable.Value)
			require.NoError(t, err)
			parameterFuncs = append(parameterFuncs, SetParameter(queryVariable.Name, queryVariable.Value))
		}

		err = stmt.SetParameters(parameterFuncs...)
		require.NoError(t, err)

		for posIndex, parameterValue := range stmt.GetPositionalParameters() {
			require.Equal(t, parameterValue, variableTest.ExpectedParameters[posIndex], variableTest.Name)
		}
	}
}
