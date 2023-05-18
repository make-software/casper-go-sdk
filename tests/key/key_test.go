package key

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/casper"
)

func Test_Key_Constructor_fromString(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		excepted string
	}{
		{"account prefix", "account-hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "account-hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"account bool prefix", "00ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "account-hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"account key", "Key::Account(00ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "account-hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"hash", "ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"hash prefix", "hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"hash key", "Key::Hash(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"uref", "uref-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760-007", "uref-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760-007"},
		{"uref key", "Key::URef(uref-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760-007)", "uref-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760-007"},
		{"transfer", "transfer-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "transfer-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"transfer key", "Key::Transfer(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "transfer-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"deploy-info", "deploy-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "deploy-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"deploy-info key", "Key::Deploy(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "deploy-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"era-id", "era-123", "era-123"},
		{"era-id key", "Key::Era(123)", "era-123"},
		{"balance", "balance-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "balance-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"balance key", "Key::Balance(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "balance-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"bid", "bid-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "bid-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"balance key", "Key::Bid(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "bid-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"withdraw", "withdraw-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "withdraw-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"withdraw key", "Key::Withdraw(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "withdraw-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"dictionary", "dictionary-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "dictionary-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"dictionary key", "Key::Dictionary(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "dictionary-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"SystemContractRegistry", "system-contract-registry-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "system-contract-registry-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"SystemContractRegistry key", "Key::SystemContractRegistry(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "system-contract-registry-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"unbond", "unbond-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "unbond-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"unbond key", "Key::Unbond(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "unbond-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"chainspec-registry", "chainspec-registry-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "chainspec-registry-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"chainspec-registry key", "Key::ChainspecRegistry(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "chainspec-registry-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"checksum-registry", "checksum-registry-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "checksum-registry-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"checksum-registry key", "Key::ChecksumRegistry(ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760)", "checksum-registry-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := casper.NewKey(test.source)
			assert.NoError(t, err)
			assert.Equal(t, test.excepted, result.String())
		})
	}
}
