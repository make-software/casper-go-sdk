package cl_value

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types"
)

func Test_ArgsParser_MapFromRawJson(t *testing.T) {
	source := `{
            "bytes": "030000000600000069737375657242000000303161333538383766333936326136613233326538653131666137643435363762363836366436383835303937346161643732383965663238373637363832356636070000006e6574776f726b2b000000746e69433248583579673279446a4d514563556f316248613434783959645a565371794b6f78323153447a0600000073746174757306000000616374697665",
            "parsed": [
              {
                "key": "issuer",
                "value": "01a35887f3962a6a232e8e11fa7d4567b6866d68850974aad7289ef287676825f6"
              },
              {
                "key": "network",
                "value": "tniC2HX5yg2yDjMQEcUo1bHa44x9YdZVSqyKox21SDz"
              },
              {
                "key": "status",
                "value": "active"
              }
            ],
            "cl_type": {
              "Map": {
                "key": "String",
                "value": "String"
              }
            }
          }`
	data, err := types.ArgsFromRawJson(json.RawMessage(source))
	require.NoError(t, err)
	assert.Equal(t, "active", data.Map.Get("status").StringVal.String())
}
