package asthelper

import (
	"github.com/stretchr/testify/assert"
	"go/parser"
	"go/token"
	"testing"
)

var exampleWithStructs = `package example

type User struct {
	Name string
}

func (u *User) GetName() string {
	return u.Name
} 

func (u *User) SetName(name string) {
	u.Name = name
}

func (u User) NotLink() {}

type Car struct {
	UserId int
}

func (c *Car) Open() error {
	return nil
}
`

func TestNewStruct(t *testing.T) {
	s := NewStruct("User")

	code, err := s.GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `type User struct {
}`, code)
}

func TestStruct_GetMethods(t *testing.T) {
	fset := token.NewFileSet()
	n, err := parser.ParseFile(fset, "example.go", exampleWithStructs, 0)

	if err != nil {
		panic(err)
	}

	f := NewFileFromAst(n)
	user := f.GetStruct("User")

	assert.Equal(t, "User", user.GetName())
	assert.Equal(t, 3, len(user.GetMethods()))
}

func TestStruct_AddField(t *testing.T) {
	s := NewStruct("User")

	s.AddField("Id", "int", false)
	s.AddField("Email", "string", true)

	code, err := s.GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `type User struct {
	Id    int
	Email *string
}`, code)
}

func TestStruct_ClearField(t *testing.T) {
	s := NewStruct("User")

	s.AddField("Id", "int", false)
	s.AddField("Email", "string", true)

	s.ClearField()

	code, err := s.GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `type User struct {
}`, code)
}

func TestStruct_GetName(t *testing.T) {
	s := NewStruct("User")

	assert.Equal(t, "User", s.GetName())
}
