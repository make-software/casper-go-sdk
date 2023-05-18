package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/types"
)

func Test_DeployInfo_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	fixture, err := os.ReadFile("../data/deploy/deploy_info_example.json")
	assert.NoError(t, err)

	var deployInfo types.DeployInfo
	err = json.Unmarshal(fixture, &deployInfo)
	assert.NoError(t, err)

	result, err := json.Marshal(deployInfo)
	assert.NoError(t, err)
	assert.JSONEq(t, string(fixture), string(result))
}
