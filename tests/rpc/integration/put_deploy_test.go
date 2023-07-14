//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
)

func Test_PutDeploy(t *testing.T) {
	facetKeys, err := casper.NewED25519PrivateKeyFromPEMFile("../../data/keys/docker-nctl-secret.pem")
	require.NoError(t, err)
	require.NoError(t, err)
	header := types.DefaultHeader()
	header.ChainName = "casper-net-1"
	header.Account = facetKeys.PublicKey()
	require.NoError(t, err)
	header.Timestamp = types.Timestamp(time.Now())
	payment := types.StandardPayment(big.NewInt(4000000000))

	moduleBytes, err := os.ReadFile("../../data/wasm/faucet.wasm")
	require.NoError(t, err)

	args := &types.Args{}
	args.AddArgument("target", clvalue.NewCLByteArray(facetKeys.PublicKey().AccountHash().Hash.Bytes())).AddArgument("amount", *clvalue.NewCLUInt512(big.NewInt(1000000000000)))
	session := types.ExecutableDeployItem{
		ModuleBytes: &types.ModuleBytes{
			ModuleBytes: hex.EncodeToString(moduleBytes),
			Args:        args,
		},
	}

	deploy, err := types.MakeDeploy(header, payment, session)
	err = deploy.SignDeploy(facetKeys)
	require.NoError(t, err)

	rpcClient := rpc.NewClient(rpc.NewHttpHandler("http://127.0.0.1:11101/rpc", http.DefaultClient))
	res, err := rpcClient.PutDeploy(context.Background(), *deploy)
	log.Println(deploy.Hash.ToHex())
	require.NoError(t, err)
	assert.NotEmpty(t, res.DeployHash)
}
