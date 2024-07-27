// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package betting

import (
	"github.com/shopspring/decimal"
)

type oscarsGrindStrategy struct {
	PlayerProps
	currentBetSize decimal.Decimal
	previousBank   decimal.Decimal
}

func NewOscarsGrindStrategy(playerProps PlayerProps) *oscarsGrindStrategy {
	return &oscarsGrindStrategy{
		PlayerProps:    playerProps,
		currentBetSize: playerProps.WinTarget.Sub(playerProps.Bankroll),
	}
}

// Implements [BettingSystem]
func (bs *oscarsGrindStrategy) Name() string {
	return "Oscar's Grind"
}

// Implements [BettingSystem]
func (bs *oscarsGrindStrategy) NextBet(currentBank decimal.Decimal) Bet {
	previousBetWon := currentBank.GreaterThan(bs.previousBank)
	bs.previousBank = currentBank
	if previousBetWon {
		bs.currentBetSize = bs.currentBetSize.Add(bs.WinTarget).Sub(bs.Bankroll)
		if currentBank.Add(bs.currentBetSize).GreaterThan(bs.WinTarget) {
			bs.currentBetSize = bs.WinTarget.Sub(currentBank)
		}
	}
	if bs.currentBetSize.GreaterThan(currentBank) {
		return NewBet(currentBank)
	}
	return NewBet(bs.currentBetSize)
}
