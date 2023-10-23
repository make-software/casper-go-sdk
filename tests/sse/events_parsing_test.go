package sse

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/sse"
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
	data, err := os.ReadFile("../data/sse/block_added_event.json")
	require.NoError(t, err)
	rawEvent := sse.RawEvent{
		Data: data,
	}
	res, err := rawEvent.ParseAsBlockAddedEvent()
	require.NoError(t, err)
	assert.Equal(t, "5809c6aacc3ac0573a67677743f4cb93cd487ade1c5132c1f806f75b6248f35f", res.BlockAdded.Block.Hash.ToHex())
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

func Test_RawEvent_ParseAsDeployExpiredEvent(t *testing.T) {
	data, err := os.ReadFile("../data/sse/deploy_expired_event.json")
	require.NoError(t, err)
	rawEvent := sse.RawEvent{
		Data: data,
	}
	res, err := rawEvent.ParseAsDeployExpiredEvent()
	require.NoError(t, err)
	assert.Equal(t, "7ecf22fc284526c6db16fb6455f489e0a9cbf782834131c010cf3078fb9be353", res.DeployExpired.DeployHash.ToHex())
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
	_, err = json.Marshal(res)
	assert.NoError(t, err)
}
