package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/rpc"
	"github.com/make-software/casper-go-sdk/v2/types"
)

func SetupServer(t *testing.T, filePath string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var requestParams rpc.RpcRequest
		err := decoder.Decode(&requestParams)
		require.NoError(t, err)
		if requestParams.Method == rpc.MethodGetStateRootHash {
			fixture, err := os.ReadFile("../data/rpc_response/get_root_state_hash.json")
			require.NoError(t, err)
			_, err = rw.Write(fixture)
			require.NoError(t, err)
			return
		}
		fixture, err := os.ReadFile(filePath)
		require.NoError(t, err)
		_, err = rw.Write(fixture)
		require.NoError(t, err)
	}))
	return server
}

func Test_DefaultClient_GetDeploy_Example(t *testing.T) {
	tests := []struct {
		filePath      string
		failed        bool
		withTransfers bool
	}{
		{
			filePath:      "../data/deploy/get_raw_rpc_deploy.json",
			withTransfers: true,
		},
		{
			failed:   true,
			filePath: "../data/deploy/get_raw_rpc_deploy_v1_failed.json",
		},
		{
			filePath: "../data/deploy/get_raw_rpc_deploy_v2.json",
		},
	}

	for _, tt := range tests {
		t.Run("GetDeploy", func(t *testing.T) {
			deployHash := "0009ea4441f4700325d9c38b0b6df415537596e1204abe4f6a94b6996aebf2f1"
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				fixture, err := os.ReadFile(tt.filePath)
				require.NoError(t, err)
				_, err = rw.Write(fixture)
				require.NoError(t, err)
			}))
			defer server.Close()

			client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
			deployResult, err := client.GetDeploy(context.Background(), deployHash)
			require.NoError(t, err)

			assert.NotEmpty(t, deployResult.ApiVersion)
			assert.NotEmpty(t, deployResult.Deploy)
			assert.NotEmpty(t, deployResult.Deploy.Session)
			assert.NotEmpty(t, deployResult.Deploy.Approvals)
			assert.NotEmpty(t, deployResult.Deploy.Hash)
			assert.NotEmpty(t, deployResult.Deploy.Payment)
			assert.NotEmpty(t, deployResult.ExecutionResults.ExecutionResult)
			if tt.failed {
				assert.NotEmpty(t, deployResult.ExecutionResults.ExecutionResult.ErrorMessage)
			}

			if tt.withTransfers {
				assert.NotEmpty(t, deployResult.ExecutionResults.ExecutionResult.Transfers)
			}
		})
	}
}

func Test_DefaultClient_GetTransaction_Example(t *testing.T) {
	tests := []struct {
		filePath          string
		isDeploy          bool
		withTransfers     bool
		executionResultV1 bool
	}{
		{
			filePath:          "../data/deploy/get_raw_rpc_deploy.json",
			isDeploy:          true,
			executionResultV1: true,
		},
		{
			filePath: "../data/deploy/get_raw_rpc_deploy_v2.json",
			isDeploy: true,
		},
		{
			filePath:          "../data/deploy/get_raw_rpc_deploy_with_transfer.json",
			isDeploy:          true,
			withTransfers:     true,
			executionResultV1: true,
		},
		{
			filePath: "../data/transaction/get_transaction.json",
		},
		{
			filePath: "../data/transaction/get_transaction_native_target.json",
		},
	}
	for _, tt := range tests {
		t.Run("GetTransaction", func(t *testing.T) {
			server := SetupServer(t, tt.filePath)
			defer server.Close()
			client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
			result, err := client.GetTransactionByTransactionHash(context.Background(), "0009ea4441f4700325d9c38b0b6df415537596e1204abe4f6a94b6996aebf2f1")
			require.NoError(t, err)
			assert.NotEmpty(t, result.APIVersion)
			assert.NotEmpty(t, result.Transaction.Hash)
			assert.NotEmpty(t, result.Transaction.Header)
			assert.NotEmpty(t, result.Transaction.Header.TTL)
			assert.NotEmpty(t, result.Transaction.Header.ChainName)
			assert.NotEmpty(t, result.Transaction.Header.PricingMode)
			assert.NotEmpty(t, result.Transaction.Header.InitiatorAddr)
			assert.NotEmpty(t, result.Transaction.Body.Target)
			assert.NotEmpty(t, result.Transaction.Body.Scheduling)
			assert.NotEmpty(t, result.ExecutionInfo.ExecutionResult.Initiator)
			assert.NotEmpty(t, result.ExecutionInfo.ExecutionResult.Effects)
			assert.NotEmpty(t, result.Transaction.Approvals)

			if tt.isDeploy {
				assert.NotEmpty(t, result.Transaction.GetDeploy())
			}

			if tt.executionResultV1 {
				assert.NotEmpty(t, result.ExecutionInfo.ExecutionResult.GetExecutionResultV1())
			} else {
				assert.NotEmpty(t, result.ExecutionInfo.ExecutionResult.GetExecutionResultV2())
			}

			if tt.withTransfers {
				assert.NotEmpty(t, result.ExecutionInfo.ExecutionResult.Transfers)
			}
		})
	}
}

