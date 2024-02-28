package keypair

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/crypto/blake2b"

	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair/ed25519"
	"github.com/make-software/casper-go-sdk/types/keypair/secp256k1"
)

var (
	ErrEmptySignature       = errors.New("empty signature")
	ErrInvalidPublicKeyAlgo = errors.New("invalid public key algorithm")
	ErrInvalidSignature     = errors.New("invalid signature")
	ErrEmptyPublicKey       = errors.New("empty public key")
)

type PublicKeyInternal interface {
	Bytes() []byte
	VerifySignature(message []byte, sig []byte) bool
}

type PublicKey struct {
	cryptoAlg keyAlgorithm
	key       PublicKeyInternal
}

func (v PublicKey) Bytes() []byte {
	if v.key == nil {
		return nil
	}

	return append([]byte{byte(v.cryptoAlg)}, v.key.Bytes()...)
}

func (v PublicKey) String() string {
	return v.ToHex()
}

func (v PublicKey) MarshalJSON() ([]byte, error) {
	return []byte(`"` + v.String() + `"`), nil
}

func (v PublicKey) ToHex() string {
	return hex.EncodeToString(v.Bytes())
}

func (v *PublicKey) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	hexData, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	tmp, err := NewPublicKeyFromBuffer(bytes.NewBuffer(hexData))
	if err != nil {
		return err
	}
	*v = tmp

	return nil
}

func (v *PublicKey) GobDecode(i []byte) error {
	tmp, err := NewPublicKeyFromBuffer(bytes.NewBuffer(i))
	if err != nil {
		return err
	}
	*v = tmp
	return nil
}

func (v PublicKey) GobEncode() ([]byte, error) {
	return v.Bytes(), nil
}

func (v PublicKey) AccountHash() key.AccountHash {
	if v.key == nil {
		return key.AccountHash{}
	}

	bytesToHash := make([]byte, 0, len(v.cryptoAlg.String())+1+len(v.key.Bytes()))

	bytesToHash = append(bytesToHash, []byte(strings.ToLower(v.cryptoAlg.String()))...)
	bytesToHash = append(bytesToHash, byte(0))
	bytesToHash = append(bytesToHash, v.key.Bytes()...)

	blakeHash := blake2b.Sum256(bytesToHash)
	data, _ := key.NewByteHashFromBuffer(bytes.NewBuffer(blakeHash[:]))

	return key.AccountHash{Hash: data}
}

func (v PublicKey) Value() (driver.Value, error) {
	return v.Bytes(), nil
}

func (v *PublicKey) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan value type")
	}

	dst := make([]byte, len(data))
	copy(dst, data)

	tmp, err := NewPublicKeyFromBuffer(bytes.NewBuffer(dst))
	if err != nil {
		return err
	}
	*v = tmp
	return nil
}

func (v PublicKey) Equals(target PublicKey) bool {
	return v.String() == target.String()
}

// VerifySignature verifies message using Casper compatible cryptographic signature, including the algorithm tag prefix
func (v PublicKey) VerifySignature(message []byte, sig []byte) error {
	if len(sig) <= 1 {
		return ErrEmptySignature
	}

	if v.key == nil {
		return ErrEmptyPublicKey
	}

	// Trim first byte with algorithm data
	sig = sig[1:]

	if v.key.VerifySignature(message, sig) {
		return nil
	}

	return ErrInvalidSignature
}

// VerifyRawSignature verifies message using raw signature
// Deprecated: won't work with Casper node, use VerifySignature method to achieve compatibility
func (v PublicKey) VerifyRawSignature(message []byte, sig []byte) error {
	if v.key == nil {
		return ErrEmptyPublicKey
	}

	if v.key.VerifySignature(message, sig) {
		return nil
	}

	return ErrInvalidSignature
}

func NewPublicKey(source string) (PublicKey, error) {
	data, err := hex.DecodeString(source)
	if err != nil {
		return PublicKey{}, err
	}
	return NewPublicKeyFromBuffer(bytes.NewBuffer(data))
}

func NewPublicKeyFromBytes(source []byte) (PublicKey, error) {
	return NewPublicKeyFromBuffer(bytes.NewBuffer(source))
}

func NewPublicKeyFromBuffer(buf *bytes.Buffer) (PublicKey, error) {
	alg, err := buf.ReadByte()
	if err != nil {
		return PublicKey{}, err
	}
	var result PublicKey
	result.cryptoAlg = keyAlgorithm(alg)
	switch result.cryptoAlg {
	case ED25519:
		if result.key, err = ed25519.NewPublicKey(buf.Next(ed25519.PublicKeySize)); err != nil {
			return PublicKey{}, err
		}
	case SECP256K1:
		if result.key, err = secp256k1.NewPublicKey(buf.Next(secp256k1.PublicKeySize)); err != nil {
			return PublicKey{}, err
		}
	default:
		return PublicKey{}, ErrInvalidPublicKeyAlgo
	}
	return result, nil
}

type PublicKeyList []PublicKey

func (p PublicKeyList) Contains(target PublicKey) bool {
	for _, one := range p {
		if one.Equals(target) {
			return true
		}
	}
	return false
}
