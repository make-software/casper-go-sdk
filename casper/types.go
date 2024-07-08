package casper

import (
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/clvalue/cltype"
)

type (
	AuctionState    = types.AuctionState
	Args            = types.Args
	ValidatorBid    = types.ValidatorBid
	AuctionBid      = types.AuctionBid
	BlockV1         = types.BlockV1
	BlockHeaderV1   = types.BlockHeaderV1
	BlockBodyV1     = types.BlockBodyV1
	BlockV2         = types.BlockV2
	BlockHeaderV2   = types.BlockHeaderV2
	BlockBodyV2     = types.BlockBodyV2
	ContractPackage = types.ContractPackage
	ContractVersion = types.ContractVersion
	Contract        = types.Contract
	Deploy          = types.Deploy
	DeployHeader    = types.DeployHeader
	DeployApproval  = types.Approval
	EntryPointV1    = types.EntryPointV1
	EntryPointV2    = types.EntryPointV2
	ExecutionResult = types.ExecutionResultStatus
	NamedKeys       = types.NamedKeys
	NamedKey        = types.NamedKey
	TransformKey    = types.TransformKey
	Transform       = types.Transform
	Argument        = types.Argument
	Account         = types.Account
	Reward          = types.EraReward
	WriteTransfer   = types.WriteTransfer
	UnbondingPurse  = types.UnbondingPurse
)

type (
	ExecutableDeployItem = types.ExecutableDeployItem
	ModuleBytes          = types.ModuleBytes
)

var (
	DefaultHeader   = types.DefaultHeader
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