func Test_DefaultClient_GetDeploy(t *testing.T) {
	server := SetupServer(t, "../data/deploy/get_raw_rpc_deploy.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	deployHash := "a2c450eb80c408105dcf5a6808786a2681d4b7ef8bffd6bb59ccbbee98b908fb"
	result, err := client.GetDeploy(context.Background(), deployHash)
	require.NoError(t, err)
	assert.Equal(t, deployHash, result.Deploy.Hash.ToHex())
}

func Test_DefaultClient_GetDeployFinalizedApproval(t *testing.T) {
	server := SetupServer(t, "../data/deploy/get_raw_rpc_deploy.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	deployHash := "a2c450eb80c408105dcf5a6808786a2681d4b7ef8bffd6bb59ccbbee98b908fb"
	result, err := client.GetDeployFinalizedApproval(context.Background(), deployHash)
	require.NoError(t, err)
	assert.Equal(t, deployHash, result.Deploy.Hash.ToHex())
}

func Test_DefaultClient_GetStateItem_GetAccount(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_state_item_account.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))

	hash := "account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59"
	stateRootHash := "fb9c42717769d72442ff17a5ff1574b4bc1c83aedf5992b14e4d071423f86240"
	result, err := client.GetStateItem(
		context.Background(),
		&stateRootHash,
		hash,
		[]string{},
	)
	require.NoError(t, err)
	assert.Equal(t, hash, result.StoredValue.Account.AccountHash.ToPrefixedString())
}

func Test_DefaultClient_GetStateItem_GetAccount_WithEmptyStateRootHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_state_item_account.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	hash := "account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59"
	result, err := client.GetStateItem(
		context.Background(),
		nil,
		hash,
		[]string{},
	)
	require.NoError(t, err)
	assert.Equal(t, hash, result.StoredValue.Account.AccountHash.ToPrefixedString())
}

func Test_DefaultClient_StateGetDictionaryItem_GetCValueUI64(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_dictionary_item_ui64.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	stateRootHash := "0808080808080808080808080808080808080808080808080808080808080808"
	uref := "uref-09480c3248ef76b603d386f3f4f8a5f87f597d4eaffd475433f861af187ab5db-007"
	result, err := client.GetDictionaryItem(context.Background(), &stateRootHash, uref, "a_unique_entry_identifier")
	require.NoError(t, err)
	value, err := result.StoredValue.CLValue.Value()
	require.NoError(t, err)
	assert.Equal(t, 1, int(value.UI64.Value()))
}

func Test_DefaultClient_StateGetDictionaryItem_GetCValueUI64_WithEmptyStateRootHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_dictionary_item_ui64.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	uref := "uref-09480c3248ef76b603d386f3f4f8a5f87f597d4eaffd475433f861af187ab5db-007"
	result, err := client.GetDictionaryItem(context.Background(), nil, uref, "a_unique_entry_identifier")
	require.NoError(t, err)
	value, err := result.StoredValue.CLValue.Value()
	require.NoError(t, err)
	assert.Equal(t, 1, int(value.UI64.Value()))
}

func Test_DefaultClient_QueryGlobalStateByBlock_GetAccount(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_era.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	blockHash := "bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59"
	accountKey := "account-hash-e94daaff79c2ab8d9c31d9c3058d7d0a0dd31204a5638dc1451fa67b2e3fb88c"
	res, err := client.QueryGlobalStateByBlockHash(context.Background(), blockHash, accountKey, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.BlockHeader.BodyHash)
	assert.NotEmpty(t, res.StoredValue.Account.AccountHash)
}

