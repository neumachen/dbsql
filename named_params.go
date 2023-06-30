package sqlstmt

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

const (
	// A constant rune (':') used as the prefix for named parameters.
	parameterPrefix = ':'
	// A constant rune ('\'') used as the escape character for quotes in the query.
	parameterEscape = '\''
	// A constant string ("$") used as the prefix for positional parameter placeholders.
	placeholderPrefix = "$"
)

// ConvertNamedToPositionalParams is a function that parses a query statement and converts named parameters
// to positional parameters. It takes a query statement as input, represented as a byte slice, and returns a Statement
// interface and an error.
func ConvertNamedToPositionalParams(queryStatement []byte) (Statement, error) {
	var revisedQuery []byte
	var namedParameter []byte

	stmt := &statement{}

	runeCount := utf8.RuneCount(queryStatement)

	var character rune
	var size int
	var positionIndex int

	for i := 0; i < runeCount; {
		character, size = utf8.DecodeRune(queryStatement[i:])
		i += size

		if character == parameterPrefix {
			// Collect the characters after the parameter prefix until a non-content rune is encountered.
			for {
				character, size = utf8.DecodeRune(queryStatement[i:])
				i += size

				if isNonContentRune(character, size) {
					break
				}

				namedParameter = append(namedParameter, string(character)...)
			}

			stmt.setPosition(string(namedParameter), positionIndex)
			positionIndex++

			// Replace the named parameter with a positional parameter placeholder.
			placeholder := strconv.Itoa(positionIndex)
			revisedQuery = append(revisedQuery, placeholderPrefix...)
			revisedQuery = append(revisedQuery, placeholder...)

			namedParameter = namedParameter[:0] // Reset the parameterBuilder

			if isEmptyRune(character, size) {
				break
			}
		}

		// Append the character to the revised query.
		revisedQuery = append(revisedQuery, byte(character))

		// If it's a quote, continue appending to the builder but do not search for parameters.
		if character == parameterEscape {
			// Append characters until the closing quote is encountered.
			for {
				character, size = utf8.DecodeRune(queryStatement[i:])
				i += size
				revisedQuery = append(revisedQuery, byte(character))

				if character == parameterEscape {
					break
				}
			}
		}
	}

	stmt.revisedQuery = revisedQuery
	stmt.parameters = make(PositionalParameters, positionIndex)

	return stmt, nil
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
