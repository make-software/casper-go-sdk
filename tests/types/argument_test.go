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

func Test_ParsePublicKeyArgument_shouldParseOk(t *testing.T) {
	source := `{"cl_type":{"List":"PublicKey"},"bytes":"03000000010e31a03ea026a8e375653573e0120c8cb96699e6c9721ae1ea98f896e6576ac30193b3800386aefe11648150f6779158f2c7e1233c8e9b423338eb71b93ae6c5a90203c90c0ee375abc85da81a982507d1f8258a380af2058b63c37abdb9c7045940f4","parsed":["010e31a03ea026a8e375653573e0120c8cb96699e6c9721ae1ea98f896e6576ac3","0193b3800386aefe11648150f6779158f2c7e1233c8e9b423338eb71b93ae6c5a9", "0203c90c0ee375abc85da81a982507d1f8258a380af2058b63c37abdb9c7045940f4"]}`
	var arg types.Argument
	err := json.Unmarshal([]byte(source), &arg)
	require.NoError(t, err)
	val, err := arg.Value()
	require.NoError(t, err)
	assert.Equal(t, "010e31a03ea026a8e375653573e0120c8cb96699e6c9721ae1ea98f896e6576ac3", val.List.Elements[0].PublicKey.String())
	assert.Equal(t, "0193b3800386aefe11648150f6779158f2c7e1233c8e9b423338eb71b93ae6c5a9", val.List.Elements[1].PublicKey.String())
	assert.Equal(t, "0203c90c0ee375abc85da81a982507d1f8258a380af2058b63c37abdb9c7045940f4", val.List.Elements[2].PublicKey.String())
}
