package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
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

func Test_DefaultClient_GetDeploy(t *testing.T) {
	server := SetupServer(t, "../data/deploy/get_raw_rpc_deploy.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	deployHash := "0009ea4441f4700325d9c38b0b6df415537596e1204abe4f6a94b6996aebf2f1"
	result, err := client.GetDeploy(context.Background(), deployHash)
	require.NoError(t, err)
	assert.Equal(t, deployHash, result.Deploy.Hash.ToHex())
}

func Test_DefaultClient_GetDeployFinalizedApproval(t *testing.T) {
	server := SetupServer(t, "../data/deploy/get_raw_rpc_deploy.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	deployHash := "0009ea4441f4700325d9c38b0b6df415537596e1204abe4f6a94b6996aebf2f1"
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

func Test_DefaultClient_QueryGlobalStateByBlockHeight_GetAccount(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/query_global_state_era.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	accountKey := "account-hash-e94daaff79c2ab8d9c31d9c3058d7d0a0dd31204a5638dc1451fa67b2e3fb88c"
	res, err := client.QueryGlobalStateByBlockHeight(context.Background(), 1000, accountKey, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.BlockHeader.BodyHash)
	assert.NotEmpty(t, res.StoredValue.Account.AccountHash)
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
	assert.Equal(t, "500468846906", strconv.Itoa(int(res.StoredValue.Withdraw[0].Amount)))
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
	result, err := client.GetAccountBalance(
		context.Background(),
		&hash,
		"uref-7b12008bb757ee32caefb3f7a1f77d9f659ee7a4e21ad4950c4e0294000492eb-007",
	)
	require.NoError(t, err)
	assert.NotEmpty(t, result.BalanceValue)
}

func Test_DefaultClient_GetStateBalance_WithEmptyStateRootHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_account_balance.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetAccountBalance(
		context.Background(),
		nil,
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
	server := SetupServer(t, "../data/rpc_response/get_block.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetBlockLatest(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, result.Block.Hash)
}

func Test_DefaultClient_GetBlockByHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_block.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetBlockByHash(context.Background(), "5dafbccc05cd3eb765ef9471a141877d8ffae306fb79c75fa4db46ab98bca370")
	require.NoError(t, err)
	assert.NotEmpty(t, result.Block.Hash)
}

func Test_DefaultClient_GetBlockByHeight(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_block.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetBlockByHeight(context.Background(), 185)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Block.Hash)
}

func Test_DefaultClient_GetBlockTransfersLatest(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_block_transfer.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetBlockTransfersLatest(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, result.BlockHash)
}

func Test_DefaultClient_GetBlockTransfersByHash(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_block_transfer.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetBlockTransfersByHash(context.Background(), "5dafbccc05cd3eb765ef9471a141877d8ffae306fb79c75fa4db46ab98bca370")
	require.NoError(t, err)
	assert.NotEmpty(t, result.BlockHash)
}

func Test_DefaultClient_GetBlockTransfersByHeight(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_block_transfer.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetBlockTransfersByHeight(context.Background(), 1412462)
	require.NoError(t, err)
	assert.NotEmpty(t, result.BlockHash)
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

func Test_DefaultClient_GetAuctionInfoLatest(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_auction_info.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetAuctionInfoLatest(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, result.AuctionState.Bids)
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
	result, err := client.QueryBalance(context.Background(), rpc.PurseIdentifier{
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
