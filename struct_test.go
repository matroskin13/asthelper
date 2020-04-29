package asthelper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStruct(t *testing.T) {
	s := NewStruct("User")

	code, err := s.GenerateCode()

	assert.Nil(t, err)
	assert.Equal(t, `type User struct {
}`, code)
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
