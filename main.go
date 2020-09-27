package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/quasilyte/gophers-and-dragons/game"
	"github.com/quasilyte/gophers-and-dragons/wasm/sim"
	"github.com/quasilyte/gophers-and-dragons/wasm/simstep"
)

var (
	iterations = flag.Int("iterations", 1000, "how many iterations to do for every strategy")
	debug      = flag.Bool("debug", false, "enable debug messages")
)

func runsim(chooseCardFn func(game.State) game.CardType) (score int) {
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
		}
	}

	if victory {
		return score
	}

	return 0
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
func computeAvgScore(chooseCardFn func(game.State) game.CardType) (avg float64, meanerr float64, winratio float64) {
	sum := 0.0
	wins := 0
	total := 0

	bests := make([]int, 0, *iterations)

	for i := 0; i < *iterations; i++ {
		best := 0

		for j := 0; j < 3; j++ {
			res := runsim(chooseCardFn)
			if res > best {
				best = res
			}
			if res > 0 {
				wins++
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

	return avg, meanerr, float64(wins) / float64(total)
}

func main() {
	flag.Parse()
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
		fmt.Printf("Avg score for strat %q: ", s.name)
		if len(s.name) < maxLen {
			fmt.Printf("%s", strings.Repeat(" ", maxLen-len(s.name)))
		}
		start := time.Now()
		avg, meanerr, winratio := computeAvgScore(s.cb)
		avgTime := time.Since(start) / time.Duration(*iterations)
		fmt.Printf("%.1f ± %.1f, winrate: %d%%, time per game: %s\n", avg, meanerr, int(winratio*100), avgTime)
	}
}
