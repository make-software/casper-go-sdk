package key

type PrefixName = string

const (
	PrefixNameAccount                PrefixName = "account-hash-"
	PrefixNameHash                   PrefixName = "hash-"
	PrefixNameContractPackageWasm    PrefixName = "contract-package-wasm"
	PrefixNameContractPackage        PrefixName = "contract-package-"
	PrefixNameContractWasm           PrefixName = "contract-wasm-"
	PrefixNameContract               PrefixName = "contract-"
	PrefixNameURef                   PrefixName = "uref-"
	PrefixNameTransfer               PrefixName = "transfer-"
	PrefixNameDeployInfo             PrefixName = "deploy-"
	PrefixNameEraId                  PrefixName = "era-"
	PrefixNameBid                    PrefixName = "bid-"
	PrefixNameBalance                PrefixName = "balance-"
	PrefixNameWithdraw               PrefixName = "withdraw-"
	PrefixNameDictionary             PrefixName = "dictionary-"
	PrefixNameSystemContractRegistry PrefixName = "system-contract-registry-"
	PrefixNameUnbond                 PrefixName = "unbond-"
	PrefixNameChainspecRegistry      PrefixName = "chainspec-registry-"
	PrefixNameChecksumRegistry       PrefixName = "checksum-registry-"
)
