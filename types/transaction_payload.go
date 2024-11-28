package types

import (
	"github.com/make-software/casper-go-sdk/v2/types/serialization"
	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

const (
	InitiatorAddrFieldIndex uint16 = iota
	TimestampFieldIndex
	TtlFieldIndex
	ChainNameFieldIndex
	PricingModeFieldIndex
	FieldsFieldIndex
)

const (
	ArgsMapKey uint16 = iota
	TargetMapKey
	EntryPointMapKey
	SchedulingMapKey
)

type NamedArgs struct {
	Args *Args `json:"Named,omitempty"`
}

func NewNamedArgs(args *Args) NamedArgs {
	return NamedArgs{
		Args: args,
	}
}

func (n NamedArgs) SerializedLength() int {
	return n.Args.SerializedLength()
}

func (n NamedArgs) Bytes() ([]byte, error) {
	argsBytes, err := n.Args.Bytes()
	if err != nil {
		return nil, err
	}

	result := make([]byte, 1) // adding extra leading byte
	return append(result, argsBytes...), nil
}

type TransactionV1Fields struct {
	// binary representation of fields
	fields serialization.Fields `json:"-"`

	NamedArgs NamedArgs `json:"args,omitempty"`
	// Execution target of a Transaction.
	Target TransactionTarget `json:"target"`
	// Entry point of a Transaction.
	TransactionEntryPoint TransactionEntryPoint `json:"entry_point"`
	// Scheduling mode of a Transaction.
	TransactionScheduling TransactionScheduling `json:"scheduling"`
}

func NewTransactionV1Fields(
	namedArgs NamedArgs,
	target TransactionTarget,
	entryPoint TransactionEntryPoint,
	scheduling TransactionScheduling,
) (TransactionV1Fields, error) {
	fields := serialization.NewFields()

	if err := fields.AddField(ArgsMapKey, namedArgs); err != nil {
		return TransactionV1Fields{}, err
	}

	if err := fields.AddField(TargetMapKey, &target); err != nil {
		return TransactionV1Fields{}, err
	}

	if err := fields.AddField(EntryPointMapKey, &entryPoint); err != nil {
		return TransactionV1Fields{}, err
	}

	if err := fields.AddField(SchedulingMapKey, &scheduling); err != nil {
		return TransactionV1Fields{}, err
	}

	return TransactionV1Fields{
		fields:                fields,
		NamedArgs:             namedArgs,
		Target:                target,
		TransactionEntryPoint: entryPoint,
		TransactionScheduling: scheduling,
	}, nil
}

type TransactionV1Payload struct {
	// The address of the initiator of a TransactionV1.
	InitiatorAddr InitiatorAddr `json:"initiator_addr"`
	// `Timestamp` formatted as per RFC 3339
	Timestamp Timestamp `json:"timestamp"`
	// Duration of the `Deploy` in milliseconds (from timestamp).
	TTL Duration `json:"ttl"`
	// Chain name
	ChainName string `json:"chain_name"`
	// Pricing mode of a Transaction.
	PricingMode PricingMode `json:"pricing_mode"`

	Fields TransactionV1Fields `json:"fields"`
}

func (f *TransactionV1Fields) Bytes() ([]byte, error) {
	return f.fields.Bytes()
}

func (f *TransactionV1Fields) SerializedLength() int {
	return f.fields.SerializedLength()
}

func NewTransactionV1Payload(
	initiatorAddr InitiatorAddr,
	timestamp Timestamp,
	ttL Duration,
	chainName string,
	pricingMode PricingMode,
	args NamedArgs,
	target TransactionTarget,
	entryPoint TransactionEntryPoint,
	scheduling TransactionScheduling,
) (TransactionV1Payload, error) {
	transactionFields, err := NewTransactionV1Fields(args, target, entryPoint, scheduling)
	if err != nil {
		return TransactionV1Payload{}, err
	}

	return TransactionV1Payload{
		InitiatorAddr: initiatorAddr,
		Timestamp:     timestamp,
		TTL:           ttL,
		ChainName:     chainName,
		PricingMode:   pricingMode,
		Fields:        transactionFields,
	}, nil
}

func (d TransactionV1Payload) Bytes() ([]byte, error) {
	builder, err := serialization.NewCallTableSerializationEnvelopeBuilder(d.serializedFieldLengths())
	if err != nil {
		return nil, err
	}

	initiatorAddrBytes, err := d.InitiatorAddr.Bytes()
	if err != nil {
		return nil, err
	}

	if err = builder.AddField(InitiatorAddrFieldIndex, initiatorAddrBytes); err != nil {
		return nil, err
	}

	timestampBytes, err := d.Timestamp.Bytes()
	if err != nil {
		return nil, err
	}

	if err = builder.AddField(TimestampFieldIndex, timestampBytes); err != nil {
		return nil, err
	}

	ttlBytes, err := d.TTL.Bytes()
	if err != nil {
		return nil, err
	}

	if err = builder.AddField(TtlFieldIndex, ttlBytes); err != nil {
		return nil, err
	}

	chainNameBytes, err := encoding.NewStringToBytesEncoder(d.ChainName).Bytes()
	if err != nil {
		return nil, err
	}

	if err = builder.AddField(ChainNameFieldIndex, chainNameBytes); err != nil {
		return nil, err
	}

	pricingModeBytes, err := d.PricingMode.Bytes()
	if err != nil {
		return nil, err
	}

	if err = builder.AddField(PricingModeFieldIndex, pricingModeBytes); err != nil {
		return nil, err
	}

	fieldsBytes, err := d.Fields.Bytes()
	if err != nil {
		return nil, err
	}

	if err = builder.AddField(FieldsFieldIndex, fieldsBytes); err != nil {
		return nil, err
	}

	return builder.BinaryPayloadBytes()
}

func (d TransactionV1Payload) serializedFieldLengths() []int {
	return []int{
		d.InitiatorAddr.SerializedLength(),
		d.Timestamp.SerializedLength(),
		d.TTL.SerializedLength(),
		encoding.StringSerializedLength(d.ChainName),
		d.PricingMode.SerializedLength(),
		d.Fields.SerializedLength(),
	}
}
