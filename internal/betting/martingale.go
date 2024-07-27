// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package betting

import (
	"github.com/shopspring/decimal"
)

type martingaleStrategy struct {
	PlayerProps
}

func NewMartingaleStrategy(playerProps PlayerProps) *martingaleStrategy {
	return &martingaleStrategy{
		PlayerProps: playerProps,
	}
}

// Implements [BettingSystem]
func (bs *martingaleStrategy) Name() string {
	return "Martingale"
}

// Implements [BettingSystem]
// ATTN: This is not implement as "double bet each time you lose", but as "bet whatever takes you to the target"
// TODO: this should take into account minimum bets, and repeatedly reach those as "sub-target".
func (bs *martingaleStrategy) NextBet(currentBank decimal.Decimal) Bet {
	betSize := bs.WinTarget.Sub(currentBank)
	if betSize.GreaterThan(currentBank) {
		return NewBet(currentBank)
	}
	return NewBet(betSize)
}
