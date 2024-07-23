package key

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/v2/casper"
)

func Test_Key_Constructor_fromString(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		excepted string
	}{
		{"account prefix", "account-hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760", "account-hash-ee83ab5d92e183e2e92c1290a5979e2b7f7fac146c215de8042e2945bbad9760"},
		{"package prefix", "package-37b75a2ad99f92cc1c3a5b22b5d52a4c34dc2aae0340d6a75184fe38f314f751", "package-37b75a2ad99f92cc1c3a5b22b5d52a4c34dc2aae0340d6a75184fe38f314f751"},
		{"entity-contract prefix", "entity-contract-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330", "entity-contract-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330"},
		{"entity-system prefix", "entity-system-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330", "entity-system-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330"},
		{"entity-account prefix", "entity-account-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330", "entity-account-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330"},
		{"byte-code prefix", "byte-code-empty-0000000000000000000000000000000000000000000000000000000000000000", "byte-code-empty-0000000000000000000000000000000000000000000000000000000000000000"},
		{"byte-code prefix", "byte-code-v1-wasm-1b1e23596b8c901a65a13e9d314ca2fff440e69df42f226a76a9bbfbb90df1fa", "byte-code-v1-wasm-1b1e23596b8c901a65a13e9d314ca2fff440e69df42f226a76a9bbfbb90df1fa"},
		{"bid-addr- prefix", "bid-addr-01da3cd8cc4c8f34e7731583e67ddc211ff9b5c3f2c52640582415c2cce9315b2a", "bid-addr-01da3cd8cc4c8f34e7731583e67ddc211ff9b5c3f2c52640582415c2cce9315b2a"},
		{"bid-addr- prefix", "bid-addr-0494f1805abf61fac1b206d35773f1d1e71be2a162b58acd29fbca6ea5e8e73bedea00000000000000", "bid-addr-0494f1805abf61fac1b206d35773f1d1e71be2a162b58acd29fbca6ea5e8e73bedea00000000000000"},
		{"message-topic prefix", "message-topic-entity-contract-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-5721a6d9d7a9afe5dfdb35276fb823bed0f825350e4d865a5ec0110c380de4e1", "message-topic-entity-contract-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-5721a6d9d7a9afe5dfdb35276fb823bed0f825350e4d865a5ec0110c380de4e1"},
		{"message-system prefix", "message-topic-entity-system-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-5721a6d9d7a9afe5dfdb35276fb823bed0f825350e4d865a5ec0110c380de4e1", "message-topic-entity-system-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-5721a6d9d7a9afe5dfdb35276fb823bed0f825350e4d865a5ec0110c380de4e1"},
		{"message-account prefix", "message-topic-entity-account-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-5721a6d9d7a9afe5dfdb35276fb823bed0f825350e4d865a5ec0110c380de4e1", "message-topic-entity-account-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-5721a6d9d7a9afe5dfdb35276fb823bed0f825350e4d865a5ec0110c380de4e1"},
		{"message-entity prefix", "message-entity-contract-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-5721a6d9d7a9afe5dfdb35276fb823bed0f825350e4d865a5ec0110c380de4e1-0", "message-entity-contract-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-5721a6d9d7a9afe5dfdb35276fb823bed0f825350e4d865a5ec0110c380de4e1-0"},
		{"named-key prefix", "named-key-entity-contract-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-5f4fd818ad44d4ae056e151759a8585de97a1c7c4d53ceecac6631f9fbb39ab6", "named-key-entity-contract-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-5f4fd818ad44d4ae056e151759a8585de97a1c7c4d53ceecac6631f9fbb39ab6"},
		{"block-message-count prefix", "block-message-count-0000000000000000000000000000000000000000000000000000000000000000", "block-message-count-0000000000000000000000000000000000000000000000000000000000000000"},
		{"block-time- prefix", "block-time-0000000000000000000000000000000000000000000000000000000000000000", "block-time-0000000000000000000000000000000000000000000000000000000000000000"},
		{"entry-point-v1- prefix", "entry-point-v1-entity-system-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330", "entry-point-v1-entity-system-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330-55d4a6915291da12afded37fa5bc01f0803a2f0faf6acb7ec4c7ca6ab76f3330"},
		{"balance-hold prefix", "balance-hold-0032ac1e06b5538ac7b8ffe135c48f984754f7f4beade0b343f71c0bc4759eec35ab23787c8f010000", "balance-hold-0032ac1e06b5538ac7b8ffe135c48f984754f7f4beade0b343f71c0bc4759eec35ab23787c8f010000"},
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
