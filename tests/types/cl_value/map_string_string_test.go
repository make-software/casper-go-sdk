package cl_value

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

func Test_MapType_Bool_To_String(t *testing.T) {
	newtype := cltype.Map{Key: cltype.Bool, Val: &cltype.Map{
		Key: cltype.Bool,
		Val: cltype.Bool,
	}}
	res := newtype.String()
	assert.Equal(t, "Map (Bool: Map (Bool: Bool))", res)
}

func Test_MapType_String_To_String_1(t *testing.T) {
	newtype := cltype.Map{Key: cltype.String, Val: cltype.Bool}
	assert.Equal(t, "Map (String: Bool)", newtype.String())
}
