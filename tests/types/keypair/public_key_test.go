package keypair

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/casper"
)

func Test_PublicKey_ToAccountHash(t *testing.T) {
	publicKey, err := casper.NewPublicKey("015a372b0e230bf9393e2df0b3de857bb0e17370884bb881f840cb1482bb2922cf")
	require.NoError(t, err)
	accountHash, err := casper.NewAccountHash("ba81d84377c5a6c285febca972fe4b531e146cb0fc0f7a5cb9f0e974fc2f6367")
	require.NoError(t, err)
	assert.Equal(t, publicKey.AccountHash().ToHex(), accountHash.ToHex())
}
