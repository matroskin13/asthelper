package asthelper

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

var exampleWithBody = `package example

type NotUser struct {
	NotId int
}

func (u *NotUser) GetName() string {
	return "superuser"
}

type User struct {
	Id int
}

func (u *User) GetName() string {
	return "superuser"
}

func Example(rows *pgx.Rows) {
	rows.Scan(&model.A, &model.B)
	rows.Scann(&model.A, &model.B)
	rowss.Scan(&model.A, &model.B)
}`

func parseExampleFile() *ast.File {
	fset := token.NewFileSet()
	n, err := parser.ParseFile(fset, "example.go", exampleWithBody, 0)

	if err != nil {
		panic(err)
	}

	return n
}

func TestBody_SearchVarMethodCall(t *testing.T) {
	f := parseExampleFile()

	funDecl := f.Scope.Lookup("Example").Decl.(*ast.FuncDecl)
	methodBody := MethodFromAst(funDecl).Body()

	calls := methodBody.SearchVarMethodCall("rows", "Scan")

	assert.Len(t, calls, 1)
}
