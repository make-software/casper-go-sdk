package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type IDValue struct {
	intValue   int
	strValue   string
	isIntValue bool
}

func NewIDFromString(value string) *IDValue {
	return &IDValue{
		strValue:   value,
		isIntValue: false,
	}
}

func NewIDFromInt(value int) *IDValue {
	return &IDValue{
		intValue:   value,
		isIntValue: true,
	}
}

func (id *IDValue) String() string {
	if id.isIntValue {
		return strconv.Itoa(id.intValue)
	}
	return id.strValue
}

func (id *IDValue) Int() int {
	if id.isIntValue {
		return id.intValue
	}
	// Try to convert string to int
	val, err := strconv.Atoi(id.strValue)
	if err == nil {
		return val
	}
	return 0
}

func (id *IDValue) MarshalJSON() ([]byte, error) {
	if id.isIntValue {
		return json.Marshal(id.intValue)
	}
	return json.Marshal(id.strValue)
}

func (id *IDValue) UnmarshalJSON(data []byte) error {
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		id.intValue = intVal
		id.isIntValue = true
		return nil
	}

	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		id.strValue = strVal
		id.isIntValue = false
		return nil
	}

	return fmt.Errorf("IDValue should be an int or string")
}