func Test_DefaultClient_QueryGlobalStateByBlock_V2(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_v2.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	blockHash := "bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59"
	accountKey := "account-hash-e94daaff79c2ab8d9c31d9c3058d7d0a0dd31204a5638dc1451fa67b2e3fb88c"
	res, err := client.QueryGlobalStateByBlockHash(context.Background(), blockHash, accountKey, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.BlockHeader.BodyHash)
	assert.NotEmpty(t, res.StoredValue.Account.AccountHash)
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_GetAccount(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_era.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	accountKey := "hash-8e08c43f144a13c915cf3681cc97bcd98c6a81d7b5da5164dc066318ec1c80a7"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, accountKey, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.BlockHeader.BodyHash)
	assert.NotEmpty(t, res.StoredValue.Account.AccountHash)
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_StoredAddressableEntity_Account(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_addressable_entity_account.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	key := "entity-account-989ca079a5e446071866331468ab949483162588d57ec13ba6bb051f1e15f8b7"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, key, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.EntityKind.Account)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.ActionThresholds)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.AssociatedKeys)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.ByteCodeHash)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.PackageHash)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.MainPurse)
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_StoredAddressableEntity_ByteCode(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_byte_code.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	key := "entity-account-989ca079a5e446071866331468ab949483162588d57ec13ba6bb051f1e15f8b7"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, key, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.ByteCode)
	assert.NotEmpty(t, res.StoredValue.ByteCode.Kind)
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_StoredAddressableEntity_EntryPoint(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_entry_point.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	key := "entry-point-v1-entity-contract-a5cf5917505ef60a6f0df395dd19e86a0f075d00f2e6ce49f5aa0e18f6e26f5d-4ca60287ae6129662475a8ce0d41c450d072b2430a8759f6178adeeff38523da"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, key, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.EntryPoint)
	assert.NotEmpty(t, res.StoredValue.EntryPoint.V1CasperVm)
	assert.NotEmpty(t, res.StoredValue.EntryPoint.V1CasperVm.Name)
	assert.NotEmpty(t, res.StoredValue.EntryPoint.V1CasperVm.Ret)
	assert.NotEmpty(t, res.StoredValue.EntryPoint.V1CasperVm.Access)
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_StoredAddressableEntity_Message(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_message.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	key := "entity-account-989ca079a5e446071866331468ab949483162588d57ec13ba6bb051f1e15f8b7"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, key, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.Message)
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_StoredAddressableEntity_MessageTopic(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_message_topic.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	key := "entity-account-989ca079a5e446071866331468ab949483162588d57ec13ba6bb051f1e15f8b7"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, key, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.MessageTopic)
	assert.NotEmpty(t, res.StoredValue.MessageTopic.MessageCount)
	assert.NotEmpty(t, res.StoredValue.MessageTopic.BlockTime)
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_StoredAddressableEntity_System(t *testing.T) {
	tests := []struct {
		filePath   string
		systemType string
	}{
		{
			filePath:   "../data/rpc_response/query_global_state_addressable_entity_system.json",
			systemType: "HandlePayment",
		},
		{
			filePath:   "../data/rpc_response/query_global_state_addressable_entity_system_mint.json",
			systemType: "Mint",
		},
		{
			filePath:   "../data/rpc_response/query_global_state_addressable_entity_system_auction.json",
			systemType: "Auction",
		},
	}
	for _, tt := range tests {
		t.Run("QueryGlobalStateByBlockHeight", func(t *testing.T) {
			server := SetupServer(t, tt.filePath)
			defer server.Close()
			client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
			key := "entity-account-989ca079a5e446071866331468ab949483162588d57ec13ba6bb051f1e15f8b7"
			res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, key, nil)
			require.NoError(t, err)
			assert.NotEmpty(t, res.StoredValue.AddressableEntity)
			assert.NotEmpty(t, res.StoredValue.AddressableEntity.EntityKind.System)
			assert.Equal(t, *res.StoredValue.AddressableEntity.EntityKind.System, types.SystemEntityType(tt.systemType))
			assert.NotEmpty(t, res.StoredValue.AddressableEntity.ActionThresholds)
			assert.NotEmpty(t, res.StoredValue.AddressableEntity.ByteCodeHash)
			assert.NotEmpty(t, res.StoredValue.AddressableEntity.PackageHash)
		})
	}
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_StoredAddressableEntity_Package(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_package.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	key := "package-8e08c43f144a13c915cf3681cc97bcd98c6a81d7b5da5164dc066318ec1c80a7"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, key, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.Package)
	assert.NotEmpty(t, res.StoredValue.Package.Versions)
	assert.NotEmpty(t, res.StoredValue.Package.LockStatus)
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_StoredAddressableEntity_Contract(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_addressable_entity_contract.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	key := "entity-contract-a5cf5917505ef60a6f0df395dd19e86a0f075d00f2e6ce49f5aa0e18f6e26f5d"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, key, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.EntityKind)
	assert.Equal(t, *res.StoredValue.AddressableEntity.EntityKind.SmartContract, types.TransactionRuntimeVmCasperV1)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.ActionThresholds)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.AssociatedKeys)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.ByteCodeHash)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.PackageHash)
	assert.NotEmpty(t, res.StoredValue.AddressableEntity.MainPurse)
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_StoredAddressableEntity_NamedKey_Account(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_named_key_account.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	key := "entity-contract-a5cf5917505ef60a6f0df395dd19e86a0f075d00f2e6ce49f5aa0e18f6e26f5d"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, key, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.NamedKey)
	assert.NotEmpty(t, res.StoredValue.NamedKey.Name)
	assert.NotEmpty(t, res.StoredValue.NamedKey.NamedKey)
	assert.Equal(t, res.StoredValue.NamedKey.Name.StringVal.String(), "cep18_contract_hash_CLICKDevNet Test")
}

