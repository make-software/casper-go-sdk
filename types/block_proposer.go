package types

import (
	"encoding/json"
	"errors"

	"github.com/make-software/casper-go-sdk/types/keypair"
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
