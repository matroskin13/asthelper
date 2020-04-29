package asthelper

import (
	"go/ast"
)

type Body struct {
	ast *ast.BlockStmt
}

func (b *Body) Ast() *ast.BlockStmt {
	return b.ast
}

func (b *Body) Clear() *Body {
	b.ast.List = []ast.Stmt{}

	return b
}

func (b *Body) SearchVarMethodCall(varName string, methodName string) []Call {
	var result []Call

	ast.Inspect(b.ast, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.CallExpr:
			if selector, ok := x.Fun.(*ast.SelectorExpr); ok {
				selVar := selector.Sel.Name
				identVar, ok := selector.X.(*ast.Ident)

				if ok && identVar.Name == varName && selVar == methodName {
					result = append(result, Call{ast: x})
				}

				return false
			}
		}

		return true
	})

	return result
}

func (b *Body) GetAssigment(varName string) *Assigment {
	var result *ast.AssignStmt

	ast.Inspect(b.ast, func(n ast.Node) bool {
		if x, ok := n.(*ast.AssignStmt); ok {
			if ident, ok := x.Lhs[0].(*ast.Ident); ok {
				if ident.Name == varName {
					result = x
				}
			}
		}

		return true
	})

	if result == nil {
		return nil
	}

	return &Assigment{Ast: result}
}

func (b *Body) Insert(st ast.Stmt) {
	b.ast.List = append(b.ast.List, st)
}

func (b *Body) InsertBodyBefore(inserted *Body) *Body {
	b.ast.List = append(inserted.ast.List, b.ast.List...)

	return b
}
