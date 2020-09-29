package main

import (
	"github.com/YuriyNasretdinov/gophdragbench/strats/alik"
	"github.com/YuriyNasretdinov/gophdragbench/strats/atercattus"
	"github.com/YuriyNasretdinov/gophdragbench/strats/yourock"
	"github.com/quasilyte/gophers-and-dragons/game"
)

type strat struct {
	name string

	// cb can be called concurrently if `-cores=N` where N>1.
	cb func(game.State) game.CardType

	// maker, if present, is called to create a logic instance
	// for each goroutine,
	// e.g. for `-cores=4` it will be called 4 times in total.
	// The returned function won't be called concurrently.
	maker func() func(game.State) game.CardType
}

var strats = []strat{
	{"yourock/hero", yourock.Hero, nil},
	{"yourock/smart", nil, yourock.NewSmart},
	{"yourock/live-coward", yourock.Live, nil},
	{"yourock/live2", yourock.Live2, nil},
	{"alik/WiningTactic", alik.WiningTactic, nil},
	{"atercattus/First", atercattus.First, nil},
}
