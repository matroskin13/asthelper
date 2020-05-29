package asthelper

import (
	"bytes"
	"github.com/go-toolsmith/astcopy"
	"go/ast"
	"go/format"
	"go/token"
)

type Method struct {
	Decl *ast.FuncDecl
}

func NewMethod(name string) *Method {
	return &Method{
		Decl: &ast.FuncDecl{
			Name: ast.NewIdent(name),
			Type: &ast.FuncType{
				Params:  &ast.FieldList{},
				Results: &ast.FieldList{},
			},
		},
	}
}

func MethodFromAst(decl *ast.FuncDecl) *Method {
	return &Method{Decl: decl}
}

func (m *Method) IsStructMethod() bool {
	return m.Decl.Recv != nil && len(m.Decl.Recv.List) > 0
}

func (m *Method) SetName(name string) *Method {
	m.Decl.Name = ast.NewIdent(name)

	return m
}

func (m *Method) Copy() *Method {
	return &Method{Decl: astcopy.FuncDecl(m.Decl)}
}

func (m *Method) Body() *Body {
	return &Body{ast: m.Decl.Body}
}

func (m *Method) GenerateCode() (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), m.Decl); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (m *Method) ClearParams() *Method {
	m.Decl.Type.Params.List = []*ast.Field{}

	return m
}

func (m *Method) ClearResult() *Method {
	m.Decl.Type.Results.List = []*ast.Field{}

	return m
}

func (m *Method) AddParam(name, paramType string, isPointer bool) *Method {
	var field ast.Field

	if name != "" {
		field.Names = []*ast.Ident{
			{Name: name},
		}
	}
	if isPointer {
		field.Type = &ast.StarExpr{X: &ast.BasicLit{Value: paramType}}
	} else {
		field.Type = &ast.Ident{Name: paramType}
	}

	m.Decl.Type.Params.List = append(m.Decl.Type.Params.List, &field)

	return m
}

func (m *Method) AddResult(name, paramType string, isPointer bool) *Method {
	var field ast.Field

	if name != "" {
		field.Names = []*ast.Ident{ast.NewIdent(name)}
	}
	if isPointer {
		field.Type = &ast.StarExpr{X: &ast.BasicLit{Value: paramType}}
	} else {
		field.Type = &ast.BasicLit{Value: paramType}
	}

	m.Decl.Type.Results.List = append(m.Decl.Type.Results.List, &field)

	return m
}

func (m *Method) AddReceive(name, forType string, isPointer bool) *Method {
	if m.Decl.Recv == nil {
		m.Decl.Recv = &ast.FieldList{}
	}

	var field ast.Field

	field.Names = []*ast.Ident{ast.NewIdent(name)}
	if isPointer {
		field.Type = &ast.StarExpr{
			X: &ast.BasicLit{Value: forType},
		}
	} else {
		field.Type = &ast.BasicLit{Value: forType}
	}

	m.Decl.Recv.List = append(m.Decl.Recv.List, &field)

	return m
}

func (m *Method) AddAssigment(assigment *Assigment) *Method {
	if m.Decl.Body == nil {
		m.Decl.Body = &ast.BlockStmt{}
	}

	m.Decl.Body.List = append(m.Decl.Body.List, assigment.Ast)

	return m
}
