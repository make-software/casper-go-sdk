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
	PricingMode PricingMode          `json:"pricing_mode"`
	Fields      serialization.Fields `json:"fields"`
}

func NewTransactionV1Payload(
	initiatorAddr InitiatorAddr,
	timestamp Timestamp,
	ttL Duration,
	chainName string,
	pricingMode PricingMode,
	args *Args,
	target TransactionTarget,
	entryPoint TransactionEntryPoint,
	scheduling TransactionScheduling,
) (TransactionV1Payload, error) {
	fields := serialization.NewFields()

	if err := fields.AddField(ArgsMapKey, args); err != nil {
		return TransactionV1Payload{}, err
	}
	if err := fields.AddField(TargetMapKey, &target); err != nil {
		return TransactionV1Payload{}, err
	}
	if err := fields.AddField(EntryPointMapKey, &entryPoint); err != nil {
		return TransactionV1Payload{}, err
	}
	if err := fields.AddField(SchedulingMapKey, &scheduling); err != nil {
		return TransactionV1Payload{}, err
	}
	return TransactionV1Payload{
		InitiatorAddr: initiatorAddr,
		Timestamp:     timestamp,
		TTL:           ttL,
		ChainName:     chainName,
		PricingMode:   pricingMode,
		Fields:        nil,
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

	pricingModeBytes, err := d.TTL.Bytes()
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
