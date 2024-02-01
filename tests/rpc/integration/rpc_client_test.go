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

func Test_DefaultClient_QueryStateByStateHash(t *testing.T) {
	accountKey := "account-hash-bf06bdb1616050cea5862333d1f4787718f1011c95574ba92378419eefeeee59"
	res, err := GetRpcClient().QueryGlobalStateByStateHash(context.Background(), nil, accountKey, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, res.StoredValue.Account.AccountHash)
}
