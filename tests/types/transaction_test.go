package types

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/casper"
	"github.com/make-software/casper-go-sdk/v2/types"
	"github.com/make-software/casper-go-sdk/v2/types/clvalue"
	"github.com/make-software/casper-go-sdk/v2/types/key"
)

func Test_TransactionSerialization_ModuleBytesSession(t *testing.T) {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile("../data/keys/docker-nctl-rc3-secret.pem")
	require.NoError(t, err)

	pubKey := keys.PublicKey()

	moduleBytes, err := os.ReadFile("../data/wasm/cep18-rc3.wasm")
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
}

func Test_TransactionSerialization_StoredTarget(t *testing.T) {
	keys, err := casper.NewED25519PrivateKeyFromPEMFile("../data/keys/docker-nctl-rc3-secret.pem")
	require.NoError(t, err)

	pubKey := keys.PublicKey()
	accountHash := pubKey.AccountHash()

	hash, err := key.NewHash("a5542d422cc7102165bde32f8c8aa460a81dc64105b03efbcd9c612a7721dadb")
	require.NoError(t, err)

	args := &types.Args{}
	args.AddArgument("name", *clvalue.NewCLString("Test")).
		AddArgument("symbol", *clvalue.NewCLString("test")).
		AddArgument("decimals", *clvalue.NewCLUint8(9)).
		AddArgument("total_supply", *clvalue.NewCLUInt256(big.NewInt(1_000_000_000_000_000))).
		AddArgument("events_mode", *clvalue.NewCLUint8(2)).
		AddArgument("enable_mint_burn", *clvalue.NewCLUint8(1))

	network := "casper-net-1"
	paymentAmount := uint64(100000000)

	createTransactionPayload := func(target types.TransactionTarget) (types.TransactionV1Payload, error) {
		return types.NewTransactionV1Payload(
			types.InitiatorAddr{AccountHash: &accountHash},
			types.Timestamp(time.Now().UTC()),
			1800000000000,
			network,
			types.PricingMode{
				Limited: &types.LimitedMode{
					GasPriceTolerance: 1,
					StandardPayment:   true,
					PaymentAmount:     paymentAmount,
				},
			},
			types.NewNamedArgs(args),
			target,
			types.TransactionEntryPoint{Call: &struct{}{}},
			types.TransactionScheduling{Standard: &struct{}{}},
		)
	}

	processTransaction := func(payload types.TransactionV1Payload) {
		require.NoError(t, err)

		transaction, err := types.MakeTransactionV1(payload)
		require.NoError(t, err)
		require.NoError(t, transaction.Sign(keys))
		require.NoError(t, transaction.Validate())
	}

	entryPoint := "apple"
	targets := []types.TransactionTarget{
		{
			Stored: &types.StoredTarget{
				ID:      types.TransactionInvocationTarget{ByHash: &hash},
				Runtime: types.NewVmCasperV1TransactionRuntime(),
			},
		},
		{
			Stored: &types.StoredTarget{
				ID:      types.TransactionInvocationTarget{ByName: &entryPoint},
				Runtime: types.NewVmCasperV1TransactionRuntime(),
			},
		},
		{
			Stored: &types.StoredTarget{
				ID: types.TransactionInvocationTarget{
					ByPackageHash: &types.ByPackageHashInvocationTarget{Addr: hash},
				},
				Runtime: types.NewVmCasperV1TransactionRuntime(),
			},
		},
		{
			Stored: &types.StoredTarget{
				ID: types.TransactionInvocationTarget{
					ByPackageHash: &types.ByPackageHashInvocationTarget{Addr: hash, Version: &[]uint32{1}[0]},
				},
				Runtime: types.NewVmCasperV1TransactionRuntime(),
			},
		},
		{
			Stored: &types.StoredTarget{
				ID: types.TransactionInvocationTarget{
					ByPackageName: &types.ByPackageNameInvocationTarget{Name: entryPoint},
				},
				Runtime: types.NewVmCasperV1TransactionRuntime(),
			},
		},
		{
			Stored: &types.StoredTarget{
				ID: types.TransactionInvocationTarget{
					ByPackageName: &types.ByPackageNameInvocationTarget{Name: entryPoint, Version: &[]uint32{1}[0]},
				},
				Runtime: types.NewVmCasperV1TransactionRuntime(),
			},
		},
	}

	for _, target := range targets {
		payload, err := createTransactionPayload(target)
		require.NoError(t, err)
		processTransaction(payload)
	}
}
