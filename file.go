package asthelper

import (
	"bytes"
	"fmt"
	"github.com/go-toolsmith/astcopy"
	"go/ast"
	"go/format"
	"go/token"
)

type File struct {
	ast *ast.File
}

func NewFileFromAst(f *ast.File) *File {
	return &File{ast: f}
}

func NewFile(packageName string) *File {
	return &File{
		ast: &ast.File{
			Name: ast.NewIdent(packageName),
		},
	}
}

func (f *File) CloneImports(from *File) *File {
	fromImportDecl := astcopy.GenDecl(from.getImportDecl())

	for _, spec := range fromImportDecl.Specs {
		if imp, ok := spec.(*ast.ImportSpec); ok {
			var name string

			if imp.Name != nil {
				name = imp.Name.Name
			}

			f.addImport(imp.Path.Value, name)
		}
	}

	return f
}

func (f *File) Copy() *File {
	copyAst := astcopy.File(f.ast)

	return &File{ast: copyAst}
}

func (f *File) Ast() *ast.File {
	return f.ast
}

func (f *File) GetVar(name string) *Value {
	for _, decl := range f.ast.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if value, ok := genDecl.Specs[0].(*ast.ValueSpec); ok {
				if value.Names[0].Name == name {
					return ValueFromAst(value)
				}
			}
		}
	}

	return nil
}

func (f *File) AddVar(name string, value string) *File {
	f.ast.Decls = append(f.ast.Decls, &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Doc:   nil,
				Names: []*ast.Ident{ast.NewIdent(name)},
				Values: []ast.Expr{
					&ast.BasicLit{
						ValuePos: 0,
						Kind:     token.STRING,
						Value:    value,
					},
				},
				Comment: nil,
			},
		},
	})

	return f
}

func (f *File) getImportDecl() *ast.GenDecl {
	var importDecl *ast.GenDecl

	for _, decl := range f.ast.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			importDecl = genDecl
			break
		}
	}

	if importDecl == nil {
		importDecl = &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: []ast.Spec{},
		}

		f.ast.Decls = append(f.ast.Decls, importDecl)
	}

	return importDecl
}

func (f *File) AddStruct(st *Struct) *File {
	f.ast.Decls = append(f.ast.Decls, st.ast)

	for _, method := range st.methods {
		f.ast.Decls = append(f.ast.Decls, method)
	}

	return f
}

func (f *File) AddInterface(i *Interface) *File {
	f.ast.Decls = append(f.ast.Decls, i.ast)

	return f
}

func (f *File) AddImport(path string, alias string) *File {
	return f.addImport(fmt.Sprintf("%q", path), alias)
}

func (f *File) addImport(path string, alias string) *File {
	if path == fmt.Sprintf("%q", "github.com/matroskin13/granule/generator/entity") {
		return f
	}

	for _, decl := range f.ast.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			for _, spec := range genDecl.Specs {
				if importSpec, ok := spec.(*ast.ImportSpec); ok {
					if importSpec.Name.Name == alias && importSpec.Path.Value == path {
						return f
					}
				}
			}
		}
	}

	importDecl := f.getImportDecl()

	importDecl.Specs = append(importDecl.Specs, &ast.ImportSpec{
		Name: ast.NewIdent(alias),
		Path: &ast.BasicLit{Value: path},
	})

	return f
}

func (f *File) GenerateCode() (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), f.ast); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (f *File) GetMethod(name string) *Method {
	for _, decl := range f.ast.Decls {
		if funDecl, ok := decl.(*ast.FuncDecl); ok {
			if funDecl.Recv == nil || len(funDecl.Recv.List) == 0 {
				if funDecl.Name.Name == name {
					return MethodFromAst(funDecl)
				}
			}
		}
	}

	return nil
}

func (f *File) DeleteMethod(name string) *File {
	for i, decl := range f.ast.Decls {
		if funDecl, ok := decl.(*ast.FuncDecl); ok {
			if funDecl.Name.Name == name {
				f.ast.Decls = append(f.ast.Decls[:i], f.ast.Decls[i+1:]...)
			}
		}
	}

	return f
}

func (f *File) GetType(name string) *Type {
	for _, decl := range f.ast.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec); ok {
				if typeSpec.Name.Name == name {
					return &Type{ast: typeSpec}
				}
			}
		}
	}

	return nil
}

func (f *File) AddType(t *Type) *File {
	f.ast.Decls = append(f.ast.Decls, &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{t.ast},
	})

	return f
}

func (f *File) AddMethod(method *Method) *File {
	f.ast.Decls = append(f.ast.Decls, method.Decl)

	return f
}

func (f *File) GetStruct(name string) *Struct {
	var methods []*ast.FuncDecl
	var s *Struct

	for _, decl := range f.ast.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec); ok {
				if typeSpec.Name.Name == name {
					s = NewStructFromAst(genDecl)
				}
			}
		}

		if funDecl, ok := decl.(*ast.FuncDecl); ok {
			if funDecl.Recv != nil {
				for _, field := range funDecl.Recv.List {
					switch fieldType := field.Type.(type) {
					case *ast.StarExpr:
						if ident, ok := fieldType.X.(*ast.Ident); ok {
							if ident.Name == name {
								funDecl.Doc = nil
								methods = append(methods, funDecl)
							}
						}
					case *ast.Ident:
						if fieldType.Name == name {
							methods = append(methods, funDecl)
						}
					}
				}
			}
		}
	}

	if s != nil {
		for _, method := range methods {
			s.AddAstMethod(method)
		}
	}

	return s
}
