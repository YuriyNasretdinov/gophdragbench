package yourock

import (
	. "github.com/quasilyte/gophers-and-dragons/game"
)

type deck struct {
	powerAttack int
	firebolt    int
	stun        int
	heal        int
	parry       int
}

type creep struct {
	typ       CreepType
	hp        int
	isFull    bool
	isStunned bool
}

type creepProbabilitiesPercent struct {
	cheepy int
	imp    int
	fairy  int
	lion   int
	mummy  int
	dragon int
}

// predict the probability of a creep based on a round,
// from https://github.com/quasilyte/gophers-and-dragons/blob/master/wasm/sim/sim.go#L151
func creepProbabilities(round int) creepProbabilitiesPercent {
	switch {
	case round == 1:
		return creepProbabilitiesPercent{cheepy: 100}
	case round == 2:
		return creepProbabilitiesPercent{imp: 100}
	case round <= 5:
		return creepProbabilitiesPercent{
			fairy:  10,
			lion:   40,
			imp:    20,
			cheepy: 30,
		}
	case round < 10:
		return creepProbabilitiesPercent{
			mummy:  30,
			fairy:  20,
			lion:   20,
			imp:    20,
			cheepy: 10,
		}
	}

	return creepProbabilitiesPercent{
		dragon: 100,
	}
}

type state struct {
	hp    int
	mp    int
	round int
	next  CreepType
	c     creep
	d     deck
}

func (s state) can(c CardType) bool {
	switch c {
	case CardAttack, CardRetreat:
		return true
	case CardMagicArrow:
		return s.mp >= MagicArrowMP
	case CardPowerAttack:
		return s.d.powerAttack > 0
	case CardFirebolt:
		return s.d.firebolt > 0 && s.mp >= FireboltMP
	case CardStun:
		return s.d.stun > 0
	case CardRest:
		return s.mp >= RestMP
	case CardHeal:
		return s.d.heal > 0 && s.mp >= HealMP
	case CardParry:
		return s.d.parry > 0
	}

	return false
}

func flattenDeck(d map[CardType]Card) deck {
	return deck{
		powerAttack: d[CardPowerAttack].Count,
		firebolt:    d[CardFirebolt].Count,
		stun:        d[CardStun].Count,
		heal:        d[CardHeal].Count,
		parry:       d[CardParry].Count,
	}
}

func firstHalf(s state) CardType {
	minHealHP := 15
	// +1 average score if we prioritize healing when we are going to receive no damage
	if s.c.typ == CreepCheepy && s.c.isFull {
		minHealHP = 25
	}

	// +5 average score, +5% dragon kills
	if s.hp <= minHealHP && s.can(CardHeal) {
		return CardHeal
	}

	// +0.5 average score
	if s.c.hp == 3 && s.can(MagicArrowMP) {
		return CardMagicArrow
	}

	return CardAttack
}

func smartCanSurviveDragon(s state) bool {
	dragonHP := float64(s.c.hp)
	myHP := float64(s.hp)

	manaLeft := s.hp
	healsLeft := s.d.heal
	parriesLeft := s.d.parry
	powersLeft := s.d.powerAttack
	stunsLeft := s.d.stun

	for dragonHP > 0 && myHP > 0 {
		if myHP <= 25 && healsLeft > 0 && manaLeft >= HealMP {
			myHP += HealAvg - DragonDmgAvg
			healsLeft--
			manaLeft -= HealMP
			continue
		}

		if parriesLeft > 0 {
			dragonHP -= DragonDmgAvg
			continue
		}

		if stunsLeft > 0 {
			stunsLeft--
			if powersLeft > 0 {
				powersLeft--
				dragonHP -= PowerAttackDmgAvg
			} else {
				dragonHP -= AttackDmgAvg
			}
			continue
		}

		myHP -= DragonDmgAvg

		if powersLeft > 0 {
			powersLeft--
			dragonHP -= PowerAttackDmgAvg
			continue
		}

		dragonHP -= AttackDmgAvg
	}

	return myHP > 0 && dragonHP < 0
}

func smartFightDragon(s state) CardType {
	if s.c.isFull && !smartCanSurviveDragon(s) {
		if s.can(CardHeal) {
			return CardHeal
		}

		return CardRetreat
	}

	// +2 average score, +3% dragon win rate
	if s.hp <= 25 && s.can(CardHeal) {
		return CardHeal
	}

	if !s.c.isStunned && s.can(CardParry) {
		return CardParry
	}

	if !s.c.isStunned && s.can(CardStun) {
		return CardStun
	}

	if s.can(CardPowerAttack) {
		return CardPowerAttack
	}

	return CardAttack
}

func canSurviveDragonWithoutPowerAttack(s state) bool {
	s.d.powerAttack--
	return smartCanSurviveDragon(s)
}

func canSurviveDragonWithoutStun(s state) bool {
	s.d.stun--
	return smartCanSurviveDragon(s)
}

func canSurviveDragonWithoutParry(s state) bool {
	s.d.parry--
	return smartCanSurviveDragon(s)
}

func canSurviveDragonWithoutFirebolt(s state) bool {
	s.mp -= FireboltMP
	return smartCanSurviveDragon(s)
}

func secondHalf(s state) CardType {
	// +0.6 average score
	if s.round == 9 && s.c.typ == CreepCheepy && s.mp >= s.d.heal*HealMP+RestMP && s.c.isFull {
		return CardRest
	}

	// +3 average score, +3% dragon kills
	if s.hp <= 25 && s.can(CardHeal) {
		return CardHeal
	}

	// +0.5 average score
	if s.c.hp <= 2 && s.can(CardAttack) {
		return CardAttack
	}

	// +0.5 average score
	if s.c.hp == 3 && s.can(MagicArrowMP) {
		return CardMagicArrow
	}

	// use up all firebolts for the last round if any are left
	// +8 score, +6% dragon kills
	if s.can(CardFirebolt) && s.round >= 9 {
		return CardFirebolt
	}

	// +5 score, +10% dragon kills
	if s.can(CardFirebolt) && s.c.typ == CreepMummy {
		return CardFirebolt
	}

	// WTF
	// +2 average score, -2% dragon kills
	if s.can(CardPowerAttack) && canSurviveDragonWithoutPowerAttack(s) {
		return CardPowerAttack
	}

	// can't use stun at all, it is essential for the dragon fight
	// can't use parry too

	return CardAttack
}

func Smart(gs State) CardType {
	s := state{
		hp:    gs.Avatar.HP,
		mp:    gs.Avatar.MP,
		next:  gs.NextCreep,
		round: gs.Round,
		d:     flattenDeck(gs.Deck),
		c: creep{
			typ:       gs.Creep.Type,
			hp:        gs.Creep.HP,
			isFull:    gs.Creep.IsFull(),
			isStunned: gs.Creep.IsStunned(),
		},
	}

	// some heuristics for the first half of the game
	if gs.Round <= 5 {
		return firstHalf(s)
	}

	// a fully-simulated dragon fight
	if s.c.typ == CreepDragon {
		return smartFightDragon(s)
	}

	// the semi-simulated second half
	return secondHalf(s)
}
