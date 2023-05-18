
## Architecture

CLType and CLValue represent values with corresponding data types in Casper network that is used by smart contracts and deploys.
This data communicate thought Network in the bytes' representation. Bytes serialization rules are described [here](https://docs.casperlabs.io/concepts/serialization-standard/#clvalue). 

The base functionality that provides by SDK is the working with bytes, string and json representation.

By the design all CLValues are represented in separate files, where is located the corresponded logic to work with a certain type. Only complex CLTypes are implemented as a separate types, other (simple types) is a values of a `SimpleType` data structure.
The common logic to parse unknown types and values is stored in the `parser` files in corresponded packages.

Also, data from network is often represented as mix of json and byte representation. 
As example in `args` of the `transfer` of a `deploy's` `succession result` could be fiend next string: `{"bytes":"01000000030000004142430a000000","cl_type":{"Map":{"key":"String","value":"I32"}}}`. In this representation of data the bytes contains only encoded value and `cl_type` field represented as json object. 

To convenient work with `args` implemented 'lazy load' behavior with delayed parsing by the direct client call, the `ArgsParser` struct responsible to do it.  

Please check examples of usage in [example_test.go](..%2F..%2Ftests%2Ftypes%2Fcl_value%2Fexample_test.go)