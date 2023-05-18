package secp256k1

import (
	"crypto/elliptic"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

var oid = asn1.ObjectIdentifier{1, 3, 132, 0, 10}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

type pkixPublicKey struct {
	Algo      pkix.AlgorithmIdentifier
	BitString asn1.BitString
}

func NewPemPair() ([]byte, []byte, error) {
	priv, err := btcec.NewPrivateKey(btcec.S256())
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

func PrivateKeyToPem(priv *btcec.PrivateKey) ([]byte, error) {
	key := priv.ToECDSA()

	privateKey := make([]byte, (key.Curve.Params().N.BitLen()+7)/8)
	privBytes, err := asn1.Marshal(ecPrivateKey{
		Version:       1,
		PrivateKey:    key.D.FillBytes(privateKey),
		NamedCurveOID: oid,
		PublicKey:     asn1.BitString{Bytes: elliptic.Marshal(key.Curve, key.X, key.Y)},
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

func PemToPrivateKey(priv []byte) (*btcec.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	if block == nil {
		return nil, fmt.Errorf("key not found")
	}

	var privKey ecPrivateKey
	if _, err := asn1.Unmarshal(block.Bytes, &privKey); err != nil {
		return nil, fmt.Errorf("x509: failed to parse EC private key: " + err.Error())
	}
	if privKey.Version != 1 {
		return nil, fmt.Errorf("x509: unknown EC private key version %d", privKey.Version)
	}

	curve := btcec.S256()

	k := new(big.Int).SetBytes(privKey.PrivateKey)
	curveOrder := curve.Params().N
	if k.Cmp(curveOrder) >= 0 {
		return nil, fmt.Errorf("x509: invalid elliptic curve private key value")
	}

	key := new(btcec.PrivateKey)
	key.Curve = curve
	key.D = k

	privateKey := make([]byte, (curveOrder.BitLen()+7)/8)

	for len(privKey.PrivateKey) > len(privateKey) {
		if privKey.PrivateKey[0] != 0 {
			return nil, fmt.Errorf("x509: invalid private key length")
		}
		privKey.PrivateKey = privKey.PrivateKey[1:]
	}

	copy(privateKey[len(privateKey)-len(privKey.PrivateKey):], privKey.PrivateKey)
	key.X, key.Y = curve.ScalarBaseMult(privateKey)

	return key, nil
}

func PublicKeyToPem(pub *btcec.PublicKey) ([]byte, error) {
	pubEDSA := pub.ToECDSA()

	var publicKeyAlgorithm pkix.AlgorithmIdentifier

	publicKeyBytes := elliptic.Marshal(pubEDSA.Curve, pubEDSA.X, pubEDSA.Y)

	publicKeyAlgorithm.Algorithm = oid
	var paramBytes []byte
	paramBytes, err := asn1.Marshal(oid)
	if err != nil {
		return nil, err
	}

	publicKeyAlgorithm.Parameters.FullBytes = paramBytes

	pubBytes, _ := asn1.Marshal(pkixPublicKey{
		Algo: publicKeyAlgorithm,
		BitString: asn1.BitString{
			Bytes:     publicKeyBytes,
			BitLength: 8 * len(publicKeyBytes),
		},
	})

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "EC PUBLIC KEY",
			Bytes: pubBytes,
		},
	), nil
}
