package asthelper

import (
	"github.com/stretchr/testify/assert"
	"go/token"
	"testing"
)

func TestNewMethod(t *testing.T) {
	method := NewMethod("Example")

	code, err := method.GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `func Example()`, code)
}

func TestMethod_AddParam(t *testing.T) {
	method := NewMethod("Example")

	method.AddParam("first", "string", false)
	method.AddParam("second", "bool", true)

	code, err := method.GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `func Example(first string, second *bool)`, code)
}

func TestMethod_AddResult(t *testing.T) {
	method := NewMethod("Example")

	method.AddResult("first", "string", false)
	method.AddResult("second", "bool", true)

	code, err := method.GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `func Example() (first string, second *bool)`, code)
}

func TestMethod_AddAssigment(t *testing.T) {
	method := NewMethod("Example")

	assigment, err := NewAssigment("user.Selector", "user.Custom", token.ASSIGN)

	assert.Nil(t, err)

	code, err := method.AddAssigment(assigment).GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `func Example() {
	user.Selector = user.Custom
}`, code)
}

func TestMethod_AddAssigment_Multiple(t *testing.T) {
	method := NewMethod("Example")

	assigment, err := NewAssigment("user.Selector", "user.Custom", token.ASSIGN)

	assert.Nil(t, err)

	code, err := method.AddAssigment(assigment).AddAssigment(assigment).GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `func Example() {
	user.Selector = user.Custom
	user.Selector = user.Custom
}`, code)
}
