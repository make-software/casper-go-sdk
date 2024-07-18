package types

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/blake2b"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/types"
	"github.com/make-software/casper-go-sdk/types/clvalue"
	"github.com/make-software/casper-go-sdk/types/key"
)

func Test_MakeDeploy(t *testing.T) {
	expectedJson := `{
      "hash": "0523104c55a708a833fa7df6d56a4dffa3d1d03885c7ddf788f70a2da2527f2e",
      "header": {
        "account": "0146c64d0506c486f2b19f9cf73479fba550f33227b6ec1c12e58b437d2680e96d",
        "timestamp": "2023-05-08T21:33:00.268Z",
        "ttl": "30m",
        "gas_price": 1,
        "body_hash": "1f3f81a43a6de45b1fa71d4d2c88271ddb10afa570376aa1fef862fc0bbdc5ce",
        "dependencies": [],
        "chain_name": "casper-net-1"
      },
      "payment": {
        "ModuleBytes": {
          "args": [
            [
              "amount",
              {
                "bytes": "0400286bee",
                "cl_type": "U512"
              }
            ]
          ],
          "module_bytes": ""
        }
      },
      "session": {
        "ModuleBytes": {
          "args": [
            [
              "target",
              {
                "bytes": "80cef79f0451fdfb21084aaab8a4811e27dfc262c970d77675e3cad5394ef1f7",
                "cl_type": {
                  "ByteArray": 32
                }
              }
            ],
            [
              "amount",
              {
                "bytes": "050010a5d4e8",
                "cl_type": "U512"
              }
            ]
          ],
          "module_bytes": "74657374"
        }
      },
      "approvals": [
        {
          "signer": "0146c64d0506c486f2b19f9cf73479fba550f33227b6ec1c12e58b437d2680e96d",
          "signature": "011ab325948bda6029b566263dcdc4837fa7487bc6e177a9480be72dc5a6dd06578fc39d2189ecaaa4d8c802bd3f7e09a591903cd29b2bac866576cbdfb2dbee0c"
        }
      ]
    }`
	associatedAccount, err := key.NewAccountHash("80cef79f0451fdfb21084aaab8a4811e27dfc262c970d77675e3cad5394ef1f7")
	require.NoError(t, err)
	privateKey, err := casper.NewED25519PrivateKeyFromPEMFile("../data/keys/secret_key.pem")
	require.NoError(t, err)
	header := types.DefaultDeployHeader()
	header.ChainName = "casper-net-1"
	header.Account = privateKey.PublicKey()
	dateTime, err := time.Parse("2006-01-02T15:04:05.999Z", "2023-05-08T21:33:00.268Z")
	require.NoError(t, err)
	header.Timestamp = types.Timestamp(dateTime)
	payment := types.StandardPayment(big.NewInt(4000000000))

	moduleBytes, err := os.ReadFile("../data/wasm/empty.wasm")
	require.NoError(t, err)

	args := &types.Args{}
	args.AddArgument("target", clvalue.NewCLByteArray(associatedAccount.Hash.Bytes())).AddArgument("amount", *clvalue.NewCLUInt512(big.NewInt(1000000000000)))
	session := types.ExecutableDeployItem{
		ModuleBytes: &types.ModuleBytes{
			ModuleBytes: hex.EncodeToString(moduleBytes),
			Args:        args,
		},
	}

	deploy, err := types.MakeDeploy(header, payment, session)
	require.NoError(t, err)
	err = deploy.Sign(privateKey)
	require.NoError(t, err)
	actual, err := json.Marshal(deploy)
	require.NoError(t, err)
	require.JSONEq(t, expectedJson, string(actual), "deploy doesn't equal")
}

func Test_Payment_ToJsonSerialization(t *testing.T) {
	expectedJson := `{
        "ModuleBytes": {
          "args": [
            [
              "amount",
              {
                "bytes": "0400286bee",
                "cl_type": "U512"
              }
            ]
          ],
          "module_bytes": ""
        }}`
	payment := types.StandardPayment(big.NewInt(4000000000))
	actual, err := json.Marshal(payment)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJson, string(actual), "payment serialization is wrong")
}

