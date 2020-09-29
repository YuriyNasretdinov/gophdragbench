package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/quasilyte/gophers-and-dragons/game"
)

var stdinReader = bufio.NewReader(os.Stdin)

func interactivePlay(s game.State) game.CardType {
	fmt.Printf("\n")
	fmt.Printf("Avatar:              HP: %d    MP: %d          ", s.Avatar.HP, s.Avatar.MP)
	for _, c := range s.Deck {
		if c.Count > 0 {
			fmt.Printf("%s: %d   ", c.Type, c.Count)
		}
	}
	fmt.Printf("\n")

	creepDmg := fmt.Sprintf("%d-%d dmg", s.Creep.Damage.Low(), s.Creep.Damage.High())

	fmt.Printf("Creep: %10s (%s)    HP: %d     Next: %s\n", s.Creep.Type, creepDmg, s.Creep.HP, s.NextCreep)

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

	for _, o := range options {
		c := s.Deck[game.CardType(o)]
		dmg := ""
		if c.IsOffensive {
			if c.Power.Low() == c.Power.High() {
				dmg = fmt.Sprintf(", %d dmg", c.Power.Low())
			} else {
				dmg = fmt.Sprintf(", %d-%d dmg", c.Power.Low(), c.Power.High())
			}

		}

		fmt.Printf("[%d] %s (%d mp%s)   ", o, c.Type, c.MP, dmg)
	}

	var cardChosen int

	for {
		fmt.Printf("\nChoose card: ")
		opt, err := stdinReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		cardChosen, err = strconv.Atoi(strings.TrimSpace(opt))
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			continue
		}

		if !s.Can(game.CardType(cardChosen)) {
			fmt.Printf("Cannot use this card\n")
			continue
		}

		break
	}

	return game.CardType(cardChosen)
}
