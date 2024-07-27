// Copyright 2024 Roel Harbers.
// Use of this source code is governed by the BEER-WARE license
// that can be found in the LICENSE file.

package betting

type PlayingCondition int

const (
	Lost PlayingCondition = iota - 1
	Playing
	Won
)
