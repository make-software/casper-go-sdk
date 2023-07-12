package rpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
)

func Test_QueryBalance_byPublicKey(t *testing.T) {
	client := rpc.NewClient(rpc.NewHttpHandler("http://127.0.0.1:11101/rpc", http.DefaultClient))
	pubKey, err := casper.NewPublicKey("0115394d1f395a87dfed4ab62bbfbc91b573bbb2bffb2c8ebb9c240c51d95bcc4d")
	require.NoError(t, err)
	res, err := client.QueryBalance(context.Background(), rpc.PurseIdentifier{
		MainPurseUnderPublicKey: &pubKey,
	})
	require.NoError(t, err)
	assert.Equal(t, "1000000000000000000000000000000000", res.Balance.String())
}
