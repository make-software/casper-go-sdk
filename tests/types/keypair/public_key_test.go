package keypair

import (
	"encoding/hex"
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

func Test_PublicKey_VerifySignature_SECP256k1(t *testing.T) {
	publicKey, err := casper.NewPublicKey("0203a0cfb64a420751875e25372050e884b4017bad8cc5b038c7c3e24261d2b09712")
	require.NoError(t, err)
	messageStr := `Casper Message:
{"email":"evgeniy@make.services","format_id":2,"public_key":"0203a0cfb64a420751875e25372050e884b4017bad8cc5b038c7c3e24261d2b09712","start_of_period":"2022-01-01T00:00:00.000Z","tax_period_id":5}`
	signatureStr := "cb211fe5e7d2a8fd12c020200bd7fab809587b9468d99d2c316545832e09f8b27e346768f99df77d7c0e5436943ff5cc18fb2cc0eef1ed6f894dde5ff1044ad3"
	signatureBytes, err := hex.DecodeString(signatureStr)
	require.NoError(t, err)
	require.NoError(t, publicKey.VerifySignature([]byte(messageStr), signatureBytes))
}
