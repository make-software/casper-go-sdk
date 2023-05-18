package types

import (
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/casper-go-sdk/types/keypair"
)

// Deploy is an item containing a smart contract along with the requester's signature(s).
type Deploy struct {
	// List of signers and signatures for this `deploy`.
	Approvals []Approval `json:"approvals"`
	// A hash over the header of the `deploy`.
	Hash key.Hash `json:"hash"`
	// Contains metadata about the `deploy`.
	Header DeployHeader `json:"header"`
	// Contains the payment amount for the `deploy`.
	Payment ExecutableDeployItem `json:"payment"`
	// Contains the session information for the `deploy`.
	Session ExecutableDeployItem `json:"session"`
}

type Approval struct {
	Signature HexBytes          `json:"signature"`
	Signer    keypair.PublicKey `json:"signer"`
}

type DeployHeader struct {
	// Public Key from the `AccountHash` owning the `Deploy`.
	Account keypair.PublicKey `json:"account"`
	// `Hash` of the body part of this `Deploy`.
	BodyHash  key.Hash `json:"body_hash"`
	ChainName string   `json:"chain_name"`
	// List of `Deploy` hashes.
	Dependencies []key.Hash `json:"dependencies"`
	GasPrice     uint64     `json:"gas_price"`
	// `Timestamp` formatted as per RFC 3339
	Timestamp Timestamp `json:"timestamp"`
	// Duration of the `Deploy` in milliseconds (from timestamp).
	TTL Duration `json:"ttl"`
}

func DefaultHeader() DeployHeader {
	return DeployHeader{
		Dependencies: []key.Hash{},
		GasPrice:     1,
		Timestamp:    Timestamp(time.Now()),
		TTL:          Duration(30 * time.Minute),
	}
}

func (d DeployHeader) Bytes() []byte {
	var result []byte
	result = append(result, d.Account.Bytes()...)
	result = append(result, clvalue.NewCLUInt64(uint64(time.Time(d.Timestamp).UnixMilli())).Bytes()...)
	result = append(result, clvalue.NewCLUInt64(uint64(time.Duration(d.TTL).Milliseconds())).Bytes()...)
	result = append(result, clvalue.NewCLUInt64(d.GasPrice).Bytes()...)
	result = append(result, d.BodyHash.Bytes()...)
	depsList := clvalue.NewCLList(cltype.NewByteArray(key.ByteHashLen))
	for _, one := range d.Dependencies {
		depsList.List.Append(clvalue.NewCLByteArray(one.Bytes()))
	}
	result = append(result, depsList.Bytes()...)
	result = append(result, clvalue.NewCLString(d.ChainName).Bytes()...)

	return result
}

func (d *Deploy) ValidateDeploy() (bool, error) {
	paymentBytes, err := d.Payment.Bytes()
	if err != nil {
		return false, err
	}
	sessionBytes, err := d.Session.Bytes()
	if err != nil {
		return false, err
	}
	if d.Header.BodyHash != blake2b.Sum256(append(paymentBytes, sessionBytes...)) {
		return false, nil
	}

	if d.Hash != blake2b.Sum256(d.Header.Bytes()) {
		return false, nil
	}

	for _, one := range d.Approvals {
		if one.Signer.VerifySignature(d.Hash.Bytes(), one.Signature) != nil {
			return false, nil
		}
	}

	return true, nil
}

func (d *Deploy) SignDeploy(keys keypair.PrivateKey) error {
	signature, err := keys.Sign(d.Hash.Bytes())
	if err != nil {
		return err
	}
	approval := Approval{
		Signer:    keys.PublicKey(),
		Signature: signature,
	}

	if d.Approvals == nil {
		d.Approvals = make([]Approval, 0, 1)
	}

	d.Approvals = append(d.Approvals, approval)
	return nil
}

func NewDeploy(hash key.Hash, header DeployHeader, payment ExecutableDeployItem, sessions ExecutableDeployItem,
	approvals []Approval) *Deploy {
	d := new(Deploy)
	d.Hash = hash
	d.Header = header
	d.Payment = payment
	d.Session = sessions
	d.Approvals = approvals
	return d
}

func MakeDeploy(deployHeader DeployHeader, payment ExecutableDeployItem, session ExecutableDeployItem) (*Deploy, error) {
	paymentBytes, err := payment.Bytes()
	if err != nil {
		return nil, err
	}
	sessionBytes, err := session.Bytes()
	if err != nil {
		return nil, err
	}
	serializedBody := append(paymentBytes, sessionBytes...)
	deployHeader.BodyHash = blake2b.Sum256(serializedBody)
	deployHash := blake2b.Sum256(deployHeader.Bytes())
	approvals := make([]Approval, 0)

	return NewDeploy(deployHash, deployHeader, payment, session, approvals), nil
}
