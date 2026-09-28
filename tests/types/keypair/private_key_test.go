package keypair

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
	ed25519kp "github.com/make-software/casper-go-sdk/v2/types/keypair/ed25519"
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

func Test_NewPrivateKeyFromHex(t *testing.T) {
	t.Run("Valid ED25519 key", func(t *testing.T) {
		// This is a sample ED25519 private key in hex format
		hexKey := "dda433e404770ebbe9cec28cd7623770ce4222c4961dc4508b076145126c200ece69876ad8154b3c4ec5c2a6ca250e88efda5008cef9ca5ec6767045ee006b53"
		privateKey, err := keypair.NewPrivateKeyFromHex(hexKey, keypair.ED25519)

		require.NoError(t, err)
		assert.NotNil(t, privateKey.PublicKey())
		assert.Equal(t, "01ce69876ad8154b3c4ec5c2a6ca250e88efda5008cef9ca5ec6767045ee006b53", privateKey.PublicKey().ToHex())
		assert.NotNil(t, privateKey)
	})

	t.Run("Valid SECP256K1 key", func(t *testing.T) {
		// This is a sample SECP256K1 private key in hex format
		hexKey := "1e99423a4ed27608a15a2616a2b0e9e52ced330ac530edcc32c8ffc6a526aedd"
		privateKey, err := keypair.NewPrivateKeyFromHex(hexKey, keypair.SECP256K1)

		require.NoError(t, err)
		assert.NotNil(t, privateKey.PublicKey())
	})

	t.Run("Invalid hex for ED25519", func(t *testing.T) {
		hexKey := "invalid_hex"
		_, err := keypair.NewPrivateKeyFromHex(hexKey, keypair.ED25519)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create private key")
	})

	t.Run("Invalid hex for SECP256K1", func(t *testing.T) {
		hexKey := "invalid_hex"
		_, err := keypair.NewPrivateKeyFromHex(hexKey, keypair.SECP256K1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create private key")
	})
}

func Test_NewPrivateKeyFromSeedHex(t *testing.T) {
	t.Run("Valid ED25519 seed", func(t *testing.T) {
		seedHex := "dda433e404770ebbe9cec28cd7623770ce4222c4961dc4508b076145126c200e"
		privateKey, err := keypair.NewPrivateKeyFromSeedHex(seedHex, keypair.ED25519)

		require.NoError(t, err)
		require.NotNil(t, privateKey)
		assert.Equal(t, "01ce69876ad8154b3c4ec5c2a6ca250e88efda5008cef9ca5ec6767045ee006b53", privateKey.PublicKey().ToHex())
	})

	t.Run("Seed matches full private key public key", func(t *testing.T) {
		fullHexKey := "dda433e404770ebbe9cec28cd7623770ce4222c4961dc4508b076145126c200ece69876ad8154b3c4ec5c2a6ca250e88efda5008cef9ca5ec6767045ee006b53"
		fullKey, err := keypair.NewPrivateKeyFromHex(fullHexKey, keypair.ED25519)
		require.NoError(t, err)

		seedHex := "dda433e404770ebbe9cec28cd7623770ce4222c4961dc4508b076145126c200e"
		seedKey, err := keypair.NewPrivateKeyFromSeedHex(seedHex, keypair.ED25519)
		require.NoError(t, err)

		assert.Equal(t, fullKey.PublicKey().ToHex(), seedKey.PublicKey().ToHex())
	})

	t.Run("Seed-derived key can sign and verify", func(t *testing.T) {
		seedHex := "dda433e404770ebbe9cec28cd7623770ce4222c4961dc4508b076145126c200e"
		privateKey, err := keypair.NewPrivateKeyFromSeedHex(seedHex, keypair.ED25519)
		require.NoError(t, err)

		message := []byte("Enigmatic Shadows Concealing Ancient Whispers")
		signature, err := privateKey.Sign(message)
		require.NoError(t, err)
		assert.NoError(t, privateKey.PublicKey().VerifySignature(message, signature))
	})

	t.Run("Invalid hex characters", func(t *testing.T) {
		seedHex := "invalid_hex_string_with_wrong_characters_and_length_to_match_seed_xxxxx"
		_, err := keypair.NewPrivateKeyFromSeedHex(seedHex, keypair.ED25519)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create private key")
	})

	t.Run("Unsupported SECP256K1 algorithm", func(t *testing.T) {
		seedHex := "dda433e404770ebbe9cec28cd7623770ce4222c4961dc4508b076145126c200e"
		_, err := keypair.NewPrivateKeyFromSeedHex(seedHex, keypair.SECP256K1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported key algorithm")
	})
}

func Test_NewPrivateKeyFromSeedBytes(t *testing.T) {
	t.Run("Valid seed bytes", func(t *testing.T) {
		seedHex := "dda433e404770ebbe9cec28cd7623770ce4222c4961dc4508b076145126c200e"
		seed, err := hex.DecodeString(seedHex)
		require.NoError(t, err)

		privateKey, err := ed25519kp.NewPrivateKeyFromSeedBytes(seed)
		require.NoError(t, err)
		require.NotNil(t, privateKey)
	})

	t.Run("Empty seed", func(t *testing.T) {
		_, err := ed25519kp.NewPrivateKeyFromSeedBytes([]byte{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wrong seed size")
	})
}

func Test_Ed25519_NewPrivateKeyFromSeedHex(t *testing.T) {
	t.Run("Invalid hex characters", func(t *testing.T) {
		seedHex := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
		_, err := ed25519kp.NewPrivateKeyFromSeedHex(seedHex)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode hex string")
	})

	t.Run("Wrong hex length", func(t *testing.T) {
		seedHex := "abcd"
		_, err := ed25519kp.NewPrivateKeyFromSeedHex(seedHex)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid hex string length")
	})
}
