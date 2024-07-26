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

// Implements BettingSystem
func (bs *oscarsGrindStrategy) Name() string {
	return "Oscar's Grind"
}

// Implements BettingSystem
func (bs *oscarsGrindStrategy) NextBet(current_bank decimal.Decimal) decimal.Decimal {
	previous_bet_won := current_bank.GreaterThan(bs.previousBank)
	bs.previousBank = current_bank
	if previous_bet_won {
		bs.currentBetSize = bs.currentBetSize.Add(bs.WinTarget).Sub(bs.Bankroll)
		if current_bank.Add(bs.currentBetSize).GreaterThan(bs.WinTarget) {
			bs.currentBetSize = bs.WinTarget.Sub(current_bank)
		}
	}
	if bs.currentBetSize.GreaterThan(current_bank) {
		return current_bank
	}
	return bs.currentBetSize
}

func (bs *oscarsGrindStrategy) NewPlayer() *Player {
	return &Player{
		PlayerProps:   bs.PlayerProps,
		CurrentBank:   bs.Bankroll,
		BettingSystem: bs,
	}
}
