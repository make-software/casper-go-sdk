package cl_value

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

func Test_BoolValueBytesParser_ToBytes(t *testing.T) {
	tests := []struct {
		name     string
		val      clvalue.CLValue
		excepted []byte
	}{
		{
			"true val to bytes",
			clvalue.NewCLBool(true),
			[]byte{1},
		},
		{
			"false val to bytes",
			clvalue.NewCLBool(false),
			[]byte{0},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.excepted, test.val.Bytes())
		})
	}
}

func Test_BoolValueBytesParser_FomBytes(t *testing.T) {
	tests := []struct {
		name     string
		val      []byte
		excepted clvalue.Bool
	}{
		{
			"true val to bytes",
			[]byte{1},
			clvalue.Bool(true),
		},
		{
			"false val to bytes",
			[]byte{0},
			clvalue.Bool(false),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer(test.val)
			res, err := clvalue.NewBoolFromBuffer(buf)
			require.NoError(t, err)
			assert.Equal(t, test.excepted, *res)

		})
	}
}
