package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_ParseResultArgument_shouldParseOk(t *testing.T) {
	source := `{"cl_type":{"Result":{"ok":"String","err":"String"}},"bytes":"010a00000068656c6c6f776f726c64","parsed":{"Ok":"helloworld"}}`
	var arg types.Argument
	err := json.Unmarshal([]byte(source), &arg)
	require.NoError(t, err)
	val, err := arg.Value()
	require.NoError(t, err)
	assert.True(t, val.Result.IsSuccess)
	assert.Equal(t, "helloworld", val.Result.Value().String())
}
