package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-software/casper-go-sdk/v2/types"
)

func Test_ExecutableItem_MarshalUnmarshal_ShouldBeSameResult(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			"item with module bytes",
			`{"ModuleBytes": {"module_bytes": "c4c411864f7b717c27839e56f6f1ebe5da3f35ec0043f437324325d65a22afa4","args": [["amount",{"bytes":"060040f538a406","parsed":"7302400000000","cl_type":"U512"}],["target",{"bytes":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","parsed":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","cl_type":"PublicKey"}],["id",{"bytes":"016dec181a95010000","parsed":1739899595885,"cl_type":{"Option":"U64"}}]]}}`,
		},
		{
			"item with stored contract by hash",
			`{"StoredContractByHash": {"hash": "c4c411864f7b717c27839e56f6f1ebe5da3f35ec0043f437324325d65a22afa4","entry_point": "pclphXwfYmCmdITj8hnh","args": [["amount",{"bytes":"060040f538a406","parsed":"7302400000000","cl_type":"U512"}],["target",{"bytes":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","parsed":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","cl_type":"PublicKey"}],["id",{"bytes":"016dec181a95010000","parsed":1739899595885,"cl_type":{"Option":"U64"}}]]}}`,
		},
		{
			"item with stored contract by name",
			`{"StoredContractByName": {"name": "U5A74bSZH8abT8HqVaK9","entry_point": "gIetSxltnRDvMhWdxTqQ","args": [["amount",{"bytes":"060040f538a406","parsed":"7302400000000","cl_type":"U512"}],["target",{"bytes":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","parsed":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","cl_type":"PublicKey"}],["id",{"bytes":"016dec181a95010000","parsed":1739899595885,"cl_type":{"Option":"U64"}}]]}}`,
		},
		{
			"item with stored versioned contract by hash without version",
			`{"StoredVersionedContractByHash": {"hash": "c4c411864f7b717c27839e56f6f1ebe5da3f35ec0043f437324325d65a22afa4","entry_point": "pclphXwfYmCmdITj8hnh","version": null,"args": [["amount",{"bytes":"060040f538a406","parsed":"7302400000000","cl_type":"U512"}],["target",{"bytes":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","parsed":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","cl_type":"PublicKey"}],["id",{"bytes":"016dec181a95010000","parsed":1739899595885,"cl_type":{"Option":"U64"}}]]}}`,
		},
		{
			"item with stored versioned contract by hash with version",
			`{"StoredVersionedContractByHash": {"hash": "c4c411864f7b717c27839e56f6f1ebe5da3f35ec0043f437324325d65a22afa4","entry_point": "pclphXwfYmCmdITj8hnh","version": 1,"args": [["amount",{"bytes":"060040f538a406","parsed":"7302400000000","cl_type":"U512"}],["target",{"bytes":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","parsed":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","cl_type":"PublicKey"}],["id",{"bytes":"016dec181a95010000","parsed":1739899595885,"cl_type":{"Option":"U64"}}]]}}`,
		},
		{
			"item with stored versioned contract by name with version",
			`{"StoredVersionedContractByName": {"name": "lWJWKdZUEudSakJzw1tn","version": 1632552656, "entry_point": "S1cXRT3E1jyFlWBAIVQ8","args": [["amount",{"bytes":"060040f538a406","parsed":"7302400000000","cl_type":"U512"}],["target",{"bytes":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","parsed":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","cl_type":"PublicKey"}],["id",{"bytes":"016dec181a95010000","parsed":1739899595885,"cl_type":{"Option":"U64"}}]]}}`,
		},
		{
			"item with stored versioned contract by name without version",
			`{"StoredVersionedContractByName": {"name": "lWJWKdZUEudSakJzw1tn","version": null, "entry_point": "S1cXRT3E1jyFlWBAIVQ8","args": [["amount",{"bytes":"060040f538a406","parsed":"7302400000000","cl_type":"U512"}],["target",{"bytes":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","parsed":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","cl_type":"PublicKey"}],["id",{"bytes":"016dec181a95010000","parsed":1739899595885,"cl_type":{"Option":"U64"}}]]}}`,
		},
		{
			"item with stored transfer",
			`{"Transfer": {"args": [["amount",{"bytes":"060040f538a406","parsed":"7302400000000","cl_type":"U512"}],["target",{"bytes":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","parsed":"011c74ebfcc1b19bc3e578bec3ecfa2d484f2a00d7e9e8152c4c70f519f6a89f6a","cl_type":"PublicKey"}],["id",{"bytes":"016dec181a95010000","parsed":1739899595885,"cl_type":{"Option":"U64"}}]]}}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var deployItem types.ExecutableDeployItem
			err := json.Unmarshal([]byte(test.data), &deployItem)
			require.NoError(t, err)

			_, err = deployItem.Bytes()
			require.NoError(t, err)

			result, err := json.Marshal(deployItem)
			require.NoError(t, err)
			assert.JSONEq(t, test.data, string(result))
		})
	}
}
