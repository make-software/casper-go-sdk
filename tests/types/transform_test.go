package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types"
)

func Test_Transform_AddUInt512(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/AddUInt512.json")
	require.NoError(t, err)
	var transformKey types.TransformKey

	err = json.Unmarshal(fixture, &transformKey)
	require.NoError(t, err)

	val, err := transformKey.Transform.ParseAsUInt512()
	require.NoError(t, err)

	assert.True(t, transformKey.Transform.IsAddUint512())
	assert.EqualValues(t, 100000000, val.Value().Int64())
}

func Test_Transform_ContractV1(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/contract_v1.json")
	require.NoError(t, err)
	var transform types.TransformKey

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)
	assert.True(t, transform.Transform.IsWriteContract())
	contract, err := transform.Transform.ParseAsWriteContract()
	require.Error(t, err)
	require.Nil(t, contract)
}

func Test_Transform_ContractV2(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/contract_v2.json")
	require.NoError(t, err)
	var transform types.Transform

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	contract, err := transform.Kind.ParseAsWriteContract()
	require.NoError(t, err)
	require.NotEmpty(t, contract.ContractPackageHash)
	require.NotEmpty(t, contract.ContractWasmHash)
	require.NotEmpty(t, contract.NamedKeys)
	require.NotEmpty(t, contract.EntryPoints)

	for _, item := range contract.EntryPoints {
		require.NotEmpty(t, item.EntryPoint)
		require.NotEmpty(t, item.EntryPoint.EntryPointType)
		require.NotEmpty(t, item.EntryPoint.Name)
		require.NotEmpty(t, item.EntryPoint.Ret)
		require.NotEmpty(t, item.EntryPoint.Access)
	}

	assert.True(t, transform.Kind.IsWriteContract())
}

func Test_Transform_ContractPackageV1(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/contract_package_v1.json")
	require.NoError(t, err)
	var transform types.TransformKey

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)
	assert.True(t, transform.Transform.IsWriteContractPackage())
}

func Test_Transform_ContractPackageV2(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/contract_package_v2.json")
	require.NoError(t, err)
	var transform types.Transform

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	contractPackage, err := transform.Kind.ParseAsWriteContractPackage()
	require.NoError(t, err)
	require.NotEmpty(t, contractPackage.Versions)
	require.NotEmpty(t, contractPackage.AccessKey)
	assert.True(t, transform.Kind.IsWriteContractPackage())
}

func Test_Transform_CLValue(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/cl_value_v1.json")
	require.NoError(t, err)
	var transform types.TransformKey

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	val, err := transform.Transform.ParseAsWriteCLValue()
	require.NoError(t, err)
	clValue, err := val.Value()
	require.NoError(t, err)

	assert.True(t, transform.Transform.IsWriteCLValue())
	assert.EqualValues(t, "9998335129799990000", clValue.UI512.Value().String())
}

func Test_Transform_CLValue_V2(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/cl_value_v2.json")
	require.NoError(t, err)
	var transform types.Transform

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	val, err := transform.Kind.ParseAsWriteCLValue()
	require.NoError(t, err)
	clValue, err := val.Value()
	require.NoError(t, err)

	assert.True(t, transform.Kind.IsWriteCLValue())
	assert.EqualValues(t, "9998335129799990000", clValue.UI512.Value().String())
}

func Test_Transform_Package(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/package.json")
	require.NoError(t, err)
	var transform types.Transform

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	packageRes, err := transform.Kind.ParseAsWritePackage()
	require.NoError(t, err)

	assert.True(t, transform.Kind.IsWritePackage())
	assert.NotEmpty(t, packageRes.LockStatus)
	assert.NotEmpty(t, packageRes.Versions)
}

