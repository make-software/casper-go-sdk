package types

import (
	"encoding/json"
	"strings"
	"time"
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
