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

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/rpc"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
)

func Test_PutTransaction(t *testing.T) {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile("../../data/keys/docker-nctl-rc3-secret.pem")
	require.NoError(t, err)

	pubKey := keys.PublicKey()

	header := types.TransactionV1Header{
		ChainName: "casper-net-1",
		Timestamp: types.Timestamp(time.Now().UTC()),
		TTL:       18000000000,
		InitiatorAddr: types.InitiatorAddr{
			PublicKey: &pubKey,
		},
		PricingMode: types.PricingMode{
			Fixed: &types.FixedMode{
				GasPriceTolerance: 3,
			},
		},
	}

	moduleBytes, err := os.ReadFile("../../data/wasm/cep18-rc3.wasm")
	require.NoError(t, err)

	args := &types.Args{}
	args.AddArgument("name", *clvalue.NewCLString("Test")).
		AddArgument("symbol", *clvalue.NewCLString("test")).
		AddArgument("decimals", *clvalue.NewCLUint8(9)).
		AddArgument("total_supply", *clvalue.NewCLUInt256(big.NewInt(1_000_000_000_000_000))).
		AddArgument("events_mode", *clvalue.NewCLUint8(2)).
		AddArgument("enable_mint_burn", *clvalue.NewCLUint8(1))

	body := types.TransactionV1Body{
		Args: args,
		Target: types.TransactionTarget{
			Session: &types.SessionTarget{
				ModuleBytes: hex.EncodeToString(moduleBytes),
				Runtime:     types.TransactionRuntimeVmCasperV1,
			},
		},
		TransactionEntryPoint: types.TransactionEntryPoint{
			Call: &struct{}{},
		},
		TransactionScheduling: types.TransactionScheduling{
			Standard: &struct{}{},
		},
		TransactionCategory: 2,
	}

	transaction, err := types.MakeTransactionV1(header, body)
	err = transaction.Sign(keys)
	require.NoError(t, err)

	rpcClient := rpc.NewClient(rpc.NewHttpHandler("http://127.0.0.1:11101/rpc", http.DefaultClient))
	res, err := rpcClient.PutTransactionV1(context.Background(), *transaction)
	require.NoError(t, err)
	assert.NotEmpty(t, res.TransactionHash.TransactionV1)
	assert.NoError(t, transaction.Validate())

	log.Println("TransactionV1 submitted:", res.TransactionHash.TransactionV1)

	time.Sleep(time.Second * 10)
	transactionRes, err := rpcClient.GetTransactionByTransactionHash(context.Background(), res.TransactionHash.TransactionV1.ToHex())
	require.NoError(t, err)
	assert.NotEmpty(t, transactionRes.Transaction)
	assert.NotEmpty(t, transactionRes.ExecutionInfo)
}
