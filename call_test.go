package asthelper

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestCall_ClearArguments(t *testing.T) {
	n := parseExampleFile()

	funDecl := n.Scope.Lookup("Example").Decl.(*ast.FuncDecl)
	methodBody := MethodFromAst(funDecl).Body()

	calls := methodBody.SearchVarMethodCall("rows", "Scan")

	calls[0].ClearArguments()

	codeResult, err := calls[0].GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `rows.Scan()`, codeResult)
}

func TestCall_AddArgument(t *testing.T) {
	n := parseExampleFile()

	funDecl := n.Scope.Lookup("Example").Decl.(*ast.FuncDecl)
	methodBody := MethodFromAst(funDecl).Body()

	calls := methodBody.SearchVarMethodCall("rows", "Scan")

	calls[0].ClearArguments()

	arg, _ := ParsExpr("&rabbit.A")
	calls[0].AddArgument(arg).AddArgument(arg)

	codeResult, err := calls[0].GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `rows.Scan(&rabbit.A, &rabbit.A)`, codeResult)
}
