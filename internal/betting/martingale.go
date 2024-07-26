// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package betting

import (
	"github.com/shopspring/decimal"
)

type martingaleStrategy struct {
	PlayerProps
	previousBank decimal.Decimal
}

func NewMartingaleStrategy(playerProps PlayerProps) *martingaleStrategy {
	return &martingaleStrategy{
		PlayerProps: playerProps,
	}
}

// Implements BettingSystem
func (bs *martingaleStrategy) Name() string {
	return "Martingale"
}

// Implements BettingSystem
func (bs *martingaleStrategy) NextBet(current_bank decimal.Decimal) decimal.Decimal {
	previous_bet_won := current_bank.GreaterThan(bs.previousBank)
	bs.previousBank = current_bank
	if previous_bet_won {
		return bs.WinTarget.Sub(bs.Bankroll)
	}
	bet_size := bs.WinTarget.Sub(current_bank)
	if bet_size.GreaterThan(current_bank) {
		return current_bank
	}
	return bet_size
}

func (bs *martingaleStrategy) NewPlayer() *Player {
	return &Player{
		PlayerProps:   bs.PlayerProps,
		CurrentBank:   bs.Bankroll,
		BettingSystem: bs,
	}
}
