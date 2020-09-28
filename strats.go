package main

import (
	"github.com/YuriyNasretdinov/gophdragbench/strats/alik"
	"github.com/YuriyNasretdinov/gophdragbench/strats/atercattus"
	"github.com/YuriyNasretdinov/gophdragbench/strats/yourock"
	"github.com/quasilyte/gophers-and-dragons/game"
)

type strat struct {
	name string
	cb   func(game.State) game.CardType
}

var strats = []strat{
	{"yourock/hero", yourock.Hero},
	{"yourock/smart", yourock.Smart},
	{"yourock/live-coward", yourock.Live},
	{"yourock/live2", yourock.Live2},
	{"alik/WiningTactic", alik.WiningTactic},
	{"atercattus/First", atercattus.First},
}
