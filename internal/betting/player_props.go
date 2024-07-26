// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package betting

import (
	"github.com/shopspring/decimal"
)

// Contains the properties of the player at the start of the gambling session
type PlayerProps struct {
	Bankroll  decimal.Decimal
	TargetWin decimal.Decimal
}

func NewPlayerProps(bankroll decimal.Decimal, target_win decimal.Decimal) PlayerProps {
	return PlayerProps{
		Bankroll:  bankroll,
		TargetWin: target_win,
	}
}
