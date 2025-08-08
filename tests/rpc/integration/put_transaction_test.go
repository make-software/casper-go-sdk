//go:build integration
// +build integration

package integration

import (
	"context"
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
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

func Test_PutTransaction_ModuleBytesSession(t *testing.T) {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile("../../data/keys/docker-nctl-rc3-secret.pem")
	require.NoError(t, err)

	pubKey := keys.PublicKey()

	moduleBytes, err := os.ReadFile("../../data/wasm/cep18-rc3.wasm")
	require.NoError(t, err)

	args := &types.Args{}
	args.AddArgument("name", *clvalue.NewCLString("Test")).
		AddArgument("symbol", *clvalue.NewCLString("test")).
		AddArgument("decimals", *clvalue.NewCLUint8(9)).
		AddArgument("total_supply", *clvalue.NewCLUInt256(big.NewInt(1_000_000_000_000_000))).
		AddArgument("events_mode", *clvalue.NewCLUint8(2)).
		AddArgument("enable_mint_burn", *clvalue.NewCLUint8(1))

	payload, err := types.NewTransactionV1Payload(
		types.InitiatorAddr{
			PublicKey: &pubKey,
		},
		types.Timestamp(time.Now().UTC()),
		1800000000000,
		"casper-net-1",
		types.PricingMode{
			Limited: &types.LimitedMode{
				GasPriceTolerance: 1,
				StandardPayment:   true,
				PaymentAmount:     100000000,
			},
		},
		types.NewNamedArgs(args),
		types.TransactionTarget{
			Session: &types.SessionTarget{
				ModuleBytes:      moduleBytes,
				Runtime:          types.NewVmCasperV1TransactionRuntime(),
				IsInstallUpgrade: true,
			},
		},
		types.TransactionEntryPoint{
			Call: &struct{}{},
		},
		types.TransactionScheduling{
			Standard: &struct{}{},
		},
	)
	require.NoError(t, err)

	transaction, err := types.MakeTransactionV1(payload)
	require.NoError(t, err)

	require.NoError(t, transaction.Sign(keys))
	require.NoError(t, transaction.Validate())

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

func Test_PutTransaction_StorageTarget(t *testing.T) {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile("../../data/keys/docker-nctl-rc3-secret.pem")
	require.NoError(t, err)

	pubKey := keys.PublicKey()
	accountHash := pubKey.AccountHash()
	entryPoint := "apple"

	args := &types.Args{}
	args.AddArgument("name", *clvalue.NewCLString("Test")).
		AddArgument("symbol", *clvalue.NewCLString("test")).
		AddArgument("decimals", *clvalue.NewCLUint8(9)).
		AddArgument("total_supply", *clvalue.NewCLUInt256(big.NewInt(1_000_000_000_000_000))).
		AddArgument("events_mode", *clvalue.NewCLUint8(2)).
		AddArgument("enable_mint_burn", *clvalue.NewCLUint8(1))

	hash, err := key.NewHash("a5542d422cc7102165bde32f8c8aa460a81dc64105b03efbcd9c612a7721dadb")

	payload, err := types.NewTransactionV1Payload(
		types.InitiatorAddr{
			AccountHash: &accountHash,
		},
		types.Timestamp(time.Now().UTC()),
		1800000000000,
		"casper-net-1",
		types.PricingMode{
			Limited: &types.LimitedMode{
				GasPriceTolerance: 1,
				StandardPayment:   true,
				PaymentAmount:     100000000,
			},
		},
		types.NewNamedArgs(args),
		types.TransactionTarget{
			Stored: &types.StoredTarget{
				ID: types.TransactionInvocationTarget{
					ByHash: &hash,
				},
				Runtime: types.NewVmCasperV1TransactionRuntime(),
			},
		},
		types.TransactionEntryPoint{
			Custom: &entryPoint,
		},
		types.TransactionScheduling{
			Standard: &struct{}{},
		},
	)
	require.NoError(t, err)

	transaction, err := types.MakeTransactionV1(payload)
	require.NoError(t, err)

	require.NoError(t, transaction.Sign(keys))
	require.NoError(t, transaction.Validate())

	rpcClient := rpc.NewClient(rpc.NewHttpHandler("http://127.0.0.1:11101/rpc", http.DefaultClient))
	_, err = rpcClient.PutTransactionV1(context.Background(), *transaction)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such contract at hash")
}

func Test_PutTransaction_NativeTransfer(t *testing.T) {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile("../../data/keys/docker-nctl-rc3-secret.pem")
	require.NoError(t, err)

	pubKey := keys.PublicKey()

	target, err := casper.NewPublicKey("0106ed45915392c02b37136618372ac8dde8e0e3b8ee6190b2ca6db539b354ede4")
	require.NoError(t, err)

	args := &types.Args{}
	args.AddArgument("target", clvalue.NewCLPublicKey(target)).
		AddArgument("amount", *clvalue.NewCLUInt512(big.NewInt(2500000000)))

	payload, err := types.NewTransactionV1Payload(
		types.InitiatorAddr{
			PublicKey: &pubKey,
		},
		types.Timestamp(time.Now().UTC()),
		1800000000000,
		"casper-net-1",
		types.PricingMode{
			Limited: &types.LimitedMode{
				GasPriceTolerance: 1,
				StandardPayment:   true,
				PaymentAmount:     100000000,
			},
		},
		types.NewNamedArgs(args),
		types.TransactionTarget{
			Native: &struct{}{},
		},
		types.TransactionEntryPoint{
			Transfer: &struct{}{},
		},
		types.TransactionScheduling{
			Standard: &struct{}{},
		},
	)
	require.NoError(t, err)

	transaction, err := types.MakeTransactionV1(payload)
	require.NoError(t, err)

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
