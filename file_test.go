package asthelper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile_GetStruct(t *testing.T) {
	n := parseExampleFile()
	f := NewFileFromAst(n)

	user := f.GetStruct("User")

	assert.NotNil(t, user)
	assert.Equal(t, "User", user.GetName())
}

func TestFile_AddImport(t *testing.T) {
	f := NewFile("example")

	f.AddImport("github.com/matroskin13/granule", "granule")
	f.AddImport("github.com/matroskin13/granule", "granule1")

	code, err := f.GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `package example

import (
	granule "github.com/matroskin13/granule"
	granule1 "github.com/matroskin13/granule"
)
`, code)
}
