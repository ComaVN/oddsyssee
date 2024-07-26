// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package betting

import (
	"github.com/shopspring/decimal"
)

// Contains the properties of the player at the start of the gambling session
type Player struct {
	PlayerProps
	CurrentBank   decimal.Decimal
	BettingSystem BettingSystem
}

func (pl *Player) PlayNextBet() (decimal.Decimal, bool) {
	if !pl.CurrentBank.IsPositive() {
		return decimal.Decimal{}, false
	}
	bet := pl.BettingSystem.NextBet(pl.CurrentBank)
	if bet.GreaterThan(pl.CurrentBank) {
		panic("Betting system proposes illegal bet")
	}
	pl.CurrentBank = pl.CurrentBank.Sub(bet)
	return bet, true
}

func (pl *Player) Win(won_amount decimal.Decimal) bool {
	pl.CurrentBank = pl.CurrentBank.Add(won_amount)
	return pl.CurrentBank.GreaterThanOrEqual(pl.Bankroll.Add(pl.TargetWin))
}

func (pl *Player) Lose() bool {
	return !pl.CurrentBank.IsPositive()
}
