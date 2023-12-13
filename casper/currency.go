package casper

import (
	"github.com/shopspring/decimal"
)

var MotesToCSPRRate = decimal.NewFromInt(1000000000)

type Motes uint64

func (m Motes) ToCSPR() decimal.Decimal {
	dec := decimal.NewFromInt(int64(m))

	return dec.Div(MotesToCSPRRate)
}
