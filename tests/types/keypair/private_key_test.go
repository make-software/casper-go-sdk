package keypair

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

func Test_ED25519_PrivateKey_Parsing(t *testing.T) {
	privateKeyData, err := casper.NewED25519PrivateKeyFromPEMFile("../../data/keys/account_test_ED25519_secret_key.pem")
	require.NoError(t, err)
	assert.Equal(t, "015a372b0e230bf9393e2df0b3de857bb0e17370884bb881f840cb1482bb2922cf", privateKeyData.PublicKey().ToHex())
}

func Test_SECPKey_From_PemFile(t *testing.T) {
	privateKeyData, err := casper.NewSECP256k1PrivateKeyFromPEMFile("../../data/keys/account_test_SECP_secret_key.pem")
	require.NoError(t, err)
	assert.Equal(t, "0203c90c0ee375abc85da81a982507d1f8258a380af2058b63c37abdb9c7045940f4", privateKeyData.PublicKey().ToHex())
}

func Test_SECPKey_ToPemFile(t *testing.T) {
	privateKeyData, err := keypair.NewPrivateKeySECP256K1("../../data/keys/account_test_SECP_secret_key.pem")
	require.NoError(t, err)
	data, err := privateKeyData.ToPem()
	require.NoError(t, err)
	privateKeyData2, err := keypair.NewPrivateKeyFromPEM(data, keypair.SECP256K1)
	require.NoError(t, err)
	assert.Equal(t, privateKeyData.PublicKey().Bytes(), privateKeyData2.PublicKey().Bytes())
}

func Test_ED25519_PrivateKey_ToPemFile(t *testing.T) {
	privateKeyData, err := keypair.NewPrivateKeyED25518("../../data/keys/account_test_ED25519_secret_key.pem")
	require.NoError(t, err)
	data, err := privateKeyData.ToPem()
	require.NoError(t, err)
	privateKeyData2, err := keypair.NewPrivateKeyFromPEM(data, keypair.ED25519)
	require.NoError(t, err)
	assert.Equal(t, privateKeyData.PublicKey().Bytes(), privateKeyData2.PublicKey().Bytes())
}

func Test_SECPKey_CreateAndValidateSignature(t *testing.T) {
	secretMessage := []byte("Enigmatic Shadows Concealing Ancient Whispers")
	privateKeyData, err := keypair.GeneratePrivateKey(keypair.SECP256K1)
	require.NoError(t, err)
	signature, err := privateKeyData.Sign(secretMessage)
	require.NoError(t, err)
	err = privateKeyData.PublicKey().VerifySignature(secretMessage, signature)
	assert.NoError(t, err)
}

func Test_SECPKey_CreateAndValidateRawSignature(t *testing.T) {
	secretMessage := []byte("Enigmatic Shadows Concealing Ancient Whispers")
	privateKeyData, err := keypair.GeneratePrivateKey(keypair.SECP256K1)
	require.NoError(t, err)
	signature, err := privateKeyData.RawSign(secretMessage)
	require.NoError(t, err)
	err = privateKeyData.PublicKey().VerifyRawSignature(secretMessage, signature)
	assert.NoError(t, err)
}

func Test_ED25PKey_CreateAndValidateSignature(t *testing.T) {
	secretMessage := []byte("Enigmatic Shadows Concealing Ancient Whispers")
	privateKeyData, err := keypair.GeneratePrivateKey(keypair.ED25519)
	require.NoError(t, err)
	signature, err := privateKeyData.Sign(secretMessage)
	require.NoError(t, err)
	err = privateKeyData.PublicKey().VerifySignature(secretMessage, signature)
	assert.NoError(t, err)
}

func Test_ED25Key_CreateAndValidateRawSignature(t *testing.T) {
	secretMessage := []byte("Enigmatic Shadows Concealing Ancient Whispers")
	privateKeyData, err := keypair.GeneratePrivateKey(keypair.ED25519)
	require.NoError(t, err)
	signature, err := privateKeyData.RawSign(secretMessage)
	require.NoError(t, err)
	err = privateKeyData.PublicKey().VerifyRawSignature(secretMessage, signature)
	assert.NoError(t, err)
}
