package types

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

type Proposer struct {
	isSystem  bool
	publicKey *keypair.PublicKey
}

func NewProposer(src string) (Proposer, error) {
	var result Proposer
	if src == "00" {
		result.isSystem = true
		return result, nil
	}
	pubKey, err := keypair.NewPublicKey(src)
	if err != nil {
		return result, err
	}
	result.publicKey = &pubKey
	return result, nil
}

func (p Proposer) IsSystem() bool {
	return p.isSystem
}

func (p Proposer) PublicKey() (keypair.PublicKey, error) {
	if p.IsSystem() {
		return keypair.PublicKey{}, errors.New("system proposer doesn't have a PublicKey")
	}
	return *p.publicKey, nil
}

func (p Proposer) PublicKeyOptional() *keypair.PublicKey {
	return p.publicKey
}

func (p Proposer) MarshalJSON() ([]byte, error) {
	if p.isSystem {
		return []byte(`"00"`), nil
	}
	return json.Marshal(p.publicKey)
}

func (p *Proposer) UnmarshalJSON(bytes []byte) error {
	var str string
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return err
	}
	*p, err = NewProposer(str)
	if err != nil {
		return err
	}
	return nil
}

func (p *Proposer) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan value type")
	}
	if hex.EncodeToString(b) == "00" {
		p.isSystem = true
		return nil
	}
	var pubKey keypair.PublicKey
	if err := pubKey.Scan(value); err != nil {
		return err
	}
	p.publicKey = &pubKey
	return nil
}

func (p Proposer) Value() (driver.Value, error) {
	if p.isSystem {
		return hex.DecodeString("00")
	}
	return p.publicKey.Value()
}
