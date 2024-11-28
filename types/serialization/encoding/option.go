package encoding

import "errors"

type Option[T any] struct {
	Some *T
	None bool
}

const (
	OptionNoneTag byte = 0
	OptionSomeTag byte = 1
)

type OptionFromBytesDecoder[T any, D FromBytes[T]] struct {
	Decoder D
}

func (o Option[T]) IsSome() bool {
	return o.Some != nil
}

func (o *OptionFromBytesDecoder[T, D]) FromBytes(data []byte) (Option[T], []byte, error) {
	if len(data) == 0 {
		return Option[T]{}, nil, errors.New("empty input")
	}

	tag := data[0]
	remaining := data[1:]

	switch tag {
	case OptionNoneTag:
		return Option[T]{None: true}, remaining, nil
	case OptionSomeTag:
		value, rem, err := o.Decoder.FromBytes(remaining)
		if err != nil {
			return Option[T]{}, nil, err
		}
		return Option[T]{Some: &value}, rem, nil
	default:
		return Option[T]{}, nil, errors.New("invalid tag")
	}
}
