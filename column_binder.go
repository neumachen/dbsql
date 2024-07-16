package dbsql

import (
	"errors"
	"fmt"

	"github.com/neumachen/dbsql/internal"
)

// ColumnBinder is an interface that represents a column binder.
// It provides methods to access the column, bind the column, and retrieve binding rules.
type ColumnBinder interface {
	// Column returns the Column associated with this ColumnBinder.
	Column() Column

	// BindColumn binds the column to the provided MappedRow.
	// It returns an error if the binding operation fails.
	BindColumn(mappedRow MappedRow) error

	// BindingRules returns the ColumnBindingRules associated with this ColumnBinder.
	BindingRules() ColumnBindingRules
}

// ColumnBindersFilterFunc is a function type used for filtering ColumnBinders.
type ColumnBindersFilterFunc func(i int, columnBinders ColumnBinders) (bool, error)

// ColumnBinders is a slice of ColumnBinder.
type ColumnBinders []ColumnBinder

// DefineColumnBinders creates a new ColumnBinders from the given ColumnBinder slice.
//
// Example:
//
//	binder1 := DefineColumnBinding(column1, binderFunc1)
//	binder2 := DefineColumnBinding(column2, binderFunc2)
//	binders := DefineColumnBinders(binder1, binder2)
func DefineColumnBinders(columnBinder ...ColumnBinder) ColumnBinders {
	columnBinders := make(ColumnBinders, len(columnBinder))
	copy(columnBinders, columnBinder)
	return columnBinders
}

// FilterUsingFunc filters the ColumnBinders using the provided ColumnBindersFilterFunc.
// It returns a new ColumnBinders containing only the elements for which the filter function returns true.
//
// Example:
//
//	binders := DefineColumnBinders(binder1, binder2, binder3)
//	filteredBinders, err := binders.FilterUsingFunc(func(i int, cb ColumnBinders) (bool, error) {
//		return cb[i].Column().Name() == "id", nil
//	})
//	if err != nil {
//		// Handle error
//	}
//	// filteredBinders now contains only the binders for columns named "id"
func (c ColumnBinders) FilterUsingFunc(
	filterFunc ColumnBindersFilterFunc,
) (
	ColumnBinders,
	error,
) {
	copied := 0
	dupColumnBinders := make(ColumnBinders, len(c))
	copy(dupColumnBinders, c)

	for i := 0; i < len(dupColumnBinders); i++ {
		b, err := filterFunc(i, dupColumnBinders)
		if err != nil {
			return nil, err
		}
		if b {
			dupColumnBinders[copied] = dupColumnBinders[i]
			copied++
		}
	}
	if copied < 1 {
		return c, nil
	}
	for i := copied; i < len(dupColumnBinders); i++ {
		dupColumnBinders[i] = nil
	}

	return dupColumnBinders[:copied], nil
}

// ColumnBinderFunc is a function type that defines the signature for binding a column to a field.
// It takes a MappedRow and a ColumnBinder as input and returns an error if the binding fails.
type ColumnBinderFunc func(mappedRow MappedRow, columnBinder ColumnBinder) error

// BindColumnToField is a generic function that creates a ColumnBinderFunc for binding a specific type T.
// It takes a bind function as input and returns a ColumnBinderFunc that handles the binding process.
//
// Parameters:
//   - bindFunc: A function that takes a value of type T and returns an error if binding fails.
//
// Returns:
//   - A ColumnBinderFunc that can be used to bind a column to a field of type T.
//
// Example usage:
//
//	intBinder := BindColumnToField(func(value int) error {
//	    // Bind the int value to a field
//	    return nil
//	})
func BindColumnToField[T any](bindFunc func(value T) error) ColumnBinderFunc {
	return func(mappedRow MappedRow, columnBinder ColumnBinder) error {
		if internal.IsNil(columnBinder) {
			return errors.New("column binder is nil")
		}

		value, found := mappedRow.Get(columnBinder.Column())
		if !found {
			if columnBinder.BindingRules().RequiredColumn() {
				return fmt.Errorf("required column %s not found", columnBinder.Column().String())
			}

			return nil
		}

		if internal.IsNilOrZeroValue(value) {
			return nil
		}

		typedValue, ok := value.(T)
		if !ok {
			return fmt.Errorf(
				"column %s has a type of %T and does not match asserted type: %T",
				columnBinder.Column().String(),
				value,
				*new(T),
			)
		}

		if internal.IsNilOrZeroValue(typedValue) {
			return nil
		}

		return bindFunc(typedValue)
	}
}