func Test_Transform_AddressableEntity(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/addressable_entity.json")
	require.NoError(t, err)
	var transform types.Transform

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	addressableEntity, err := transform.Kind.ParseAsWriteAddressableEntity()
	require.NoError(t, err)

	assert.True(t, transform.Kind.IsWriteAddressableEntity())
	assert.NotEmpty(t, addressableEntity.EntityKind)
	assert.True(t, addressableEntity.EntityKind.Account != nil)
	assert.NotEmpty(t, addressableEntity.MainPurse)
	assert.NotEmpty(t, addressableEntity.AssociatedKeys)
	assert.NotEmpty(t, addressableEntity.PackageHash)
}

func Test_Transform_BidKind(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/bid_kind.json")
	require.NoError(t, err)
	var transform types.Transform

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	bidKind, err := transform.Kind.ParseAsWriteBidKind()
	require.NoError(t, err)

	assert.True(t, transform.Kind.IsWriteBidKind())
	assert.True(t, bidKind.Credit != nil)
	assert.Equal(t, bidKind.Credit.Amount.String(), "100000000")
}

func Test_Transform_NamedKey(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/named_key.json")
	require.NoError(t, err)
	var transform types.Transform

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	namedKey, err := transform.Kind.ParseAsWriteNamedKey()
	require.NoError(t, err)

	nameCLValue, err := namedKey.Name.Value()
	require.NoError(t, err)

	nameKeyCLValue, err := namedKey.NamedKey.Value()
	require.NoError(t, err)

	assert.True(t, transform.Kind.IsWriteNamedKey())
	assert.Equal(t, nameCLValue.StringVal.String(), "my-key-name")
	assert.True(t, nameKeyCLValue.Key != nil)
}

func Test_Transform_Message(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/message.json")
	require.NoError(t, err)
	var transform types.Transform

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	message, err := transform.Kind.ParseAsWriteMessage()
	require.NoError(t, err)

	assert.True(t, transform.Kind.IsWriteMessage())
	assert.Equal(t, string(*message), "message-checksum-987976f5ed2b2843976aaeb5b6d4a810eed7b5a1d7934ef01d167bf46c9c7a8f")
}

func Test_Transform_MessageTopic(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/message_topic.json")
	require.NoError(t, err)
	var transform types.Transform

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	messageTopic, err := transform.Kind.ParseAsWriteMessageTopic()
	require.NoError(t, err)

	assert.True(t, transform.Kind.IsWriteMessageTopic())
	assert.Equal(t, messageTopic.MessageCount, uint32(1))
}

func Test_Transform_WriteDeployInfo(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/WriteDeployInfo.json")
	require.NoError(t, err)
	var transformKey types.TransformKey

	err = json.Unmarshal(fixture, &transformKey)
	require.NoError(t, err)

	val, err := transformKey.Transform.ParseAsWriteDeployInfo()
	require.NoError(t, err)

	assert.True(t, transformKey.Transform.IsWriteDeployInfo())
	assert.EqualValues(t, 1, len(val.Transfers))
}

func Test_Transform_WriteAccountV1(t *testing.T) {
	invalidFixture, err := os.ReadFile("../data/transform/write_clvalue_v1.json")
	require.NoError(t, err)

	var invalidTransform types.TransformKey
	err = json.Unmarshal(invalidFixture, &invalidTransform)
	require.NoError(t, err)
	assert.False(t, invalidTransform.Transform.IsWriteAccount())

	fixture, err := os.ReadFile("../data/transform/write_account_v1.json")
	require.NoError(t, err)
	var transform types.TransformKey

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	writeAccount, err := transform.Transform.ParseAsWriteAccount()
	require.NoError(t, err)
	require.NotEmpty(t, writeAccount.Hash)
	assert.True(t, transform.Transform.IsWriteAccount())
}

func Test_Transform_WriteAccountV2(t *testing.T) {
	fixture, err := os.ReadFile("../data/transform/write_account_v2.json")
	require.NoError(t, err)
	var transform types.Transform

	err = json.Unmarshal(fixture, &transform)
	require.NoError(t, err)

	writeAccount, err := transform.Kind.ParseAsWriteAccount()
	require.NoError(t, err)
	require.NotEmpty(t, writeAccount.Hash)
	assert.True(t, transform.Kind.IsWriteAccount())
}
