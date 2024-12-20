package secp256k1

import (
	"crypto/elliptic"
	"encoding/asn1"
	"encoding/pem"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

var oid = asn1.ObjectIdentifier{1, 3, 132, 0, 10}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

func NewPemPair() ([]byte, []byte, error) {
	priv, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("creating new S256 private key")
	}

	privKeyPem, errPrivExp := PrivateKeyToPem(priv)
	if errPrivExp != nil {
		return nil, nil, fmt.Errorf("export priv key: %v", errPrivExp)
	}

	pubKeyPem, err := PublicKeyToPem(priv.PubKey())
	if err != nil {
		return nil, nil, fmt.Errorf("generating public pem: %s", err)
	}

	return privKeyPem, pubKeyPem, nil
}

func GeneratePrivateKey() (PrivateKey, error) {
	priv, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return PrivateKey{}, fmt.Errorf("creating new S256 private key")
	}
	return PrivateKey{key: priv}, nil
}

func PrivateKeyToPem(priv *secp256k1.PrivateKey) ([]byte, error) {
	key := priv.ToECDSA()

	privateKey := make([]byte, (key.Curve.Params().N.BitLen()+7)/8)
	privBytes, err := asn1.Marshal(ecPrivateKey{
		Version:       1,
		PrivateKey:    key.D.FillBytes(privateKey),
		NamedCurveOID: oid,
		PublicKey:     asn1.BitString{Bytes: elliptic.Marshal(key.Curve, key.X, key.Y)}, //nolint
	})
	if err != nil {
		return nil, fmt.Errorf("marshalling EC private key: %s", err)
	}

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privBytes,
		},
	), nil
}

func PemToPrivateKey(content []byte) (*secp256k1.PrivateKey, error) {
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("key not found")
	}

	var privKey ecPrivateKey
	if _, err := asn1.Unmarshal(block.Bytes, &privKey); err != nil {
		return nil, fmt.Errorf("x509: failed to parse EC private key: %s", err.Error())
	}
	if privKey.Version != 1 {
		return nil, fmt.Errorf("x509: unknown EC private key version %d", privKey.Version)
	}

	return secp256k1.PrivKeyFromBytes(privKey.PrivateKey), nil
}

func PublicKeyToPem(pub *secp256k1.PublicKey) ([]byte, error) {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "EC PUBLIC KEY",
			Bytes: pub.SerializeUncompressed(),
		},
	), nil
}
