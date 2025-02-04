package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/make-software/casper-go-sdk/v2/types"
)

func Test_EraSummary_MarshalUnmarshal_ShouldReturnSameResult(t *testing.T) {
	tests := []struct {
		name        string
		fixturePath string
		isPurse     bool
	}{
		{
			"V1 EraSummary",
			"../data/era/era_summary_example.json",
			false,
		},
		{
			"V2 EraSummary",
			"../data/era/era_summary_v2_delegator_kind_purse.json",
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture, err := os.ReadFile(test.fixturePath)
			assert.NoError(t, err)

			var era types.EraSummary
			err = json.Unmarshal(fixture, &era)
			assert.NoError(t, err)

			for _, summary := range era.StoredValue.EraInfo.SeigniorageAllocations {
				if summary.Delegator != nil {
					assert.NotEmpty(t, summary.Delegator.DelegatorKind.ToHex())
					assert.NotEmpty(t, summary.Delegator.ValidatorPublicKey.ToHex())
					assert.NotEmpty(t, summary.Delegator.Amount)
				}

				if summary.Validator != nil {
					assert.NotEmpty(t, summary.Validator.ValidatorPublicKey.ToHex())
					assert.NotEmpty(t, summary.Validator.Amount)
				}
			}
		})
	}
}
