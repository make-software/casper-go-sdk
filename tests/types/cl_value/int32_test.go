package cl_value

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/types/clvalue"
)

func Test_NewInt32FromBuffer_maxValue(t *testing.T) {
	maxInBytes := clvalue.NewCLInt32(math.MaxInt32).Bytes()
	res := clvalue.NewInt32FromBytes(maxInBytes)
	assert.Equal(t, int32(math.MaxInt32), res.Value())
}

func Test_NewInt32FromBuffer_negative(t *testing.T) {
	maxInBytes := clvalue.NewCLUInt32(math.MaxUint32).Bytes()
	res := clvalue.NewInt32FromBytes(maxInBytes)
	assert.Equal(t, int32(-1), res.Value())
}
