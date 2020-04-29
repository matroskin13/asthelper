package asthelper

import (
	"bytes"
	"github.com/go-toolsmith/astcopy"
	"go/ast"
	"go/format"
	"go/token"
)

type Struct struct {
	ast *ast.GenDecl

	methods []*ast.FuncDecl
}

func NewStruct(name string) *Struct {
	return &Struct{
		ast: &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(name),
					Type: &ast.StructType{
						Fields: &ast.FieldList{},
					},
				},
			},
			Rparen: 0,
		},
	}
}

func NewStructFromAst(genDecl *ast.GenDecl) *Struct {
	return &Struct{
		ast: genDecl,
	}
}

func (s *Struct) GetName() string {
	return s.ast.Specs[0].(*ast.TypeSpec).Name.Name
}

func (s *Struct) SetName(name string) *Struct {
	s.ast.Specs[0].(*ast.TypeSpec).Name = ast.NewIdent(name)

	for _, method := range s.methods {
		switch t := method.Recv.List[0].Type.(type) {
		case *ast.Ident:
			t.Name = name
		case *ast.StarExpr:
			t.X = ast.NewIdent(name)
		}
	}

	return s
}

func (s *Struct) AddAstMethod(method *ast.FuncDecl) {
	s.methods = append(s.methods, method)
}

func (s *Struct) ClearField() *Struct {
	astStruct := s.ast.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)

	astStruct.Fields.List = []*ast.Field{}

	return s
}

func (s *Struct) RemoveField(name string) {
	astStruct := s.ast.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)

	for i, field := range astStruct.Fields.List {
		if field.Names[0].Name == name {
			astStruct.Fields.List = append(astStruct.Fields.List[:i], astStruct.Fields.List[i+1:]...)
		}
	}
}

func (s *Struct) AddField(name string, paramType string, isPointer bool) *Struct {
	var field ast.Field

	field.Names = []*ast.Ident{
		{Name: name},
	}
	if isPointer {
		field.Type = &ast.StarExpr{X: &ast.BasicLit{Value: paramType}}
	} else {
		field.Type = &ast.Ident{Name: paramType}
	}

	astStruct := s.ast.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)

	astStruct.Fields.List = append(astStruct.Fields.List, &field)

	return s
}

func (s *Struct) GetMethods() []*Method {
	var methods []*Method

	for _, method := range s.methods {
		methods = append(methods, MethodFromAst(method))
	}

	return methods
}

func (s *Struct) GetMethod(name string) *Method {
	for _, method := range s.methods {
		if method.Name.Name == name {
			return MethodFromAst(method)
		}
	}

	return nil
}

func (s *Struct) GenerateCode() (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), s.ast); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *Struct) Copy() *Struct {
	newS := &Struct{
		ast: astcopy.GenDecl(s.ast),
	}

	for _, method := range s.methods {
		newS.methods = append(newS.methods, astcopy.FuncDecl(method))
	}

	return newS
}
