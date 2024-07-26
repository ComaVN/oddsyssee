// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"

	"github.com/ComaVN/oddsyssee/internal/betting"
	"github.com/shopspring/decimal"
)

func main() {
	// Game properties
	// TODO: this should be way more parameterized:
	//  - games can have other odds than 1/1
	//  - games can have different odds and probabilities, and betting strategies can take different odds at different times
	single_win_probability := 18.0 / 37
	min_bet := decimal.RequireFromString("5")
	max_bet := decimal.RequireFromString("5")
	bet_multiplier := decimal.NewFromInt(2)

	// Player properties
	player_props := betting.NewPlayerProps(
		decimal.RequireFromString("1000"),
		decimal.RequireFromString("100"),
	)
	betting_strategies := []betting.BettingSystem{
		betting.NewSingleBetStrategy(player_props),
	}

	// Simulation properties
	sim_repeats := 10 // Number of tines the full simulation is repeated

	fmt.Printf("Game properties:\n"+
		"  Probability to win one round: %.2f%%\n"+
		"  Minimum bet: €%v\n"+
		"  Maximum bet: €%v\n"+
		"Player properties:\n"+
		"  Bankroll: €%v\n"+
		"  Target winnings: €%v\n"+
		"Simulation properties:\n"+
		"  runs: %d\n"+
		"\n",
		single_win_probability*100,
		min_bet.StringFixedBank(2),
		max_bet.StringFixedBank(2),
		player_props.Bankroll.StringFixedBank(2),
		player_props.TargetWin.StringFixedBank(2),
		sim_repeats,
	)

	fmt.Println("Starting simulation")
	for i := 0; i < sim_repeats; i++ {
		fmt.Printf("  Run %d\n", i)
		// Initialize players for each strategy
		var players []*betting.Player
		for _, bs := range betting_strategies {
			players = append(players, bs.NewPlayer())
		}
		for len(players) > 0 {
			// TODO: there should be a way to catch strategies that never win or lose
			for idx := len(players) - 1; idx >= 0; idx-- {
				pl := players[idx]
				fmt.Printf("    Player %d with betting system '%v' current bank €%v\n", idx, pl.BettingSystem.Name(), pl.CurrentBank)
				bet_size, ok := pl.PlayNextBet()
				if !ok {
					fmt.Printf("    Player %d with betting system '%v' is out of money\n", idx, pl.BettingSystem.Name())
					players = append(players[:idx], players[idx+1:]...)
					continue
				}
				fmt.Printf("    Player %d with betting system '%v' current bank €%v, betted €%v\n", idx, pl.BettingSystem.Name(), pl.CurrentBank, bet_size)
				if rand.Float64() < single_win_probability {
					// Player won
					bet_payout := bet_size.Mul(bet_multiplier)
					target_reached := pl.Win(bet_payout)
					fmt.Printf("    Player %d with betting system '%v' current bank €%v, won €%v\n", idx, pl.BettingSystem.Name(), pl.CurrentBank, bet_payout)
					if target_reached {
						fmt.Printf("    Player %d with betting system '%v' current bank €%v, has reached their target win of €%v\n", idx, pl.BettingSystem.Name(), pl.CurrentBank, pl.TargetWin)
						players = append(players[:idx], players[idx+1:]...)
					}
				} else {
					// Player lost
					fmt.Printf("    Player %d with betting system '%v' current bank €%v, lost\n", idx, pl.BettingSystem.Name(), pl.CurrentBank)
					if pl.Lose() {
						fmt.Printf("    Player %d with betting system '%v' is out of money\n", idx, pl.BettingSystem.Name())
						players = append(players[:idx], players[idx+1:]...)
						continue
					}
				}
			}
		}
	}
}