func Test_DefaultClient_QueryGlobalStateByBlockHeight_StoredAddressableEntity_NamedKey_Contract(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_named_key_contract.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	key := "entity-contract-a5cf5917505ef60a6f0df395dd19e86a0f075d00f2e6ce49f5aa0e18f6e26f5d"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, key, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.NamedKey)
	assert.NotEmpty(t, res.StoredValue.NamedKey.Name)
	assert.NotEmpty(t, res.StoredValue.NamedKey.NamedKey)
	assert.Equal(t, res.StoredValue.NamedKey.NamedKey.Key.URef.String(), "uref-c24ea2f7b569632dc0038bf87a8b9a4c720426fd177ea53f615f7723cd056202-007")
	assert.Equal(t, res.StoredValue.NamedKey.Name.StringVal.String(), "enable_mint_burn")
}

func Test_DefaultClient_QueryGlobalStateByBlock_GetWithdraw(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_withdraw.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	blockHash := "4ecd164a29a4e42d40e14ab823a107fe1869dbd3a309202d729589cfa85265ec"
	withdrawKey := "withdraw-b98aeebe049967c8273088fac9301978ac1edbd8bff07c5256e0957e52a0ccac"
	res, err := client.QueryGlobalStateByBlockHash(context.Background(), blockHash, withdrawKey, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.BlockHeader.BodyHash)
	assert.Equal(t, "500468846906", res.StoredValue.Withdraw[0].Amount.String())
}

func Test_DefaultClient_QueryGlobalStateByStateRoot_GetAccount(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_era.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	stateRootHash := "bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59"
	accountKey := "account-hash-e94daaff79c2ab8d9c31d9c3058d7d0a0dd31204a5638dc1451fa67b2e3fb88c"
	res, err := client.QueryGlobalStateByStateHash(context.Background(), &stateRootHash, accountKey, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.Account.AccountHash)
}

func Test_DefaultClient_QueryGlobalStateByStateRoot_GetAccount_WithEmptyStateRootHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_era.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	accountKey := "account-hash-e94daaff79c2ab8d9c31d9c3058d7d0a0dd31204a5638dc1451fa67b2e3fb88c"
	res, err := client.QueryGlobalStateByStateHash(context.Background(), nil, accountKey, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.Account.AccountHash)
}

