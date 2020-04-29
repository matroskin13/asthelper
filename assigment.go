package asthelper

import (
	"go/ast"
	"go/token"
)

type Assigment struct {
	Ast *ast.AssignStmt
}

func NewAssigment(left, right string, sign token.Token) (*Assigment, error) {
	leftExpr, err := ParsExpr(left)
	if err != nil {
		return nil, err
	}

	rightExpr, err := ParsExpr(right)
	if err != nil {
		return nil, err
	}

	return &Assigment{
		Ast: &ast.AssignStmt{
			Lhs: []ast.Expr{leftExpr.Ast()},
			Tok: sign,
			Rhs: []ast.Expr{rightExpr.Ast()},
		},
	}, nil
}

func NewAssigmentFromExpr(left *Expr, right *Expr, sign token.Token) *Assigment {
	return &Assigment{
		Ast: &ast.AssignStmt{
			Lhs: []ast.Expr{left.Ast()},
			Tok: sign,
			Rhs: []ast.Expr{right.Ast()},
		},
	}
}

func (a *Assigment) ReplaceValue(expr *Expr) *Assigment {
	a.Ast.Rhs = []ast.Expr{expr.Ast()}

	return nil
}
