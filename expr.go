package asthelper

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Expr struct {
	ast ast.Expr
}

func (exp *Expr) Ast() ast.Expr {
	return exp.ast
}

func SelectorExpr(left, right string, isLink bool) *Expr {
	selector := &ast.SelectorExpr{
		X:   ast.NewIdent(left),
		Sel: ast.NewIdent(right),
	}

	if isLink {
		return &Expr{ast: &ast.UnaryExpr{Op: token.AND, X: selector}}
	}

	return &Expr{ast: selector}
}

func BasicExpr(value string) *Expr {
	return &Expr{
		ast: &ast.BasicLit{
			ValuePos: 0,
			Kind:     0,
			Value:    value,
		},
	}
}

func ExprFromAst(expr ast.Expr) *Expr {
	return &Expr{ast: expr}
}

func ParsExpr(expr string) (*Expr, error) {
	exprAst, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, err
	}

	return &Expr{ast: exprAst}, nil
}
