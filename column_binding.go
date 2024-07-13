package dbsql

// // DefineColumnBinding ...
// func DefineColumnBinding(column Column, binderFunc ColumnBinderFunc) ColumnBinder {
// 	return &ColumnBinding{
// 		column:     column,
// 		binderFunc: binderFunc,
// 	}
// }
//
// // ColumnBinding ...
// type ColumnBinding struct {
// 	column     Column
// 	binderFunc func(bindFunc func(value any) error) ColumnBinderFunc
// }
//
// func (c ColumnBinding) Column() Column {
// 	return c.column
// }
//
// func (c ColumnBinding) BinderFunc() ColumnBinderFunc {
// 	return c.binderFunc
// }
//
// type ColumnBinder[T any] interface {
// 	Column() Column
// 	BinderFunc() func[T any](bindFunc func(value T) error) ColumnBinderFunc
// }
//
// // ColumnBindings ...
// type (
// 	ColumnBindings []ColumnBinding
// 	ColumnBinders  []ColumnBinder
// )
//
// // Count returns the number of column bindings.
// func (c ColumnBindings) Count() int {
// 	return len(c)
// }
//
// // Count returns the number of column bindings.
// func (c ColumnBinders) Count() int {
// 	return len(c)
// }