func Test_DefaultClient_GetAccountInfoByBlochHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_account_info.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	pubKey, err := casper.NewPublicKey("01018525deae6091abccab6704a0fa44e12c495eec9e8fe6929862e1b75580e715")
	require.NoError(t, err)
	blockHash := "bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59"
	res, err := client.GetAccountInfoByBlochHash(context.Background(), blockHash, pubKey)
	require.NoError(t, err)
	assert.Equal(t, "account-hash-e94daaff79c2ab8d9c31d9c3058d7d0a0dd31204a5638dc1451fa67b2e3fb88c", res.Account.AccountHash.ToPrefixedString())
}

func Test_DefaultClient_GetAccountInfoByBlochHeight(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_account_info.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	pubKey, err := casper.NewPublicKey("01018525deae6091abccab6704a0fa44e12c495eec9e8fe6929862e1b75580e715")
	require.NoError(t, err)
	res, err := client.GetAccountInfoByBlochHeight(context.Background(), 185, pubKey)
	require.NoError(t, err)
	assert.Equal(t, "account-hash-e94daaff79c2ab8d9c31d9c3058d7d0a0dd31204a5638dc1451fa67b2e3fb88c", res.Account.AccountHash.ToPrefixedString())
}

func Test_DefaultClient_GetStateBalance(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_account_balance.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	hash := "fb9c42717769d72442ff17a5ff1574b4bc1c83aedf5992b14e4d071423f86240"
	result, err := client.GetBalanceByStateRootHash(
		context.Background(),
		"uref-7b12008bb757ee32caefb3f7a1f77d9f659ee7a4e21ad4950c4e0294000492eb-007",
		hash,
	)
	require.NoError(t, err)
	assert.Equal(t, "93000000000", result.BalanceValue.String())
}

func Test_DefaultClient_GetStateBalance_WithEmptyStateRootHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_account_balance.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetLatestBalance(
		context.Background(),
		"uref-7b12008bb757ee32caefb3f7a1f77d9f659ee7a4e21ad4950c4e0294000492eb-007",
	)
	require.NoError(t, err)
	assert.NotEmpty(t, result.BalanceValue)
}

func Test_DefaultClient_GetEraInfoLatest(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_era_info_latest.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	_, err := client.GetEraInfoLatest(
		context.Background(),
	)
	assert.NoError(t, err)
}

func Test_DefaultClient_GetEraInfoByBlockHeight(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_era_info_by_block.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetEraInfoByBlockHeight(
		context.Background(),
		1412462,
	)
	require.NoError(t, err)
	assert.Equal(t, "5dafbccc05cd3eb765ef9471a141877d8ffae306fb79c75fa4db46ab98bca370", result.EraSummary.BlockHash.ToHex())
}

func Test_DefaultClient_GetEraInfoByBlockHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_era_info_by_block.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetEraInfoByBlockHash(
		context.Background(),
		"5dafbccc05cd3eb765ef9471a141877d8ffae306fb79c75fa4db46ab98bca370",
	)
	require.NoError(t, err)
	assert.Equal(t, "5dafbccc05cd3eb765ef9471a141877d8ffae306fb79c75fa4db46ab98bca370", result.EraSummary.BlockHash.ToHex())
}

func Test_DefaultClient_GetValidatorChanges(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_validator_changes.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	res, err := client.GetValidatorChangesInfo(context.Background())
	require.NoError(t, err)
	assert.Equal(t, rpc.ValidatorStateAdded, res.Changes[0].StatusChanges[0].ValidatorState)
}

func Test_DefaultClient_GetBlockLatest(t *testing.T) {
	tests := []struct {
		filePath string
	}{
		{
			filePath: "../data/rpc_response/get_block_v2.json",
		},
		{
			filePath: "../data/rpc_response/get_block_v2_era_end.json",
		},
		{
			filePath: "../data/rpc_response/get_block_v1.json",
		},
	}
	for _, tt := range tests {
		t.Run("GetLatestBlock", func(t *testing.T) {
			server := SetupServer(t, tt.filePath)
			defer server.Close()
			client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
			result, err := client.GetLatestBlock(context.Background())
			require.NoError(t, err)
			assert.NotEmpty(t, result.APIVersion)
			assert.NotEmpty(t, result.Block.Hash)
			assert.NotEmpty(t, result.Block.Proposer)
			assert.NotEmpty(t, result.Block.Height)
			assert.NotEmpty(t, result.Block.ParentHash)
			assert.NotEmpty(t, result.Block.StateRootHash)
			assert.NotEmpty(t, result.Block.Timestamp)
			assert.NotEmpty(t, result.Block.EraID)
			assert.NotEmpty(t, result.Block.ProtocolVersion)
			assert.NotEmpty(t, result.Block.Proofs)
			assert.NotEmpty(t, result.Block.Transactions)
			assert.NotEmpty(t, result.GetRawJSON())
		})
	}
}

