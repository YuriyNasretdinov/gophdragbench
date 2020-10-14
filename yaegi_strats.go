package main

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/golangconf/gophers-and-dragons/game"
	"github.com/traefik/yaegi/interp"
)

func loadYaegiStrats(dir string) []strat {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Panicf("read strats dir: %v", err)
	}

	var strats []strat
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filename := filepath.Join(dir, f.Name())
		code, err := ioutil.ReadFile(filename)
		if err != nil {
			strats = append(strats, strat{name: f.Name(), err: err})
			continue
		}
		strats = append(strats, newYaegiStrat(f.Name(), string(code)))
	}

	return strats
}

func newYaegiStrat(name, code string) strat {
	i := interp.New(interp.Options{})

	out := strat{name: name}

	i.Use(map[string]map[string]reflect.Value{
		"github.com/golangconf/gophers-and-dragons/game": {
			"Creeps": reflect.ValueOf(game.Creeps),
			"Cards":  reflect.ValueOf(game.Cards),

			"State":          reflect.ValueOf((*game.State)(nil)),
			"Avatar":         reflect.ValueOf((*game.Avatar)(nil)),
			"AvatarStats":    reflect.ValueOf((*game.AvatarStats)(nil)),
			"Card":           reflect.ValueOf((*game.Card)(nil)),
			"CardStats":      reflect.ValueOf((*game.CardStats)(nil)),
			"CardType":       reflect.ValueOf((*game.CardType)(nil)),
			"Creep":          reflect.ValueOf((*game.Creep)(nil)),
			"CreepStats":     reflect.ValueOf((*game.CreepStats)(nil)),
			"CreepType":      reflect.ValueOf((*game.CreepType)(nil)),
			"CreepTrait":     reflect.ValueOf((*game.CreepTrait)(nil)),
			"CreepTraitList": reflect.ValueOf((*game.CreepTraitList)(nil)),
			"IntRange":       reflect.ValueOf((*game.IntRange)(nil)),

			"CreepCheepy": reflect.ValueOf(game.CreepCheepy),
			"CreepImp":    reflect.ValueOf(game.CreepImp),
			"CreepLion":   reflect.ValueOf(game.CreepLion),
			"CreepFairy":  reflect.ValueOf(game.CreepFairy),
			"CreepClaws":  reflect.ValueOf(game.CreepClaws),
			"CreepMummy":  reflect.ValueOf(game.CreepMummy),
			"CreepDragon": reflect.ValueOf(game.CreepDragon),

			"TraitCoward":        reflect.ValueOf(game.TraitCoward),
			"TraitMagicImmunity": reflect.ValueOf(game.TraitMagicImmunity),
			"TraitWeakToFire":    reflect.ValueOf(game.TraitWeakToFire),
			"TraitSlow":          reflect.ValueOf(game.TraitSlow),
			"TraitRanged":        reflect.ValueOf(game.TraitRanged),
			"TraitBloodlust":     reflect.ValueOf(game.TraitBloodlust),

			"CardMagicArrow":  reflect.ValueOf(game.CardMagicArrow),
			"CardAttack":      reflect.ValueOf(game.CardAttack),
			"CardPowerAttack": reflect.ValueOf(game.CardPowerAttack),
			"CardStun":        reflect.ValueOf(game.CardStun),
			"CardFirebolt":    reflect.ValueOf(game.CardFirebolt),
			"CardRetreat":     reflect.ValueOf(game.CardRetreat),
			"CardRest":        reflect.ValueOf(game.CardRest),
			"CardHeal":        reflect.ValueOf(game.CardHeal),
			"CardParry":       reflect.ValueOf(game.CardParry),
		},
	})

	if _, err := i.Eval(code); err != nil {
		out.err = err
		return out
	}

	pkg := inferPackage(code)
	userFuncSym := "ChooseCard"
	if pkg != "" {
		userFuncSym = pkg + ".ChooseCard"
	}

	res, err := i.Eval(userFuncSym)
	if err != nil {
		out.err = errors.New("can't find proper ChooseCard definition")
		return out
	}
	userFunc, ok := res.Interface().(func(game.State) game.CardType)
	if !ok {
		out.err = errors.New("can't find proper ChooseCard definition")
		return out
	}

	out.cb = userFunc
	return out
}

func inferPackage(s string) string {
	newline := strings.IndexByte(s, '\n')
	if newline == -1 {
		return ""
	}
	line := s[:newline]
	if !strings.HasPrefix(line, "package ") {
		return ""
	}
	packageName := line[len("package "):]
	return packageName
}
