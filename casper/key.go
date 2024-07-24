package casper

import (
	"github.com/make-software/casper-go-sdk/v2/types/key"
	"github.com/make-software/casper-go-sdk/v2/types/keypair"
)

type (
	Key                 = key.Key
	Hash                = key.Hash
	DeployHash          = key.Hash
	AccountHash         = key.AccountHash
	ContractHash        = key.ContractHash
	ContractPackageHash = key.ContractPackageHash
	TransferHash        = key.TransferHash
	Uref                = key.URef
	PublicKey           = keypair.PublicKey
	PublicKeyList       = keypair.PublicKeyList
	PrivateKey          = keypair.PrivateKey
)

var (
	NewAccountHash                    = key.NewAccountHash
	NewTransferHash                   = key.NewTransferHash
	NewHash                           = key.NewHash
	NewHashFromBytes                  = key.NewHashFromBytes
	NewContractHash                   = key.NewContract
	NewContractPackageHash            = key.NewContractPackage
	NewUref                           = key.NewURef
	NewKey                            = key.NewKey
	NewKeyFromBytes                   = key.NewKeyFromBytes
	NewPublicKey                      = keypair.NewPublicKey
	NewED25519PrivateKeyFromPEMFile   = keypair.NewPrivateKeyED25518
	NewSECP256k1PrivateKeyFromPEMFile = keypair.NewPrivateKeySECP256K1
)
