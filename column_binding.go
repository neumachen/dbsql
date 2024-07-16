package dbsql

import (
	"errors"

	"github.com/neumachen/dbsql/internal"
)

// DefineColumnBinding creates a new ColumnBinder with the given Column, ColumnBinderFunc,
// and optional DefineColumnBindingRulesFunc. It applies the default binding rules and
// any additional rules specified by the DefineColumnBindingRulesFunc.
//
// Example:
//
//	column := NewColumn("id", "INT")
//	binderFunc := func(cb ColumnBinder, mr MappedRow) error {
//		// Binding logic here
//		return nil
//	}
//	customRuleFunc := func(rules *ColumnBindingRules) {
//		rules.AllowNull = true
//	}
//	binder := DefineColumnBinding(column, binderFunc, customRuleFunc)
func DefineColumnBinding(
	column Column,
	binderFunc ColumnBinderFunc,
	defineRulesFuncs ...DefineColumnBindingRulesFunc,
) ColumnBinder {
	bindingRules := defaultColumnBindingRules
	for i := range defineRulesFuncs {
		defineRulesFuncs[i](bindingRules)
	}

	return &columnBinding{
		column:       column,
		binderFunc:   binderFunc,
		bindingRules: bindingRules,
	}
}

type columnBinding struct {
	column       Column
	binderFunc   ColumnBinderFunc
	bindingRules ColumnBindingRules
}

func (c columnBinding) Column() Column {
	return c.column
}

func (c columnBinding) BindColumn(mappedRow MappedRow) error {
	if internal.IsNil(c.binderFunc) {
		return errors.New("column binding has nil binder func")
	}
	return c.binderFunc(mappedRow, c)
}

func (c columnBinding) BindingRules() ColumnBindingRules {
	return c.bindingRules
}
