package asthelper

import (
	"fmt"
	"go/ast"
)

type Value struct {
	ast *ast.ValueSpec
}

func ValueFromAst(n *ast.ValueSpec) *Value {
	return &Value{ast: n}
}

func (v *Value) ToArray() *ArrayValue {
	return &ArrayValue{ast: v.ast}
}

func (v *Value) AddValuesToArray(exprs ...*Expr) (*Value, error) {
	if arr, ok := v.ast.Values[0].(*ast.CompositeLit); ok {
		for _, exp := range exprs {
			arr.Elts = append(arr.Elts, exp.Ast())
		}

		return nil, nil
	}

	return nil, fmt.Errorf("value is not array")
}

type BaseValue struct {
	ast *ast.ValueSpec
}

func NewBaseValue(name string, value string) *BaseValue {
	b := &BaseValue{ast: &ast.ValueSpec{
		Doc:     nil,
		Names:   []*ast.Ident{ast.NewIdent(name)},
		Values:  nil,
		Comment: nil,
	}}

	b.SetValue(fmt.Sprintf("%q", value))

	return b
}

func (b *BaseValue) SetValue(value string) *BaseValue {
	b.ast.Values = []ast.Expr{&ast.BasicLit{Value: value}}

	return b
}

func (b *BaseValue) ToValue() *Value {
	return &Value{ast: b.ast}
}

type ArrayValue struct {
	ast *ast.ValueSpec
}

func (v *ArrayValue) AddValues(exprs ...*Expr) *ArrayValue {
	if arr, ok := v.ast.Values[0].(*ast.CompositeLit); ok {
		for _, exp := range exprs {
			arr.Elts = append(arr.Elts, exp.Ast())
		}
	}

	return v
}

func (v *ArrayValue) ClearValues() *ArrayValue {
	if arr, ok := v.ast.Values[0].(*ast.CompositeLit); ok {
		arr.Elts = []ast.Expr{}
	}

	return v
}
