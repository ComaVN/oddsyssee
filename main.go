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
	singleWinProbability := 18.0 / 37
	var game betting.NextOutcomer
	game = betting.NewFixedProbabilityGame(singleWinProbability)
	// game = betting.NewAlternatingGame(false)
	// game = betting.WinningGame{}
	// game = betting.LosingGame{}
	minBet := decimal.RequireFromString("5")
	maxBet := decimal.RequireFromString("5")
	betMultiplier := decimal.NewFromInt(2)

	// Player properties
	playerProps := betting.NewPlayerProps(
		decimal.RequireFromString("1000"),
		decimal.RequireFromString("1100"),
	)
	bettingStrategies := []betting.BettingSystem{
		betting.NewSingleBetStrategy(playerProps),
		betting.NewMartingaleStrategy(playerProps),
		betting.NewOscarsGrindStrategy(playerProps),
	}

	// Simulation properties
	simRepeats := 10 // Number of tines the full simulation is repeated

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
		singleWinProbability*100,
		minBet.StringFixedBank(2),
		maxBet.StringFixedBank(2),
		playerProps.Bankroll.StringFixedBank(2),
		playerProps.WinTarget.StringFixedBank(2),
		simRepeats,
	)

	fmt.Println("Starting simulation")
	for i := 1; i <= simRepeats; i++ {
		fmt.Printf("  Run %d\n", i)
		// Initialize players for each strategy
		players := make([]betting.Player, 0, len(bettingStrategies))
		for idx, bs := range bettingStrategies {
			players = append(players, betting.NewPlayer(playerProps, fmt.Sprintf("%d", idx), bs))
		}
		roundsCnt := 0
		for {
			// TODO: there should be a way to catch strategies that never win or lose
			roundsCnt++
			fmt.Printf("    Round %d\n", roundsCnt)
			thisRoundPlayers := make([]betting.Player, 0, len(players))
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
				thisRoundPlayers = append(thisRoundPlayers, pl)
				bets = append(bets, bet)
				fmt.Printf("    Player %q with strategy %q current bank €%v, betted €%v\n", pl.Name(), pl.StrategyName(), pl.Bank(), bet.Size())
			}
			players = thisRoundPlayers
			if len(players) == 0 {
				break
			}
			betWon := game.NextOutcome()
			if betWon {
				fmt.Println("    Players won")
			} else {
				fmt.Println("    Players lost")
			}
			for _, bet := range bets {
				pl := bet.Player()
				var outcome betting.Outcome
				if betWon {
					outcome = betting.NewOutcome(bet, bet.Size().Mul(betMultiplier))
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