func Test_DefaultClient_GetReward(t *testing.T) {
	tests := []struct {
		filePath string
	}{
		{
			filePath: "../data/rpc_response/info_get_reward.json",
		},
	}

	pubKey, err := casper.NewPublicKey("0115394d1f395a87dfed4ab62bbfbc91b573bbb2bffb2c8ebb9c240c51d95bcc4d")
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run("GetValidatorRewardByEraID", func(t *testing.T) {
			server := SetupServer(t, tt.filePath)
			defer server.Close()
			client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
			result, err := client.GetValidatorRewardByEraID(context.Background(), pubKey, 100)
			require.NoError(t, err)
			assert.NotEmpty(t, result.APIVersion)
			assert.NotEmpty(t, result.EraID)
			assert.NotEmpty(t, result.RewardAmount)
			assert.Equal(t, result.RewardAmount.Value().Int64(), int64(62559062048560))
			assert.NotEmpty(t, result.GetRawJSON())
		})
	}
}

func Test_DefaultClient_GetEntity(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/state_get_entity_account_example.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))

	result, err := client.GetLatestEntity(context.Background(), rpc.EntityIdentifier{})
	require.NoError(t, err)
	assert.NotEmpty(t, result.Entity.AddressableEntity)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.EntityKind.Account)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.PackageHash)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.ByteCodeHash)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.AssociatedKeys)
}

func Test_DefaultClient_GetEntity_SystemKind(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/state_get_entity_system_example.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))

	result, err := client.GetLatestEntity(context.Background(), rpc.EntityIdentifier{})
	require.NoError(t, err)
	assert.NotEmpty(t, result.Entity.AddressableEntity)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.EntityKind.System)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.PackageHash)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.ByteCodeHash)
	assert.NotEmpty(t, result.Entity.AddressableEntity.EntryPoints)
	assert.NotEmpty(t, result.Entity.AddressableEntity.NamedKeys)
}

func Test_DefaultClient_GetEntity_SmartContractKind(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/state_get_entity_contract_example.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))

	result, err := client.GetLatestEntity(context.Background(), rpc.EntityIdentifier{})
	require.NoError(t, err)
	assert.NotEmpty(t, result.Entity.AddressableEntity)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.EntityKind.SmartContract)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.PackageHash)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.ByteCodeHash)
	assert.NotEmpty(t, result.Entity.AddressableEntity.Entity.MainPurse)
	assert.NotEmpty(t, result.Entity.AddressableEntity.EntryPoints)
	assert.NotEmpty(t, result.Entity.AddressableEntity.NamedKeys)
}

func Test_DefaultClient_GetBlockTransfersLatest_V2(t *testing.T) {
	tests := []struct {
		filePath string
	}{
		{
			filePath: "../data/rpc_response/get_block_transfers_v2.json",
		},
		{
			filePath: "../data/rpc_response/get_block_transfers_v2_old.json",
		},
		{
			filePath: "../data/rpc_response/get_block_transfers_v1.json",
		},
	}

	for _, tt := range tests {
		t.Run("GetBlockLatest", func(t *testing.T) {
			server := SetupServer(t, tt.filePath)
			defer server.Close()
			client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
			result, err := client.GetBlockTransfersByHeight(context.Background(), 3144186)
			require.NoError(t, err)
			assert.NotEmpty(t, result.BlockHash)
			assert.NotEmpty(t, result.Transfers[0])
			assert.NotEmpty(t, result.Transfers[0].TransactionHash)
			assert.NotEmpty(t, result.Transfers[0].From.AccountHash)
			assert.NotEmpty(t, result.Transfers[0].Source)
			assert.NotEmpty(t, result.Transfers[0].Amount)
		})
	}
}

func Test_DefaultClient_GetEraSummaryLatest(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_era_info_summary.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetEraSummaryLatest(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, result.EraSummary.StoredValue.EraInfo.SeigniorageAllocations)
}

