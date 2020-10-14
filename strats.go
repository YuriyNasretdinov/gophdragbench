package main

import (
	"github.com/golangconf/gophers-and-dragons/game"
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

	// if non-nil, represents strat load or execution error.
	err error
}

var strats = []strat{}

// var strats = []strat{
// 	{name: "yourock/hero", cb: yourock.Hero},
// 	{name: "yourock/smart", cb: nil, maker: yourock.NewSmart},
// 	{name: "yourock/live-coward", cb: yourock.Live},
// 	{name: "yourock/live2", cb: yourock.Live2},
// 	{name: "alik/WiningTactic", cb: alik.ChooseCard},
// 	{name: "atercattus/First", cb: atercattus.First},
// }
