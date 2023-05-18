# RPC
Package rpc provides access to the exported methods of RPC Client and data structures where serialized response.

The `Client` implements RPC methods according to [spec](https://docs.casperlabs.io/developers/json-rpc/json-rpc-informational/).
`Client` unites implementation of `ClientInformational` interface related to [spec](https://docs.casperlabs.io/developers/json-rpc/json-rpc-informational/) and `ClientPOS` interface related to [spec](https://docs.casperlabs.io/developers/json-rpc/json-rpc-pos/).


## Usage

The configuration is flexible, and caller can provide custom configurations, the most common usage is to create `NewClient` which depends on `Handler` interface.
```
    handler := rpc.NewHttpHandler("<<NODE_RPC_API_URL>>", http.DefaultClient)
    client := rpc.NewClient(handler)
    deployResult, err := client.GetDeploy(context.Background(), deployHash)
```

For more examples check the [examples](../tests/rpc/client_example_test.go)

## Architecture

#### `Client` interface unites `ClientInformational` and `ClientPOS` interfaces. 

The motivation for separating interfaces correlates with requirements from [specification](https://docs.casperlabs.io/developers/json-rpc/).
Following this document, the different clients have different responsibilities, which could mean different patterns of usage.

#### `client` struct that implements interface `Client` depends on `Handler`

The motivation for separating this logic is splitting responsibilities. 

In this case, `client` knows nothing about the underlying network protocol used for communication with a remote server. It is responsible to builds RpcRequest corresponding to function name and params and serializes RpcResponse to the corresponding data structures.

From the other side, the `httpHandler` that implements interface `Handler` takes on responsibilities to operate with an external RPC server through HTTP. It builds and processes the request, reads a response and handles errors. All logic with HTTP interaction is isolated here and can be replaced with other (more efficient) protocols. 

Same time it knows nothing about ended data structure in which should be serialized `RpcResponse.Result` and don't care how to parameters for the RpcRequest were built. The interface supposes to use the Handler wider (not bounded by `Client` implementation). 

In the summary, the separating of interfaces supposes to extend current functionality without code modification. For examples of how to do it check the [examples](../tests/rpc/client_example_test.go)