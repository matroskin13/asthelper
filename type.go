package asthelper

import "go/ast"

type Type struct {
	ast *ast.TypeSpec
}

func TypeFromAst(n *ast.TypeSpec) *Type {
	return &Type{ast: n}
}
