package casper

import (
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue/cltype"
)

type (
	AuctionStateV1       = types.AuctionStateV1
	AuctionStateV2       = types.AuctionStateV2
	Args                 = types.Args
	ValidatorBid         = types.ValidatorBid
	Bid                  = types.Bid
	Block                = types.Block
	BlockV1              = types.BlockV1
	BlockHeaderV1        = types.BlockHeaderV1
	BlockBodyV1          = types.BlockBodyV1
	BlockV2              = types.BlockV2
	BlockHeaderV2        = types.BlockHeaderV2
	BlockBodyV2          = types.BlockBodyV2
	ContractPackage      = types.ContractPackage
	ContractVersion      = types.ContractVersion
	Contract             = types.Contract
	Transaction          = types.Transaction
	TransactionV1        = types.TransactionV1
	TransactionV1Payload = types.TransactionV1Payload
	Deploy               = types.Deploy
	DeployHeader         = types.DeployHeader
	DeployApproval       = types.Approval
	EntryPointV1         = types.EntryPointV1
	EntryPointV2         = types.EntryPointV2
	ExecutionResult      = types.ExecutionResult
	NamedKeys            = types.NamedKeys
	NamedKey             = types.NamedKey
	TransformKey         = types.TransformKey
	Transform            = types.Transform
	Argument             = types.Argument
	Account              = types.Account
	PublicKeyAndBid      = types.PublicKeyAndBid
	Reward               = types.EraReward
	WriteTransfer        = types.WriteTransfer
	UnbondingPurse       = types.UnbondingPurse
)

type (
	ExecutableDeployItem = types.ExecutableDeployItem
	ModuleBytes          = types.ModuleBytes
)

var (
	DefaultHeader   = types.DefaultDeployHeader
	MakeDeploy      = types.MakeDeploy
	StandardPayment = types.StandardPayment
)

type (
	CLType  = cltype.CLType
	CLValue = clvalue.CLValue
	CLMap   = clvalue.Map
)

type (
	StoredContractByHash          = types.StoredContractByHash
	StoredContractByName          = types.StoredContractByName
	StoredValue                   = types.StoredValue
	StoredVersionedContractByHash = types.StoredVersionedContractByHash
	StoredVersionedContractByName = types.StoredVersionedContractByName
)

var (
	ErrArgumentNotFound = types.ErrArgumentNotFound
)
