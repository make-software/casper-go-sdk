//go:build integration
// +build integration

package integration

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
)

func GetRpcClient() rpc.Client {
	url, found := os.LookupEnv("NODE_URL")
	if !found {
		panic("NODE_URL env variable is not set")
	}
	return rpc.NewClient(rpc.NewHttpHandler(url, http.DefaultClient))
}

func Test_DefaultClient_GetAccountInfo_ByPublicKey(t *testing.T) {
	pubKey, err := casper.NewPublicKey("01018525deae6091abccab6704a0fa44e12c495eec9e8fe6929862e1b75580e715")
	require.NoError(t, err)
	res, err := GetRpcClient().GetAccountInfo(context.Background(), nil, rpc.AccountIdentifier{PublicKey: &pubKey})
	require.NoError(t, err)
	assert.Equal(t, "account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59", res.Account.AccountHash.ToPrefixedString())
}

func Test_DefaultClient_GetAccountInfo_ByAccountKey(t *testing.T) {
	accountKey, err := casper.NewAccountHash("account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59")
	require.NoError(t, err)
	res, err := GetRpcClient().GetAccountInfo(context.Background(), nil, rpc.AccountIdentifier{AccountHash: &accountKey})
	require.NoError(t, err)
	assert.Equal(t, "account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59", res.Account.AccountHash.ToPrefixedString())
}

func Test_DefaultClient_QueryBalanceDetails(t *testing.T) {
	pubKey, err := casper.NewPublicKey("0111bc2070a9af0f26f94b8549bffa5643ead0bc68eba3b1833039cfa2a9a8205d")
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		name       string
		identifier rpc.BalanceStateIdentifier
	}{
		{
			name: "ByBlock",
			identifier: rpc.BalanceStateIdentifier{
				Block: &rpc.BlockIdentifier{
					Height: &[]uint64{165520}[0],
				},
			},
		},
		{
			name: "ByBlockHash",
			identifier: rpc.BalanceStateIdentifier{
				Block: &rpc.BlockIdentifier{
					Hash: &[]string{"ba48f6c4211a98ec0db3e62e95133a5b3cbd521107112cc115ed0ad84bd1083f"}[0],
				},
			},
		},
		{
			name: "ByStateRoot",
			identifier: rpc.BalanceStateIdentifier{
				StateRoot: &rpc.StateRootInfo{
					StateRootHash: "281da1208effdfe4df196ddb83b862e4d11c5f58ddb1de4bfca4f8d43a51b6b4",
					Timestamp:     "2024-05-14T11:39:50.921Z",
				},
			},
		},
	}

	entityAddr := "entity-account-" + pubKey.AccountHash().ToHex()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetRpcClient().QueryBalanceDetails(ctx,
				rpc.PurseIdentifier{
					MainPurseUnderEntityAddr: &entityAddr,
				},
				tt.identifier,
			)
			require.NoError(t, err)
			assert.NotEmpty(t, result.AvailableBalance)
			assert.NotEmpty(t, result.TotalBalance)
			assert.NotEmpty(t, result.TotalBalanceProof)
			assert.Empty(t, result.Holds)
		})
	}
}
func Test_DefaultClient_QueryStateByStateHash(t *testing.T) {
	accountKey := "account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59"
	res, err := GetRpcClient().QueryGlobalStateByStateHash(context.Background(), nil, accountKey, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.Account.AccountHash)
}
