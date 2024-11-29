package types

// BidKind Auction bid variants.
type BidKind struct {
	// A unified record indexed on validator data, with an embedded collection of all delegator bids assigned to that validator.
	Unified *Bid `json:"Unified,omitempty"`
	// A bid record containing only validator data.
	Validator *ValidatorBid `json:"Validator,omitempty"`
	// A bid record containing only delegator data.
	Delegator *Delegator `json:"Delegator,omitempty"`
	// A bridge record pointing to a new `ValidatorBid` after the public key was changed.
	Bridge *Bridge `json:"Bridge,omitempty"`
	// New validator public key associated with the bid.
	Credit *Credit `json:"Credit,omitempty"`
	//// Represents a validator reserving a slot for specific delegator
	//Reservation *Reservation `json:"Reservation,omitempty"`
	//Unbond      *Unbond       `json:"Unbond,omitempty"`
}
