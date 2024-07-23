package sse

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/sse"
)

func Test_RawEvent_ParseAsDeployAcceptedEvent(t *testing.T) {
	data, err := os.ReadFile("../data/sse/deploy_accepted_event.json")
	require.NoError(t, err)
	rawEvent := sse.RawEvent{
		Data: data,
	}
	res, err := rawEvent.ParseAsDeployAcceptedEvent()
	require.NoError(t, err)
	assert.Equal(t, "99483863a391510b8d3447dd5cfc446b42d65e598672d569abc4cdded85b81e6", res.DeployAccepted.Hash.ToHex())
}

func Test_RawEvent_ParseAsFinalitySignatureEvent(t *testing.T) {
	data, err := os.ReadFile("../data/sse/finality_signature_event.json")
	require.NoError(t, err)
	rawEvent := sse.RawEvent{
		Data: data,
	}
	res, err := rawEvent.ParseAsFinalitySignatureEvent()
	require.NoError(t, err)
	assert.Equal(t, "abbcdc782a18a9ba31826b07c838a69a6b790c8b36a0fd5f0818f757834d82f5", res.FinalitySignature.BlockHash.ToHex())
}

func Test_RawEvent_ParseAsApiVersionEvent(t *testing.T) {
	data, err := os.ReadFile("../data/sse/api_version_event.json")
	require.NoError(t, err)
	rawEvent := sse.RawEvent{
		Data: data,
	}
	res, err := rawEvent.ParseAsAPIVersionEvent()
	require.NoError(t, err)
	assert.Equal(t, "1.0.0", res.APIVersion)
}

func Test_RawEvent_ParseAsBlockAddedEvent(t *testing.T) {
	tests := []struct {
		filePath     string
		expectedHash string
	}{
		{
			filePath:     "../data/sse/block_added_event.json",
			expectedHash: "5809c6aacc3ac0573a67677743f4cb93cd487ade1c5132c1f806f75b6248f35f",
		},
		{
			filePath:     "../data/sse/block_added_event_v2.json",
			expectedHash: "cb64adaae660d227b7d7579039a75e21f1022e2c044a5ed0ce0beefb12b95758",
		},
	}

	for _, tc := range tests {
		t.Run(tc.filePath, func(t *testing.T) {
			data, err := os.ReadFile(tc.filePath)
			require.NoError(t, err)

			rawEvent := sse.RawEvent{
				Data: data,
			}

			res, err := rawEvent.ParseAsBlockAddedEvent()
			require.NoError(t, err)
			assert.Equal(t, tc.expectedHash, res.BlockAdded.Block.Hash.ToHex())
			assert.NotEmpty(t, res.BlockAdded.Block.Proposer)
			assert.NotEmpty(t, res.BlockAdded.Block.Height)
			assert.NotEmpty(t, res.BlockAdded.Block.EraID)
			assert.NotEmpty(t, res.BlockAdded.Block.StateRootHash)
			assert.NotEmpty(t, res.BlockAdded.Block.RandomBit)
			assert.NotEmpty(t, res.BlockAdded.Block.ParentHash)
			assert.NotEmpty(t, res.BlockAdded.Block.AccumulatedSeed)
			assert.NotEmpty(t, res.BlockAdded.Block.CurrentGasPrice)
			assert.NotEmpty(t, res.BlockAdded.Block.ProtocolVersion)
		})
	}
}

func Test_RawEvent_FinalitySignature(t *testing.T) {
	tests := []struct {
		filePath     string
		expectedHash string
		isV2         bool
	}{
		{
			filePath:     "../data/sse/finality_signature_event.json",
			expectedHash: "abbcdc782a18a9ba31826b07c838a69a6b790c8b36a0fd5f0818f757834d82f5",
		},
		{
			filePath:     "../data/sse/finality_signature_event_v2.json",
			expectedHash: "5122b986d344e4cd0ab5a2df1b3bb398de9c557a2e463391b58de82626154428",
			isV2:         true,
		},
		{
			filePath:     "../data/sse/finality_signature_event_v2_old.json",
			expectedHash: "abbcdc782a18a9ba31826b07c838a69a6b790c8b36a0fd5f0818f757834d82f5",
		},
	}

	for _, tc := range tests {
		t.Run(tc.filePath, func(t *testing.T) {
			data, err := os.ReadFile(tc.filePath)
			require.NoError(t, err)

			rawEvent := sse.RawEvent{
				Data: data,
			}

			res, err := rawEvent.ParseAsFinalitySignatureEvent()
			require.NoError(t, err)
			assert.Equal(t, tc.expectedHash, res.FinalitySignature.BlockHash.ToHex())
			assert.NotEmpty(t, res.FinalitySignature.Signature)
			assert.NotEmpty(t, res.FinalitySignature.PublicKey)
			assert.NotEmpty(t, res.FinalitySignature.EraID)
		})
	}
}

func Test_RawEvent_ParseAsDeployProcessedEvent(t *testing.T) {
	data, err := os.ReadFile("../data/sse/deploy_processed_event.json")
	require.NoError(t, err)
	rawEvent := sse.RawEvent{
		Data: data,
	}
	res, err := rawEvent.ParseAsDeployProcessedEvent()
	require.NoError(t, err)
	assert.Equal(t, "f19e3b63678ca5aa9fa8b30377275c83f8c1a041902b38ce7f4de50f02dbf396", res.DeployProcessed.BlockHash.ToHex())
}

