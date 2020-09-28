package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase_SetID(t *testing.T) {
	type Thing struct {
		Base
		Field string
	}

	thing := new(Thing)
	assert.Equal(t, thing.ID, "")

	thing.SetID("omg")
	assert.Equal(t, thing.ID, "omg")
}
