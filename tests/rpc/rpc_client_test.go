package rpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/casper"
)

func SetupServer(t *testing.T, filePath string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
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

func Test_DefaultClient_GetStateItem_GetAccount(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_state_item_account.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))

	hash := "account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59"
	result, err := client.GetStateItem(
		context.Background(),
		"fb9c42717769d72442ff17a5ff1574b4bc1c83aedf5992b14e4d071423f86240",
		hash,
		[]string{},
	)
	require.NoError(t, err)
	assert.Equal(t, hash, result.StoredValue.Account.AccountHash.ToPrefixedString())
}

func Test_DefaultClient_GetStateBalance(t *testing.T) {
	server := SetupServer(t, "../data/rpc_response/get_account_balance.json")
	defer server.Close()
	client := casper.NewRPCClient(casper.NewRPCHandler(server.URL, http.DefaultClient))
	result, err := client.GetAccountBalance(
		context.Background(),
		"fb9c42717769d72442ff17a5ff1574b4bc1c83aedf5992b14e4d071423f86240",
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