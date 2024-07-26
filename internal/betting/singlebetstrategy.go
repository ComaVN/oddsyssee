// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package betting

import (
	"github.com/shopspring/decimal"
)

type singleBetStrategy struct {
	PlayerProps
}

func NewSingleBetStrategy(playerProps PlayerProps) *singleBetStrategy {
	return &singleBetStrategy{
		PlayerProps: playerProps,
	}
}

// Implements BettingSystem
func (bs *singleBetStrategy) Name() string {
	return "Single Bet"
}

// Implements BettingSystem
func (bs *singleBetStrategy) NextBet(current_bank decimal.Decimal) decimal.Decimal {
	return bs.TargetWin
}

func (bs *singleBetStrategy) NewPlayer() *Player {
	return &Player{
		PlayerProps:   bs.PlayerProps,
		CurrentBank:   bs.Bankroll,
		BettingSystem: bs,
	}
}
