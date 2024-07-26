// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package betting

import (
	"github.com/shopspring/decimal"
)

type BettingSystem interface {
	Name() string
	NextBet(current_bank decimal.Decimal) decimal.Decimal
	NewPlayer() *Player
}
