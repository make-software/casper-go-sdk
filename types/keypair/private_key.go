package keypair

import (
	"errors"

	"github.com/make-software/casper-go-sdk/types/keypair/ed25519"
	"github.com/make-software/casper-go-sdk/types/keypair/secp256k1"
)

type PrivateKeyInternal interface {
	PublicKeyBytes() []byte
	Sign(mes []byte) ([]byte, error)
}

type PrivateKey struct {
	alg  keyAlgorithm
	pub  PublicKey
	priv PrivateKeyInternal
}

func (v PrivateKey) PublicKey() PublicKey {
	return v.pub
}

// Sign creates a Casper compatible cryptographic signature, including the algorithm tag prefix
func (v PrivateKey) Sign(msg []byte) ([]byte, error) {
	sign, err := v.priv.Sign(msg)
	if err != nil {
		return nil, err
	}
	return append([]byte{v.alg.Byte()}, sign...), nil
}

// RawSign returns raw bytes of signature to sign off chain data
// Deprecated: won't work with Casper node, use Sign method instead
func (v PrivateKey) RawSign(mes []byte) ([]byte, error) {
	return v.priv.Sign(mes)
}

func NewPrivateKeyED25518(path string) (PrivateKey, error) {
	return NewPrivateKeyFromFile(path, ED25519)
}

func NewPrivateKeySECP256K1(path string) (PrivateKey, error) {
	return NewPrivateKeyFromFile(path, SECP256K1)
}

func NewPrivateKeyFromFile(path string, algorithm keyAlgorithm) (PrivateKey, error) {
	var priv PrivateKeyInternal
	var err error
	switch algorithm {
	case ED25519:
		priv, err = ed25519.NewPrivateKeyFromPEMFile(path)
		if err != nil {
			return PrivateKey{}, err
		}
	case SECP256K1:
		priv, err = secp256k1.NewPrivateKeyFromPemFile(path)
		if err != nil {
			return PrivateKey{}, err
		}
	default:
		return PrivateKey{}, errors.New("")
	}
	publicKey, err := NewPublicKeyFromBytes(append([]byte{byte(algorithm)}, priv.PublicKeyBytes()...))
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		alg:  algorithm,
		pub:  publicKey,
		priv: priv,
	}, nil
}

func GeneratePrivateKey(algorithm keyAlgorithm) (PrivateKey, error) {
	var priv PrivateKeyInternal
	var err error
	switch algorithm {
	case ED25519:
		if priv, err = ed25519.GeneratePrivateKey(); err != nil {
			return PrivateKey{}, err
		}
	case SECP256K1:
		priv, err = secp256k1.GeneratePrivateKey()
		if err != nil {
			return PrivateKey{}, err
		}
	}
	publicKey, err := NewPublicKeyFromBytes(append([]byte{byte(algorithm)}, priv.PublicKeyBytes()...))
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		alg:  algorithm,
		pub:  publicKey,
		priv: priv,
	}, nil
}
