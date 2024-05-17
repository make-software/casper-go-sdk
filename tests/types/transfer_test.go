package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_TransferV1_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/transfer/transfer_example.json")
	assert.NoError(t, err)

	var transfer types.TransferV1
	err = json.Unmarshal(fixture, &transfer)
	assert.NoError(t, err)

	result, err := json.Marshal(transfer)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}

func Test_TransferV2_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/transfer/transfer_example_v2.json")
	assert.NoError(t, err)

	var transfer types.TransferV2
	err = json.Unmarshal(fixture, &transfer)
	assert.NoError(t, err)

	result, err := json.Marshal(transfer)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}

func Test_TransferV1_WithEmptyReceiver_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/transfer/transfer_empty_receiver.json")
	assert.NoError(t, err)

	var transfer types.TransferV1
	err = json.Unmarshal(fixture, &transfer)
	assert.NoError(t, err)

	result, err := json.Marshal(transfer)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
