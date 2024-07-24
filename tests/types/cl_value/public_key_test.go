package cl_value

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

func Test_PublicKeyHash_ToString(t *testing.T) {
	res, err := keypair.NewPublicKey("01387cda157981c8b9de8aa1cb1eeaadba0c10f24b207ae168f3a231d1523b835a")
	require.NoError(t, err)
	assert.Equal(t, "01387cda157981c8b9de8aa1cb1eeaadba0c10f24b207ae168f3a231d1523b835a", res.String())
}

func Test_PublicKey_FromBytesByType(t *testing.T) {
	source := "02037292af42f13f1f49507c44afe216b37013e79a062d7e62890f77b8adad60501e"
	decoded, err := hex.DecodeString(source)
	require.NoError(t, err)
	res, err := clvalue.FromBytesByType(decoded, cltype.PublicKey)
	require.NoError(t, err)
	assert.Equal(t, "02037292af42f13f1f49507c44afe216b37013e79a062d7e62890f77b8adad60501e", res.String())
}
