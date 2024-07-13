package dbsql

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type PrepareStatementTest struct {
	Name                       string
	UnpreparedStatement        string
	ExpectedStatement          string
	ExpectedParameterPositions *NamedParameterPositions
}

func TestPrepareStatement(t *testing.T) {
	tests := []PrepareStatementTest{
		{
			UnpreparedStatement:        "SELECT * FROM table WHERE col1 = 1",
			ExpectedStatement:          "SELECT * FROM table WHERE col1 = 1",
			ExpectedParameterPositions: nil,
			Name:                       "No Parameter",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @name",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 = $1",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"name": {0},
				},
				totalPositions: 1,
			},
			Name: "Single Parameter",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @name AND col2 = @occupation",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"name":       {0},
					"occupation": {1},
				},
				totalPositions: 2,
			},
			Name: "Two Parameters",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @name AND col2 = @name",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"name": {0, 1},
				},
				totalPositions: 2,
			},
			Name: "Repeated Named Parameter",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 IN (@something, @else)",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 IN ($1, $2)",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"something": {0},
					"else":      {1},
				},
				totalPositions: 2,
			},
			Name: "Parameters In Parenthesis",
		},
		{
			UnpreparedStatement:        "SELECT * FROM table WHERE col1 = '@literal' AND col2 LIKE '@literal'",
			ExpectedStatement:          "SELECT * FROM table WHERE col1 = '@literal' AND col2 LIKE '@literal'",
			ExpectedParameterPositions: nil,
			Name:                       "Escaped Parameters",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = '@literal' AND col2 = @literal AND col3 LIKE '@literal'",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 = '@literal' AND col2 = $1 AND col3 LIKE '@literal'",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"literal": {0},
				},
				totalPositions: 1,
			},
			Name: "Escaping and non-escaping parameters",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @foo AND col2 IN (SELECT id FROM tabl2 WHERE col10 = @bar)",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 = $1 AND col2 IN (SELECT id FROM tabl2 WHERE col10 = $2)",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"foo": {0},
					"bar": {1},
				},
				totalPositions: 2,
			},
			Name: "Parameters In Subclause",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @1234567890 AND col2 = @0987654321",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"1234567890": {0},
					"0987654321": {1},
				},
				totalPositions: 2,
			},
			Name: "Numeric Parameters",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 = $1",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"ABCDEFGHIJKLMNOPQRSTUVWXYZ": {0},
				},
				totalPositions: 1,
			},
			Name: "Upcased Parameter",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @abc123ABC098",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 = $1",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"abc123ABC098": {0},
				},
				totalPositions: 1,
			},
			Name: "Multicased alphanumeric parameter",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 LIKE %@t%",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 LIKE %$1%",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"t": {0},
				},
				totalPositions: 1,
			},
			Name: "Pattern Matching",
		},
		{
			UnpreparedStatement: "ST_GeomFromText('POINT(' || @long @lat || ',4326)'",
			ExpectedStatement:   "ST_GeomFromText('POINT(' || $1 $2 || ',4326)'",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"long": {0},
					"lat":  {1},
				},
				totalPositions: 2,
			},
			Name: "Concated parameters",
		},
		{
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @first_name AND col2 = @last_name",
			ExpectedStatement:   "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameterPositions: &NamedParameterPositions{
				parameterPositions: map[string][]int{
					"first_name": {0},
					"last_name":  {1},
				},
				totalPositions: 2,
			},
			Name: "Snake Case",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			preparedStatement, err := PrepareStatement(test.UnpreparedStatement)
			require.NoError(t, err)
			require.Equal(t, test.ExpectedStatement, preparedStatement.Revised())
			actualParamPositions := preparedStatement.NamedParameterPositions()
			require.Equal(t, test.ExpectedParameterPositions, actualParamPositions, test.Name)
		})
	}
}

type BoundNamedParameterValuesTest struct {
	Name                     string
	UnpreparedStatement      string
	BindNamedParameterValues []TestNamedParameterValue
	ExpectedBoundValues      []any
}

type TestNamedParameterValue struct {
	Name  string
	Value any
}

func TestBindNamedParameterValues(t *testing.T) {
	tests := []BoundNamedParameterValuesTest{
		{
			Name:                "Single String Parameter",
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @foo",
			BindNamedParameterValues: []TestNamedParameterValue{
				{
					Name:  "foo",
					Value: "bar",
				},
			},
			ExpectedBoundValues: []any{
				"bar",
			},
		},
		{
			Name:                "Two String Parameter",
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @foo AND col2 = @foo2",
			BindNamedParameterValues: []TestNamedParameterValue{
				{
					Name:  "foo",
					Value: "bar",
				},
				{
					Name:  "foo2",
					Value: "bart",
				},
			},
			ExpectedBoundValues: []any{
				"bar",
				"bart",
			},
		},
		{
			Name:                "Repeated Parameters",
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @foo AND col2 = @foo",
			BindNamedParameterValues: []TestNamedParameterValue{
				{
					Name:  "foo",
					Value: "bar",
				},
			},
			ExpectedBoundValues: []any{
				"bar",
				"bar",
			},
		},
		{
			Name:                "Type Parameters",
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @str AND col2 = @int AND col3 = @pi",
			BindNamedParameterValues: []TestNamedParameterValue{
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
			ExpectedBoundValues: []any{
				"foo",
				1,
				3.14,
			},
		},
		{
			Name:                "Ordered Parameters",
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @foo AND col2 = @bar AND col3 = @foo AND col4 = @foo AND col5 = @bar",
			BindNamedParameterValues: []TestNamedParameterValue{
				{
					Name:  "foo",
					Value: "something",
				},
				{
					Name:  "bar",
					Value: "else",
				},
			},
			ExpectedBoundValues: []any{
				"something", "else", "something", "something", "else",
			},
		},
		{
			Name:                "Case Sensitive",
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @foo AND col2 = @FOO",
			BindNamedParameterValues: []TestNamedParameterValue{
				{
					Name:  "foo",
					Value: "baz",
				},
				{
					Name:  "FOO",
					Value: "quux",
				},
			},
			ExpectedBoundValues: []any{
				"baz", "quux",
			},
		},
		{
			Name:                "Nil Parameter",
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @foo",
			BindNamedParameterValues: []TestNamedParameterValue{
				{
					Name:  "foo",
					Value: pq.Array([]string{}),
				},
			},
			ExpectedBoundValues: []any{
				pq.Array([]string{}),
			},
		},
		{
			Name:                "Casted Type Parameter",
			UnpreparedStatement: "SELECT * FROM table WHERE col1 = @foo",
			BindNamedParameterValues: []TestNamedParameterValue{
				{
					Name:  "foo",
					Value: "'testing'::varchar",
				},
			},
			ExpectedBoundValues: []any{
				"'testing'::varchar",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			parameterFuncs := make([]BindNamedParameterValueFunc, 0)
			prepraredStatement, err := PrepareStatement(test.UnpreparedStatement)
			require.NoError(t, err)

			for _, bindValue := range test.BindNamedParameterValues {
				err = prepraredStatement.BindNamedParameterValue(bindValue.Name, bindValue.Value)
				require.NoError(t, err)
				parameterFuncs = append(parameterFuncs, BindNamedParameterValue(bindValue.Name, bindValue.Value))
			}

			err = prepraredStatement.BindNamedParameterValues(parameterFuncs...)
			require.NoError(t, err)

			for posIndex, boundValue := range prepraredStatement.BoundNamedParameterValues() {
				require.Equal(t, boundValue, test.ExpectedBoundValues[posIndex], test.Name)
			}
		})
	}
}