func Test_RawEvent_ParseAsTransactionExpiredEvent(t *testing.T) {
	tests := []struct {
		filePath      string
		isTransaction bool
	}{
		{
			filePath: "../data/sse/deploy_expired_event.json",
		},
		{
			filePath:      "../data/sse/transaction_expired_event.json",
			isTransaction: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.filePath, func(t *testing.T) {
			data, err := os.ReadFile(tc.filePath)
			require.NoError(t, err)

			rawEvent := sse.RawEvent{
				Data: data,
			}

			res, err := rawEvent.ParseAsTransactionExpiredEvent()
			require.NoError(t, err)
			require.NotEmpty(t, res.TransactionExpiredPayload.TransactionHash)

			if tc.isTransaction {
				require.NotEmpty(t, res.TransactionExpiredPayload.TransactionHash.TransactionV1)
			} else {
				require.NotEmpty(t, res.TransactionExpiredPayload.TransactionHash.Deploy)
			}
		})
	}
}

func Test_RawEvent_ParseAsTransactionAcceptedEvent(t *testing.T) {
	tests := []struct {
		filePath      string
		isTransaction bool
	}{
		{
			filePath: "../data/sse/deploy_accepted_event.json",
		},
		{
			filePath:      "../data/sse/transaction_accepted_event.json",
			isTransaction: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.filePath, func(t *testing.T) {
			data, err := os.ReadFile(tc.filePath)
			require.NoError(t, err)

			rawEvent := sse.RawEvent{
				Data: data,
			}

			res, err := rawEvent.ParseAsTransactionAcceptedEvent()
			require.NoError(t, err)

			if tc.isTransaction {
				require.NotEmpty(t, res.TransactionAcceptedPayload.Transaction.GetTransactionV1())
			} else {
				require.NotEmpty(t, res.TransactionAcceptedPayload.Transaction.GetDeploy())
			}

			require.NotEmpty(t, res.TransactionAcceptedPayload.Transaction)
			require.NotEmpty(t, res.TransactionAcceptedPayload.Transaction.Header)
			require.NotEmpty(t, res.TransactionAcceptedPayload.Transaction.Body)
		})
	}
}

func Test_RawEvent_ParseAsTransactionProcessedEvent(t *testing.T) {
	tests := []struct {
		filePath      string
		isTransaction bool
	}{
		{
			filePath: "../data/sse/deploy_processed_event.json",
		},
		{
			filePath:      "../data/sse/transaction_processed_event.json",
			isTransaction: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.filePath, func(t *testing.T) {
			data, err := os.ReadFile(tc.filePath)
			require.NoError(t, err)

			rawEvent := sse.RawEvent{
				Data: data,
			}

			res, err := rawEvent.ParseAsTransactionProcessedEvent()
			require.NoError(t, err)
			require.NotEmpty(t, res.TransactionProcessedPayload.TransactionHash)

			if tc.isTransaction {
				require.NotEmpty(t, res.TransactionProcessedPayload.TransactionHash.TransactionV1)
				require.NotEmpty(t, res.TransactionProcessedPayload.Messages)
			} else {
				require.NotEmpty(t, res.TransactionProcessedPayload.TransactionHash.Deploy)
			}

			require.NotEmpty(t, res.TransactionProcessedPayload.ExecutionResult)
			require.NotEmpty(t, res.TransactionProcessedPayload.TTL)
			require.NotEmpty(t, res.TransactionProcessedPayload.Timestamp)
			require.NotEmpty(t, res.TransactionProcessedPayload.BlockHash)
		})
	}
}

func Test_RawEvent_ParseAsFaultEvent(t *testing.T) {
	data, err := os.ReadFile("../data/sse/fault_event.json")
	require.NoError(t, err)
	rawEvent := sse.RawEvent{
		Data: data,
	}
	res, err := rawEvent.ParseAsFaultEvent()
	require.NoError(t, err)
	assert.Equal(t, "012fa85eb06279da42e68530e1116be04bfd2aaa5ed8d63401ebff4d9153a609a9", res.Fault.PublicKey.ToHex())
}

func Test_RawEvent_ParseAsStepEvent(t *testing.T) {
	data, err := os.ReadFile("../data/sse/step_event.json")
	require.NoError(t, err)
	rawEvent := sse.RawEvent{
		Data: data,
	}
	res, err := rawEvent.ParseAsStepEvent()
	require.NoError(t, err)
	assert.False(t, res.Step.ExecutionEffect.Transforms[0].Transform.IsWriteTransfer())
}

func Test_RawEvent_ParseAndMarshalStepEvent(t *testing.T) {
	data, err := os.ReadFile("../data/sse/step_event.json")
	require.NoError(t, err)
	rawEvent := sse.RawEvent{
		Data: data,
	}
	res, err := rawEvent.ParseAsStepEvent()
	require.NoError(t, err)
	_, err = json.Marshal(res)
	assert.NoError(t, err)
}
