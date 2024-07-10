package types

// ByteCodeKind The type of Byte code.
type ByteCodeKind string

func (b ByteCodeKind) IsEmpty() bool {
	return b == "Empty"
}

func (b ByteCodeKind) IsV1CasperWasm() bool {
	return b == "V1CasperWasm"
}

// ByteCode A container for contract's Wasm bytes.
type ByteCode struct {
	Kind ByteCodeKind `json:"kind"`
	// Bytes array representation of underlying data
	Bytes HexBytes `json:"bytes"`
}
