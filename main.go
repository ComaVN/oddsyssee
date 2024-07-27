// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/ComaVN/oddsyssee/internal/betting"
	"github.com/shopspring/decimal"
)

func main() {
	// Game properties
	// TODO: this should be way more parameterized:
	//  - games can have other odds than 1/1
	//  - games can have multiple odds and probabilities, and betting strategies can take different odds at different times
	single_win_probability := 18.0 / 37
	var game betting.NextOutcomer
	game = betting.NewFixedProbabilityGame(single_win_probability)
	// game = betting.NewAlternatingGame(false)
	// game = betting.WinningGame{}
	// game = betting.LosingGame{}
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
		players := make([]betting.Player, 0, len(betting_strategies))
		for idx, bs := range betting_strategies {
			players = append(players, betting.NewPlayer(player_props, fmt.Sprintf("%d", idx), bs))
		}
		rounds_cnt := 0
		for {
			// TODO: there should be a way to catch strategies that never win or lose
			rounds_cnt++
			fmt.Printf("    Round %d\n", rounds_cnt)
			this_round_players := make([]betting.Player, 0, len(players))
			bets := make([]betting.PlayerBet, 0, len(players))
			for _, pl := range players {
				switch pl.PlayingCondition() {
				case betting.Won:
					fmt.Printf("    Player %q with strategy %q current bank €%v, has reached their target\n", pl.Name(), pl.StrategyName(), pl.Bank())
					continue
				case betting.Lost:
					fmt.Printf("    Player %q with strategy %q is out of money\n", pl.Name(), pl.StrategyName())
					continue
				default:
					fmt.Printf("    Player %q with strategy %q current bank €%v\n", pl.Name(), pl.StrategyName(), pl.Bank())
				}
				bet, ok := pl.PlaceNextBet()
				if !ok {
					fmt.Printf("    Player %q with strategy %q passed\n", pl.Name(), pl.StrategyName())
					continue
				}
				this_round_players = append(this_round_players, pl)
				bets = append(bets, bet)
				fmt.Printf("    Player %q with strategy %q current bank €%v, betted €%v\n", pl.Name(), pl.StrategyName(), pl.Bank(), bet.Size())
			}
			players = this_round_players
			if len(players) == 0 {
				break
			}
			bet_won := game.NextOutcome()
			if bet_won {
				fmt.Println("    Players won")
			} else {
				fmt.Println("    Players lost")
			}
			for _, bet := range bets {
				pl := bet.Player()
				var outcome betting.Outcome
				if bet_won {
					outcome = betting.NewOutcome(bet, bet.Size().Mul(bet_multiplier))
				} else {
					outcome = betting.NewOutcome(bet, decimal.Zero)
				}
				pl.HandleOutcome(outcome)
				if outcome.Payout().IsPositive() {
					fmt.Printf("    Player %q with strategy %q current bank €%v, won €%v\n", pl.Name(), pl.StrategyName(), pl.Bank(), outcome.Payout().Sub(bet.Size()))
				} else {
					fmt.Printf("    Player %q with strategy %q current bank €%v, lost €%v\n", pl.Name(), pl.StrategyName(), pl.Bank(), bet.Size())
				}
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}
