// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package betting

import (
	"github.com/shopspring/decimal"
)

type Bet interface {
	Size() decimal.Decimal
}

type bet decimal.Decimal

// Implements [Bet]
func (b *bet) Size() decimal.Decimal {
	if b == nil {
		return decimal.Zero
	}
	return decimal.Decimal(*b)
}

func NewBet(size decimal.Decimal) *bet {
	return (*bet)(&size)
}

type PlayerBet interface {
	Bet
	Player() Player
}

type playerBet struct {
	Bet
	player Player
}

// Implements [PlayerBet]
func (b *playerBet) Player() Player {
	return b.player
}

func NewPlayerBet(player Player, bet Bet) *playerBet {
	return &playerBet{
		Bet:    bet,
		player: player,
	}
}

type Outcome interface {
	Bet() Bet
	Payout() decimal.Decimal
}

type outcome struct {
	bet    Bet
	payout decimal.Decimal
}

// Implements [Outcome]
func (o *outcome) Bet() Bet {
	return o.bet
}

// Implements [Outcome]
func (o *outcome) Payout() decimal.Decimal {
	return o.payout
}

func NewOutcome(bet Bet, payout decimal.Decimal) *outcome {
	return &outcome{
		bet:    bet,
		payout: payout,
	}
}
