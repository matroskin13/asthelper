package asthelper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInterface(t *testing.T) {
	i := NewInterface("IExample")

	code, err := i.GenerateCode()
	assert.Nil(t, err)
	assert.Equal(t, `type IExample interface {
}`, code)
}

func TestInterface_AddMethod(t *testing.T) {
	i := NewInterface("IExample")

	m := NewMethod("")
	m.AddParam("", "string", false)
	m.AddParam("", "string", false)
	m.AddResult("", "int", false)
	m.AddResult("", "error", false)

	i.AddMethod("example", m)

	code, err := i.GenerateCode()
	assert.Nil(t, err)
	assert.Equal(t, `type IExample interface {
	example(string, string) (int, error)
}`, code)
}
