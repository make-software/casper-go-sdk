This package unites all types of identification used in the Casper network. The majority of them are represented as 32-byte hashes, although some of them contain prefixes (such as account, contract, contract package, and transfer).

The Key represents a complex and universal key type used in the Casper network. [(See documentation for more information.)](https://docs.casper.network/concepts/serialization-standard/#serialization-standard-state-keys)

The Unforgeable Reference (URef) stores access rights apart from the hash and prefix. [(See documentation for more information.)]((https://docs.casper.network/concepts/design/casper-design/#uref-head))

The exception to the 32-byte hash rule is the identifier for the Era, which is represented as a uint value. [(See documentation for more information.)](https://docs.casper.network/concepts/serialization-standard/#serialization-standard-era-info-key)

