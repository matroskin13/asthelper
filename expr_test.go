package asthelper

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestParsExpr(t *testing.T) {
	exp, err := ParsExpr("user.Selector")

	assert.Nil(t, err)
	assert.IsType(t, &ast.SelectorExpr{}, exp.Ast())
}
