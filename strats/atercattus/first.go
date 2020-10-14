package atercattus

import (
	"github.com/golangconf/gophers-and-dragons/game"
)

const restHP = 3

var (
	creepyDps = map[game.CreepType]map[bool]int{
		game.CreepCheepy: {false: 1, true: 4}, // hp 4, dp 1-4 (Coward)
		game.CreepImp:    {false: 3, true: 4}, // hp 5, dp 3-4
		game.CreepLion:   {false: 2, true: 3}, // hp 10, dp 2-3
		game.CreepFairy:  {false: 4, true: 5}, // hp 9, dp 4-5 (Ranged)
		game.CreepMummy:  {false: 3, true: 4}, // hp 18, dp 3-4 (WeakToFire, Slow)
		game.CreepDragon: {false: 5, true: 6}, // hp 30, dp 5-6 (MagicImmunity)
	}
)

func First(s game.State) game.CardType {
	card := chooseCard(s)
	println(card, ` !!!`)
	return card
}

func chooseCard(s game.State) game.CardType {
	A := s.Avatar
	C := s.Creep

	println("HP:", A.HP, " MP:", A.MP, " C:", C.Type, " CHP:", C.HP, " CD:", getCreepyDamage(C.Type, false))

	// Лечение без аптечек, если можно
	if (A.HP < A.MaxHP) && (A.MP > A.MaxMP/2) {
		if C.Traits.Has(game.TraitCoward) && (C.HP == C.MaxHP) {
			return game.CardRest
		}
		if (C.Damage.High() < restHP) && (getCreepyDamage(C.Type, false) < restHP) {
			return game.CardRest
		}
	}

	// Добивание?
	if damage, card := getDamage(s, true); C.HP <= damage {
		return card
	}

	// Лечение
	if s.Can(game.CardHeal) && ((A.HP <= A.MaxHP-15) || (C.Type == game.CreepDragon)) {
		return game.CardHeal
	}

	// Обычная атака
	myDamage, card := getDamage(s, false)
	if myDamage >= C.HP {
		return card
	}

	needToRetreat := false
	cDamage := getCreepyDamage(C.Type, true)
	if (myDamage > 0) && ((C.HP+myDamage-1)/myDamage > (A.HP+cDamage-1)/cDamage) {
		needToRetreat = true
	} else if (2*myDamage < C.HP) && (2*cDamage >= A.HP) {
		needToRetreat = true
	}

	if needToRetreat {
		if (C.Type != game.CreepMummy) && s.Can(game.CardStun) && (C.Stun == 0) {
			//println(`STUN BEFORE RETREAT`)
			return game.CardStun
		}

		//println(`RETREAT`)
		return game.CardRetreat
	}

	return card
}

func getCreepyDamage(creepType game.CreepType, max bool) int {
	return creepyDps[creepType][max]
}

func getDamage(s game.State, max bool) (int, game.CardType) {
	A := s.Avatar
	C := s.Creep

	myDps := map[game.CardType]map[bool]int{
		game.CardAttack:      {false: 2, true: 4},
		game.CardMagicArrow:  {false: 3, true: 3},
		game.CardPowerAttack: {false: 4, true: 5},
		game.CardFirebolt:    {false: 4, true: 6},
	}

	dps := func(cardType game.CardType) int {
		return myDps[cardType][max]
	}

	dpst := func(cardType game.CardType) (int, game.CardType) {
		return dps(cardType), cardType
	}

	stuntCount := s.Deck[game.CardStun].Count
	fireboltCount := s.Deck[game.CardFirebolt].Count

	isDragon := C.Type == game.CreepDragon
	nextDragon := s.NextCreep == game.CreepDragon

	// parry
	if s.Can(game.CardParry) && C.Type != game.CreepFairy {
		return getCreepyDamage(C.Type, false), game.CardParry
	}

	// firebolt for mummy
	if (C.Type == game.CreepMummy) && s.Can(game.CardFirebolt) && (A.MP >= 7 || max) {
		return dpst(game.CardFirebolt) // mp 3
	}

	// stun
	if s.Can(game.CardStun) && (C.Stun == 0) && (isDragon || stuntCount >= 2) {
		switch C.Type {
		case game.CreepLion, game.CreepFairy, game.CreepMummy, game.CreepDragon:
			return 0, game.CardStun
		default:
		}
	}

	// power attack
	if s.Can(game.CardPowerAttack) {
		switch C.Type {
		case game.CreepMummy, game.CreepDragon:
			return dpst(game.CardPowerAttack)
		case game.CreepImp, game.CreepLion, game.CreepFairy:
			if s.Turn <= 3 && A.MP >= A.MaxMP-5 {
				return dpst(game.CardPowerAttack)
			}
		default:
		}
	}

	// magic arrow
	if s.Can(game.CardMagicArrow) {
		canCast := max || (A.MP >= 10)

		if canCast {
			switch C.Type {
			case game.CreepFairy, game.CreepMummy, game.CreepLion:
				return dpst(game.CardMagicArrow) // mp 1
			default:
			}
		}
	}

	// firebolt
	if s.Can(game.CardFirebolt) {
		if (fireboltCount >= 3 || A.MP >= 6) && !nextDragon {
			return dpst(game.CardFirebolt) // mp 3
		}
	}

	return dpst(game.CardAttack)
}