func Test_DefaultClient_GetEraSummaryLatest_byHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_era_info_summary.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetEraSummaryByHash(context.Background(), "9bfa58709058935882a095ca6adf844b72a2ddf0f49b8575ef1ceda987452fb8")
	require.NoError(t, err)
	assert.NotEmpty(t, result.EraSummary.StoredValue.EraInfo.SeigniorageAllocations)
}

func Test_DefaultClient_GetEraSummaryLatest_byHeight(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_era_info_summary.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetEraSummaryByHeight(context.Background(), 1412462)
	require.NoError(t, err)
	assert.NotEmpty(t, result.EraSummary.StoredValue.EraInfo.SeigniorageAllocations)
}

func Test_DefaultClient_GetAuctionInfoLatest_Comaptible(t *testing.T) {
	tests := []struct {
		filePath string
		isV2     bool
	}{
		{
			filePath: "../data/rpc_response/get_auction_info.json",
		},
		{
			filePath: "../data/rpc_response/get_auction_info_v2.json",
			isV2:     true,
		},
	}
	for _, tt := range tests {
		t.Run("GetLatestAuctionInfo", func(t *testing.T) {
			server := SetupServer(t, tt.filePath)
			defer server.Close()
			client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
			result, err := client.GetLatestAuctionInfo(context.Background())
			require.NoError(t, err)
			assert.NotEmpty(t, result.Version)
			assert.NotEmpty(t, result.AuctionState.Bids)
			for _, bid := range result.AuctionState.Bids {
				assert.NotEmpty(t, bid.PublicKey)
				if tt.isV2 {
					assert.NotEmpty(t, bid.Bid.Delegators)
					assert.NotEmpty(t, bid.Bid.ValidatorPublicKey)
				}
				assert.NotEmpty(t, bid.Bid.BondingPurse)
				assert.NotEmpty(t, bid.Bid.StakedAmount)
			}
			assert.NotEmpty(t, result.AuctionState.StateRootHash)
			assert.NotEmpty(t, result.AuctionState.BlockHeight)
			assert.NotEmpty(t, result.AuctionState.EraValidators)
			assert.NotEmpty(t, result.GetRawJSON())
		})
	}
}

func Test_DefaultClient_GetAuctionInfoByHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_auction_info.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetAuctionInfoByHash(context.Background(), "5dafbccc05cd3eb765ef9471a141877d8ffae306fb79c75fa4db46ab98bca370")
	require.NoError(t, err)
	assert.NotEmpty(t, result.AuctionState.Bids)
}

func Test_DefaultClient_GetAuctionInfoByHeight(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_auction_info.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetAuctionInfoByHeight(context.Background(), 1412462)
	require.NoError(t, err)
	assert.NotEmpty(t, result.AuctionState.Bids)
}

func Test_DefaultClient_GetStateRootHashLatest(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_root_state_hash.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetStateRootHashLatest(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, result.StateRootHash)
}

func Test_DefaultClient_GetStateRootHashByHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_root_state_hash.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetStateRootHashByHash(context.Background(), "5dafbccc05cd3eb765ef9471a141877d8ffae306fb79c75fa4db46ab98bca370")
	require.NoError(t, err)
	assert.NotEmpty(t, result.StateRootHash)
}

func Test_DefaultClient_GetStateRootHashByHeight(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_root_state_hash.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetStateRootHashByHeight(context.Background(), 1412462)
	require.NoError(t, err)
	assert.NotEmpty(t, result.StateRootHash)
}

func Test_DefaultClient_GetStatus(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_status.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetStatus(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, result.ChainSpecName)
	assert.NotEmpty(t, result.LatestSwitchBlockHash)
}

func Test_DefaultClient_GetPeers(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_peers.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetPeers(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, result.Peers)
}

func Test_DefaultClient_QueryBalance_byPublicKey(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_balance.json")
	defer server.Close()
	pubKey, err := casper.NewPublicKey("0115394d1f395a87dfed4ab62bbfbc91b573bbb2bffb2c8ebb9c240c51d95bcc4d")
	require.NoError(t, err)
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.QueryLatestBalance(context.Background(), rpc.PurseIdentifier{
		MainPurseUnderPublicKey: &pubKey,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.Balance)
}

func Test_DefaultClient_GetChainspec(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_chainspec.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetChainspec(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, result.ChainspecBytes.ChainspecBytes)
}
