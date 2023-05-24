package keypair

type (
	keyAlgorithm        byte
	KeyAlgorithmSetting struct {
		name string
	}
)

const (
	ED25519   keyAlgorithm = 1
	SECP256K1 keyAlgorithm = 2
)

var KeySettings = map[keyAlgorithm]KeyAlgorithmSetting{
	ED25519: {
		name: "ED25519",
	},
	SECP256K1: {
		name: "SECP256K1",
	},
}

func (a keyAlgorithm) String() string {
	return KeySettings[a].name
}

func (a keyAlgorithm) Byte() byte {
	return byte(a)
}
