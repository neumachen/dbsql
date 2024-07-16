package dbsql

// ColumnBindingRules defines the interface for column binding rules.
// It provides methods to determine if a column is required.
// NOTE: This interface is subject to change as it is not yet finalized.
type ColumnBindingRules interface {
	// RequiredColumn returns true if the column is required, false otherwise.
	RequiredColumn() bool
}

// DefineColumnBindingRulesFunc is a function type used to define column binding rules.
// It takes a pointer to columnBindingRules as an argument and modifies it.
type DefineColumnBindingRulesFunc func(rules *columnBindingRules)

// ColumnRequired returns a DefineColumnBindingRulesFunc that sets the column as required.
// This function is used to create a rule that marks a column as required.
func ColumnRequired() DefineColumnBindingRulesFunc {
	return func(rules *columnBindingRules) {
		rules.RequireColumn(true)
	}
}

// preBindingRules contains the pre-binding rules for a column.
type preBindingRules struct {
	requireColumn bool // Indicates whether the column is required
}

// columnBindingRules implements the ColumnBindingRules interface.
// It contains the pre-binding rules for a column.
type columnBindingRules struct {
	preBinding *preBindingRules
}

// RequireColumn sets whether the column is required or not.
// If the preBinding is nil, it initializes it before setting the value.
func (c *columnBindingRules) RequireColumn(require bool) {
	if c.preBinding == nil {
		c.preBinding = &preBindingRules{}
	}
	c.preBinding.requireColumn = require
}

// RequiredColumn returns whether the column is required.
// If the preBinding is nil, it initializes it before returning the value.
func (c *columnBindingRules) RequiredColumn() bool {
	if c.preBinding == nil {
		c.preBinding = &preBindingRules{}
	}
	return c.preBinding.requireColumn
}

// defaultColumnBindingRules is the default set of column binding rules.
// By default, columns are required.
var defaultColumnBindingRules = &columnBindingRules{
	preBinding: &preBindingRules{
		requireColumn: true,
	},
}
