package asthelper

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

type Interface struct {
	ast *ast.GenDecl

	methods []*ast.FuncDecl
}

func NewInterface(name string) *Interface {
	return &Interface{
		ast: &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(name),
					Type: &ast.InterfaceType{
						Methods: &ast.FieldList{},
					},
				},
			},
			Rparen: 0,
		},
	}
}

func (i *Interface) AddMethod(name string, method *Method) *Interface {
	if typeSpec, ok := i.ast.Specs[0].(*ast.TypeSpec); ok {
		if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
			interfaceType.Methods.List = append(interfaceType.Methods.List, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(name)},
				Type:  method.Decl.Type,
			})
		}
	}

	return i
}

func (i *Interface) GenerateCode() (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), i.ast); err != nil {
		return "", err
	}

	return buf.String(), nil
}
