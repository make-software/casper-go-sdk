package types

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/make-software/casper-go-sdk/v2/types/serialization/encoding"
)

type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).UTC().Format("2006-01-02T15:04:05.999Z"))
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var dataString string

	if err := json.Unmarshal(data, &dataString); err != nil {
		return err
	}

	parse, err := time.Parse("2006-01-02T15:04:05.999Z", dataString)
	if err != nil {
		return err
	}

	*t = Timestamp(parse)
	return nil
}

func (t *Timestamp) ToTime() time.Time {
	return time.Time(*t)
}

func (d Timestamp) Bytes() ([]byte, error) {
	return encoding.NewU64ToBytesEncoder(uint64(d.ToTime().Unix())).Bytes()
}

func (d Timestamp) SerializedLength() int {
	return encoding.U64SerializedLength
}

type TimestampFromBytesDecoder struct{}

func (addr *TimestampFromBytesDecoder) FromBytes(bytes []byte) (*Timestamp, []byte, error) {
	u64Decoder := encoding.U64FromBytesDecoder{}
	seconds, remainder, err := u64Decoder.FromBytes(bytes)
	if err != nil {
		return nil, nil, err
	}

	t := time.Unix(int64(seconds), 0).UTC()
	timestamp := Timestamp(t)

	return &timestamp, remainder, nil
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	s := time.Duration(d).String()
	if s == "24h0m0s" {
		s = "1day"
	}
	if strings.HasSuffix(s, "h0m0s") {
		s = strings.TrimSuffix(s, "0m0s")
	}
	if strings.HasSuffix(s, "m0s") {
		s = strings.TrimSuffix(s, "0s")
	}
	return json.Marshal(s)
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var dataString string

	if err := json.Unmarshal(data, &dataString); err != nil {
		return err
	}
	dataString = strings.ReplaceAll(dataString, " ", "")
	if dataString == "1day" {
		dataString = "24h"
	}
	duration, err := time.ParseDuration(dataString)
	if err != nil {
		return err
	}

	*d = Duration(duration)

	return nil
}

func (d Duration) Bytes() ([]byte, error) {
	return encoding.NewU64ToBytesEncoder(uint64(d)).Bytes()
}

type DurationFromBytesDecoder struct{}

func (addr *DurationFromBytesDecoder) FromBytes(bytes []byte) (*Duration, []byte, error) {
	u64Decoder := encoding.U64FromBytesDecoder{}
	raw, remainder, err := u64Decoder.FromBytes(bytes)
	if err != nil {
		return nil, nil, err
	}

	t := time.Duration(raw)
	duration := Duration(t)

	return &duration, remainder, nil
}

func (d Duration) SerializedLength() int {
	return encoding.U64SerializedLength
}
