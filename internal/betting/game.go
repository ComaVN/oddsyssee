package betting

import "math/rand/v2"

// An iterator, representing an infinite stream of boolean game outcomes.
type NextOutcomer interface {
	NextOutcome() bool
}

type fixedProbabilityGame struct {
	probability float64
}

func NewFixedProbabilityGame(probability float64) *fixedProbabilityGame {
	return &fixedProbabilityGame{
		probability: probability,
	}
}

// Implements NextOutcomer
func (g *fixedProbabilityGame) NextOutcome() bool {
	return rand.Float64() < g.probability
}

type LosingGame struct{}

// Implements NextOutcomer
func (g LosingGame) NextOutcome() bool {
	return false
}

type WinningGame struct{}

// Implements NextOutcomer
func (g WinningGame) NextOutcome() bool {
	return true
}

type alternatingGame struct {
	nextOutcome bool
}

// Implements NextOutcomer
func (g *alternatingGame) NextOutcome() bool {
	outcome := g.nextOutcome
	g.nextOutcome = !outcome
	return outcome
}

func NewAlternatingGame(firstOutcome bool) *alternatingGame {
	return &alternatingGame{
		nextOutcome: firstOutcome,
	}
}
