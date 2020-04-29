package asthelper

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

type Call struct {
	ast *ast.CallExpr
}

func NewCall(what *Expr) *Call {
	return &Call{
		ast: &ast.CallExpr{
			Fun:  what.ast,
			Args: []ast.Expr{},
		},
	}
}

func (c *Call) ToExpr() *Expr {
	return ExprFromAst(c.ast)
}

func (c *Call) GenerateCode() (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), c.ast); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *Call) Ast() *ast.CallExpr {
	return c.ast
}

func (c *Call) ClearArguments() *Call {
	c.ast.Args = []ast.Expr{}

	return c
}

func (c *Call) AddArgument(expr *Expr) *Call {
	c.ast.Args = append(c.ast.Args, expr.Ast())

	return c
}

func (c *Call) AddArguments(expr []*Expr) *Call {
	for _, exp := range expr {
		c.ast.Args = append(c.ast.Args, exp.Ast())
	}

	return c
}
