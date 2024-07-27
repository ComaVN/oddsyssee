// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package betting

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Player interface {
	Name() string
	StrategyName() string
	Bank() decimal.Decimal
	PlayingCondition() PlayingCondition
	PlaceNextBet() (PlayerBet, bool)
	HandleOutcome(outcome Outcome)
}

// Contains the properties of the player at the start of the gambling session
type player struct {
	PlayerProps
	name          string
	currentBank   decimal.Decimal
	bettingSystem BettingSystem
}

// Implements [Player]
func (pl *player) Name() string {
	return pl.name
}

// Implements [Player]
func (pl *player) StrategyName() string {
	return pl.bettingSystem.Name()
}

// Implements [Player]
func (pl *player) Bank() decimal.Decimal {
	return pl.currentBank
}

// Implements [Player] using a betting system
// This assumes all betting systems only pass when out of money
func (pl *player) PlaceNextBet() (PlayerBet, bool) {
	if !pl.currentBank.IsPositive() {
		return nil, false
	}
	bet := pl.bettingSystem.NextBet(pl.currentBank)
	if bet.Size().GreaterThan(pl.currentBank) {
		panic(fmt.Sprintf("Betting system '%v' proposes illegal bet of €%v", pl.bettingSystem.Name(), bet))
	}
	pl.currentBank = pl.currentBank.Sub(bet.Size())
	return NewPlayerBet(pl, bet), true
}

// Implements [Player]
func (pl *player) PlayingCondition() PlayingCondition {
	if pl.currentBank.GreaterThanOrEqual(pl.WinTarget) {
		return Won
	}
	if !pl.currentBank.IsPositive() {
		return Lost
	}
	return Playing
}

// Implements [Player]
func (pl *player) HandleOutcome(outcome Outcome) {
	pl.currentBank = pl.currentBank.Add(outcome.Payout())
}

func NewPlayer(playerProps PlayerProps, name string, bettingSystem BettingSystem) Player {
	return &player{
		PlayerProps:   playerProps,
		name:          name,
		currentBank:   playerProps.Bankroll,
		bettingSystem: bettingSystem,
	}
}
