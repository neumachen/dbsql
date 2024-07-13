package dbsql

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

const (
	// A constant rune (':') used as the prefix for named parameters.
	parameterPrefix = '@'
	// A constant rune ('\'') used as the escape character for quotes in the query.
	parameterEscape = '\''
	// A constant string ("$") used as the prefix for positional parameter placeholders.
	// TODO: add support for other drivers, currently only supports PostgreSQL.
	pgPlaceholderPrefix = "$"
)

// PrepareStatement takes an unprepared SQL statement and returns a PreparedStatement interface.
// The PreparedStatement interface provides methods for managing named parameters, binding parameter
// values, and executing the prepared statement.
//
// The PrepareStatement function does the following:
//
// 1. Replaces named parameters (indicated by a '@' prefix) with positional placeholders ($1, $2, etc.).
// 2. Stores the positions of the named parameters in a NamedParameterPositions struct.
// 3. Returns a preparedStatement struct that implements the PreparedStatement interface.
//
// Example usage:
//
//	preparedStmt, err := PrepareStatement("SELECT * FROM users WHERE name = @name AND age > @age")
//	if err != nil {
//		// handle error
//	}
//
//	rows, err := preparedStmt.Query(db,
//		BindNamedParameterValue("name", "John"),
//		BindNamedParameterValue("age", 30),
//	)
//	if err != nil {
//		// handle error
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		// process rows
//	}
func PrepareStatement(unpreparedStatement string) (PreparedStatement, error) {
	var revisedStatement []byte
	var namedParameter []byte
	unpreparedStatementByte := []byte(unpreparedStatement)

	namedParamPositions := ParameterPositions{}

	runeCount := utf8.RuneCount(unpreparedStatementByte)

	var character rune
	var size int
	var positionIndex int

	// Iterate through the unprepared statement, replacing named parameters with positional placeholders
	for i := 0; i < runeCount; {
		character, size = utf8.DecodeRune(unpreparedStatementByte[i:])
		i += size

		if character == parameterPrefix {
			// Collect the characters after the parameter prefix until a non-content rune is encountered
			for {
				character, size = utf8.DecodeRune(unpreparedStatementByte[i:])
				i += size

				if isNonContentRune(character, size) {
					break
				}

				namedParameter = append(namedParameter, string(character)...)
			}

			// Set the position of the named parameter in the NamedParameterPositions struct
			namedParamPositions.insert(string(namedParameter), positionIndex)
			positionIndex++

			// Replace the named parameter with a positional parameter placeholder
			placeholder := strconv.Itoa(positionIndex)
			revisedStatement = append(revisedStatement, pgPlaceholderPrefix...)
			revisedStatement = append(revisedStatement, placeholder...)

			namedParameter = namedParameter[:0] // Reset the parameterBuilder

			if isEmptyRune(character, size) {
				break
			}
		}

		// Append the character to the revised query
		revisedStatement = append(revisedStatement, byte(character))

		// If it's a quote, continue appending to the builder but do not search for parameters
		if character == parameterEscape {
			// Append characters until the closing quote is encountered
			for {
				character, size = utf8.DecodeRune(unpreparedStatementByte[i:])
				i += size
				revisedStatement = append(revisedStatement, byte(character))

				if character == parameterEscape {
					break
				}
			}
		}
	}

	// Return a new preparedStatement struct with the revised statement, named parameter positions, and other information
	return &preparedStatement{
		originalStatement:     unpreparedStatement,
		namedParamPositions:   &namedParamPositions,
		revisedStatement:      string(revisedStatement),
		boundNamedParamValues: make(BoundParameterValues, positionIndex),
	}, nil
}

// isEmptyRune is a helper function that checks if a rune is empty.
// It takes a rune (r) and its size (size) as input and returns a boolean value indicating whether the rune is empty.
func isEmptyRune(r rune, size int) bool {
	return r == utf8.RuneError && size == 0
}

// runeUnderscore a constant rune ('_') used to specify the underscore character, which is treated as punctuation in Unicode.
const runeUnderscore = '_'

// isNonContentRune is a helper function that checks if a rune is a non-content rune.
// It takes a rune (r) and its size (size) as input and returns a boolean value indicating whether the rune is a non-content rune.
func isNonContentRune(r rune, size int) bool {
	if unicode.IsSpace(r) || unicode.IsPunct(r) {
		return r != runeUnderscore
	}

	return isEmptyRune(r, size)
}
