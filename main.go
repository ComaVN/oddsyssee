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
	//  - games can have multiple odds and probabilities, and betting strategies can take different odds at different times
	single_win_probability := 18.0 / 37
	min_bet := decimal.RequireFromString("5")
	max_bet := decimal.RequireFromString("5")
	bet_multiplier := decimal.NewFromInt(2)

	// Player properties
	player_props := betting.NewPlayerProps(
		decimal.RequireFromString("1000"),
		decimal.RequireFromString("1100"),
	)
	betting_strategies := []betting.BettingSystem{
		betting.NewSingleBetStrategy(player_props),
		betting.NewMartingaleStrategy(player_props),
		betting.NewOscarsGrindStrategy(player_props),
	}

	// Simulation properties
	sim_repeats := 10 // Number of tines the full simulation is repeated

	fmt.Printf("Game properties:\n"+
		"  Probability to win one bet: %.2f%%\n"+
		"  Minimum bet: €%v\n"+
		"  Maximum bet: €%v\n"+
		"Player properties:\n"+
		"  Bankroll: €%v\n"+
		"  Win Target: €%v\n"+
		"Simulation properties:\n"+
		"  runs: %d\n"+
		"\n",
		single_win_probability*100,
		min_bet.StringFixedBank(2),
		max_bet.StringFixedBank(2),
		player_props.Bankroll.StringFixedBank(2),
		player_props.WinTarget.StringFixedBank(2),
		sim_repeats,
	)

	fmt.Println("Starting simulation")
	for i := 1; i <= sim_repeats; i++ {
		fmt.Printf("  Run %d\n", i)
		// Initialize players for each strategy
		players := make([]*betting.Player, 0, len(betting_strategies))
		for _, bs := range betting_strategies {
			players = append(players, bs.NewPlayer())
		}
		rounds_cnt := 0
		for len(players) > 0 {
			// TODO: there should be a way to catch strategies that never win or lose
			rounds_cnt++
			fmt.Printf("    Round %d\n", rounds_cnt)
			this_round_players := make([]*betting.Player, 0, len(players))
			bet_sizes := make([]decimal.Decimal, 0, len(players))
			for idx, pl := range players {
				bet_size, ok := pl.PlayNextBet()
				if !ok {
					fmt.Printf("    Player %d with betting system '%v' is out of money\n", idx, pl.BettingSystem.Name())
					continue
				}
				this_round_players = append(this_round_players, pl)
				bet_sizes = append(bet_sizes, bet_size)
				fmt.Printf("    Player %d with betting system '%v' current bank €%v, betted €%v\n", idx, pl.BettingSystem.Name(), pl.CurrentBank.Add(bet_size), bet_size)
			}
			bet_won := rand.Float64() < single_win_probability
			if bet_won {
				fmt.Println("    Players won")
			} else {
				fmt.Println("    Players lost")
			}
			players = make([]*betting.Player, 0, len(this_round_players))
			for idx, pl := range this_round_players {
				if bet_won {
					bet_payout := bet_sizes[idx].Mul(bet_multiplier)
					target_reached := pl.Win(bet_payout)
					if target_reached {
						fmt.Printf("    Player %d with betting system '%v' current bank €%v, has reached their win target of €%v\n", idx, pl.BettingSystem.Name(), pl.CurrentBank, pl.WinTarget)
						continue
					}
				} else {
					if pl.Lose() {
						fmt.Printf("    Player %d with betting system '%v' is out of money\n", idx, pl.BettingSystem.Name())
						continue
					}
				}
				players = append(players, pl)
				fmt.Printf("    Player %d with betting system '%v' current bank €%v\n", idx, pl.BettingSystem.Name(), pl.CurrentBank)
			}

			fmt.Println("")
		}
	}
}