func Test_Session_ToJsonSerialization(t *testing.T) {
	expectedJson := `{
        "ModuleBytes": {
          "args": [
            [
              "target",
              {
                "bytes": "80cef79f0451fdfb21084aaab8a4811e27dfc262c970d77675e3cad5394ef1f7",
                "cl_type": {
                  "ByteArray": 32
                }
              }
            ],
            [
              "amount",
              {
                "bytes": "050010a5d4e8",
                "cl_type": "U512"
              }
            ]
          ],
          "module_bytes": "74657374"
        }
      }`

	session := createTestSession(t)
	actual, err := json.Marshal(session)
	require.NoError(t, err)
	assert.JSONEq(t, expectedJson, string(actual), "payment serialization is wrong")
}

func createTestSession(t *testing.T) types.ExecutableDeployItem {
	associatedAccount, err := key.NewAccountHash("80cef79f0451fdfb21084aaab8a4811e27dfc262c970d77675e3cad5394ef1f7")
	require.NoError(t, err)
	moduleBytes, err := os.ReadFile("../data/wasm/empty.wasm")
	require.NoError(t, err)

	args := &types.Args{}
	args.AddArgument("target", clvalue.NewCLByteArray(associatedAccount.Hash.Bytes())).
		AddArgument("amount", *clvalue.NewCLUInt512(big.NewInt(1000000000000)))
	session := types.ExecutableDeployItem{
		ModuleBytes: &types.ModuleBytes{
			ModuleBytes: hex.EncodeToString(moduleBytes),
			Args:        args,
		},
	}
	return session
}

func Test_BodyHash_calculate(t *testing.T) {
	expectedHash := "1f3f81a43a6de45b1fa71d4d2c88271ddb10afa570376aa1fef862fc0bbdc5ce"
	payment := types.StandardPayment(big.NewInt(4000000000))
	session := createTestSession(t)
	paymentBytes, err := payment.Bytes()
	require.NoError(t, err)
	sessionBytes, err := session.Bytes()
	require.NoError(t, err)
	serializedBody := append(paymentBytes, sessionBytes...)
	bodyHash := blake2b.Sum256(serializedBody)
	assert.Equal(t, expectedHash, key.Hash(bodyHash).ToHex())
}

func Test_Payment_ToBytesSerialization(t *testing.T) {
	expectedHex := "00000000000100000006000000616d6f756e74050000000400286bee08"
	payment := types.StandardPayment(big.NewInt(4000000000))
	bytes, err := payment.Bytes()
	require.NoError(t, err)
	assert.Equal(t, expectedHex, hex.EncodeToString(bytes))
}

func Test_Session_ToBytesSerialization(t *testing.T) {
	expectedHex := "00040000007465737402000000060000007461726765742000000080cef79f0451fdfb21084aaab8a4811e27dfc262c970d77675e3cad5394ef1f70f2000000006000000616d6f756e7406000000050010a5d4e808"
	session := createTestSession(t)
	bytes, err := session.Bytes()
	require.NoError(t, err)
	assert.Equal(t, expectedHex, hex.EncodeToString(bytes))
}

func Test_BodyHashCalculate_WithStoredVersionedContract(t *testing.T) {
	expectedHash := "3b88d9461f8f16e6975db372fffaa4198bbd5f818fb95c262eecdf3d22b77918"

	paymentArgs := &types.Args{}
	paymentArgs.AddArgument("amount", *clvalue.NewCLUInt512(big.NewInt(3000000000)))
	payment := types.ExecutableDeployItem{
		ModuleBytes: &types.ModuleBytes{
			ModuleBytes: "",
			Args:        paymentArgs,
		},
	}
	paymentBytes, err := payment.Bytes()
	require.NoError(t, err)

	sessionArgs := &types.Args{}
	sessionArgs.AddArgument("amount", *clvalue.NewCLUInt256(big.NewInt(2500000000)))
	contractHash, err := key.NewContract("8ff7a1c49017400013dcf78305343fa07c31b04292b7928845ed59764e1ee512")
	require.NoError(t, err)
	varVal := json.Number("2")
	session := types.ExecutableDeployItem{
		StoredVersionedContractByHash: &types.StoredVersionedContractByHash{
			Hash:       contractHash,
			EntryPoint: "get_message",
			Version:    &varVal,
			Args:       sessionArgs,
		},
	}

	sessionBytes, err := session.Bytes()
	require.NoError(t, err)

	serializedBody := append(paymentBytes, sessionBytes...)
	bodyHash := blake2b.Sum256(serializedBody)
	assert.Equal(t, expectedHash, key.Hash(bodyHash).ToHex())
}
