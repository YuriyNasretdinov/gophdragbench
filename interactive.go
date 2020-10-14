package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/golangconf/gophers-and-dragons/game"
)

var stdinReader = bufio.NewReader(os.Stdin)

func interactivePlay(s game.State) game.CardType {
	fmt.Printf("\n")
	fmt.Printf("--- Round %d ---\n\n", s.Round)
	fmt.Printf("Avatar:                         HP: %d    MP: %d\n", s.Avatar.HP, s.Avatar.MP)
	creepDmg := fmt.Sprintf("(%d-%d dmg)", s.Creep.Damage.Low(), s.Creep.Damage.High())

	stunText := ""
	if s.Creep.IsStunned() {
		stunText = " (stunned)"
	}

	fmt.Printf("Creep: %10s %10s%s    HP: %d     %s\nNext:  %10s\n", s.Creep.Type, creepDmg, stunText, s.Creep.HP, s.Creep.Traits, s.NextCreep)

	fmt.Printf("Deck:\n")

	var options []int
	for _, c := range s.Deck {
		if !s.Can(c.Type) {
			continue
		}

		if c.Count > 0 || c.Count == -1 {
			options = append(options, int(c.Type))
		}
	}
	sort.Ints(options)

	optMap := make(map[string]game.CardType)

	for _, o := range options {
		c := s.Deck[game.CardType(o)]
		optMap[c.Type.String()] = c.Type

		dmg := ""
		if c.IsOffensive && c.Type != game.CardStun {
			if c.Power.Low() == c.Power.High() {
				dmg = fmt.Sprintf("%d dmg", c.Power.Low())
			} else {
				dmg = fmt.Sprintf("%d-%d dmg", c.Power.Low(), c.Power.High())
			}

		} else if c.Effect != "" {
			if c.Power.Low() == c.Power.High() {
				dmg = fmt.Sprintf("%d %s", c.Power.Low(), c.Effect)
			} else {
				dmg = fmt.Sprintf("%d-%d %s", c.Power.Low(), c.Power.High(), c.Effect)
			}
		}

		left := ""
		if c.Count == -1 {
			left = "∞"
		} else {
			left = fmt.Sprintf("%d", c.Count)
		}

		available := true
		if c.Type == game.CardStun && s.Creep.IsStunned() {
			available = false
		}

		if c.Type == game.CardParry && s.Creep.Traits.Has(game.TraitRanged) {
			available = false
		}

		if s.Creep.Traits.Has(game.TraitMagicImmunity) && (c.Type == game.CardMagicArrow || c.Type == game.CardFirebolt) {
			available = false
		}

		if !available {
			dmg = "nothing"
		}

		fmt.Printf("   [%d] %14s: %s   %d mp   %s\n", o, c.Type, left, c.MP, dmg)
	}
	fmt.Printf("\n")

	var card game.CardType

	for {
		fmt.Printf("\nChoose card [Attack]: ")
		opt, err := stdinReader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		var ok bool

		opt = strings.TrimSpace(opt)
		if opt == "" {
			card = 0
			break
		}

		card, ok = optMap[opt]
		if !ok {
			optNum, err := strconv.Atoi(opt)
			if err != nil {
				fmt.Printf("Unrecognized card: %s\n", opt)
				continue
			}

			card = game.CardType(optNum)
		}

		if !s.Can(card) {
			fmt.Printf("Cannot use this card\n")
			continue
		}

		break
	}

	fmt.Printf("Using %s\n", card)

	return card
}
