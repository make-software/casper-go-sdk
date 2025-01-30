package types

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types"
)

func Test_Bid_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		fixturePath string
	}{
		{
			"Auction Bid V1",
			"../data/bid/auction_bid_example_v1.json",
		},
		{
			"Auction Bid V2",
			"../data/bid/auction_bid_example_v2.json",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := os.ReadFile(test.fixturePath)
			require.NoError(t, err)

			var bid types.Bid
			err = json.Unmarshal(data, &bid)
			require.NoError(t, err)

			require.Equal(t, bid.StakedAmount.String(), "900000000000")
			require.Equal(t, 1, len(bid.Delegators))
			require.Equal(t, bid.Delegators[0].DelegatorKind.ToHex(), "01d829cbfb66b2b11ef8d8feb6d3f2155789fc22f407bb57f89b05f6ba4b9ae070")
			require.Equal(t, bid.Delegators[0].ValidatorPublicKey.ToHex(), "01197f6b23e16c8532c6abc838facd5ea789be0c76b2920334039bfa8b3d368d61")
		})
	}
}
