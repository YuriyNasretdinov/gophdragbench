package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/quasilyte/gophers-and-dragons/game"
	"github.com/quasilyte/gophers-and-dragons/wasm/sim"
	"github.com/quasilyte/gophers-and-dragons/wasm/simstep"
)

var (
	iterations  = flag.Int("iterations", 1000, "how many iterations to do for every strategy")
	filterRegex = flag.String("filter-regex", "", "regexp for strategy names to leave")
	debug       = flag.Bool("debug", false, "enable debug messages")
	samples     = flag.Bool("samples", false, "show best and worst matches")
)

func runsim(chooseCardFn func(game.State) game.CardType) (log []simstep.Action, score int, dragonDefeated bool) {
	conf := &sim.Config{
		AvatarHP: 40,
		AvatarMP: 20,
		Rounds:   10,
	}

	victory := false

	actions := sim.Run(conf, chooseCardFn)
	for _, act := range actions {
		switch act := act.(type) {
		case simstep.UpdateScore:
			score += act.Delta
		case simstep.Victory:
			victory = true
		case simstep.GreenLog:
			if strings.Contains(act.Message, "Dragon is defeated") {
				dragonDefeated = true
			}
		}
	}

	if victory {
		return actions, score, dragonDefeated
	}

	return nil, 0, false
}

func tryDisableDebugMessages() {
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err == nil {
		syscall.Dup2(int(devNull.Fd()), int(os.Stderr.Fd()))
		devNull.Close()
	}
}

// computeAvgScore computes average for N rounds of best-of-three
// (it does three launches and takes the best score of three launches
//  and then returns the average of those scores)
func computeAvgScore(chooseCardFn func(game.State) game.CardType) (avg, meanerr, winratio, dragonDefeatRatio float64, worstActions, bestActions []simstep.Action, minScore, maxScore int) {
	sum := 0.0
	wins := 0
	dragonDefeats := 0
	total := 0
	minScore = 100000

	bests := make([]int, 0, *iterations)

	for i := 0; i < *iterations; i++ {
		best := 0

		for j := 0; j < 3; j++ {
			actions, res, dragonDefeated := runsim(chooseCardFn)
			if res > best {
				best = res
			}
			if res > 0 {
				wins++
			}
			if dragonDefeated {
				dragonDefeats++
			}
			if res > maxScore {
				maxScore = res
				bestActions = actions
			}
			if res > 0 && res < minScore {
				minScore = res
				worstActions = actions
			}
			total++
		}

		bests = append(bests, best)
		sum += float64(best)
	}

	sumsquares := 0.0
	avg = sum / float64(*iterations)
	for _, best := range bests {
		sumsquares += (float64(best) - avg) * (float64(best) - avg)
	}

	// mean error is calculated as standard deviation / sqrt(N)
	// and standard deviation is calculated as
	//    sqrt( (x - avg(x))^2 / (N-1) )
	// if N is large enough we can just say that meanerr = sqrt( (x - avg(x))^2 ) / N
	meanerr = math.Sqrt(sumsquares) / float64(*iterations)

	return avg, meanerr, float64(wins) / float64(total), float64(dragonDefeats) / float64(total), worstActions, bestActions, minScore, maxScore
}

func printAction(act simstep.Action) {
	switch act := act.(type) {
	case simstep.Log:
		fmt.Printf("  %s\n", act.Message)
	case simstep.GreenLog:
		fmt.Printf("  \033[1;32m%s\033[0m\n", act.Message)
	case simstep.RedLog:
		fmt.Printf("  \033[1;31m%s\033[0m\n", act.Message)
	}
}

func main() {
	flag.Parse()

	var filter *regexp.Regexp
	if *filterRegex != "" {
		filter = regexp.MustCompile(*filterRegex)
	}

	rand.Seed(time.Now().UnixNano())

	if !*debug {
		tryDisableDebugMessages()
	}

	maxLen := 0
	for _, s := range strats {
		if len(s.name) > maxLen {
			maxLen = len(s.name)
		}
	}

	for _, s := range strats {
		if filter != nil && !filter.MatchString(s.name) {
			continue
		}

		fmt.Printf("Avg score for strat %q: ", s.name)
		if len(s.name) < maxLen {
			fmt.Printf("%s", strings.Repeat(" ", maxLen-len(s.name)))
		}
		start := time.Now()
		avg, meanerr, winratio, dragonRatio, worst, best, min, max := computeAvgScore(s.cb)
		avgTime := time.Since(start) / time.Duration(*iterations)
		avgStr := fmt.Sprintf("%.2f", avg)
		fmt.Printf("%6s ± %.2f, wins: %2d%%, dragon kills: %2d%%, time per game: %s\n", avgStr, meanerr, int(winratio*100), int(dragonRatio*100), avgTime)

		if *samples {
			fmt.Printf("Worst game (%d points):\n", min)
			for _, a := range worst {
				printAction(a)
			}

			fmt.Printf("\nBest game (%d points):\n", max)
			for _, a := range best {
				printAction(a)
			}
		}
	}
}
