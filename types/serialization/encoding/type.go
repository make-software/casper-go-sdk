package encoding

type FromBytes[T any] interface {
	FromBytes([]byte) (T, []byte, error)
}

type ToBytes interface {
	Bytes() ([]byte, error)
	SerializedLength() int
}
