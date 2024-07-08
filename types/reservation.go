package types

import "github.com/make-software/casper-go-sdk/types/key"

// ReservationKind Container for bytes recording location, type and data for a gas reservation
type ReservationKind struct {
	Receipt         key.Hash `json:"receipt"`
	ReservationData HexBytes `json:"reservation_data"`
	ReservationKind uint8    `json:"reservation_kind"`
}
