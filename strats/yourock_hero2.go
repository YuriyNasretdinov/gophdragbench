package yourock

import . "github.com/golangconf/gophers-and-dragons/game"

const MagicArrowMP = 1
const FireboltMP = 3
const HealMP = 4
const RestMP = 2

const HealAvg = 12.5

const DragonDmgAvg = 5.5
const PowerAttackDmgAvg = 4.5
const AttackDmgAvg = 3

func heroCanSurviveDragon(s State) bool {
	dragonHP := float64(s.Creep.HP) - 4
	myHP := float64(s.Avatar.HP)

	manaLeft := s.Avatar.MP
	healsLeft := s.Deck[CardHeal].Count
	parriesLeft := s.Deck[CardParry].Count
	powersLeft := s.Deck[CardPowerAttack].Count
	stunsLeft := s.Deck[CardStun].Count

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

func fightDragon(s State) CardType {
	if s.Creep.IsFull() && !heroCanSurviveDragon(s) {
		if !s.Creep.IsStunned() && s.Can(CardStun) {
			return CardStun
		}

		if s.Avatar.HP < 40 && s.Can(CardHeal) {
			return CardHeal
		}

		return CardRetreat
	}

	if s.Avatar.HP <= 25 && s.Can(CardHeal) {
		return CardHeal
	}

	if !s.Creep.IsStunned() && s.Can(CardParry) {
		return CardParry
	}

	if !s.Creep.IsStunned() && s.Can(CardStun) {
		return CardStun
	}

	if s.Can(CardPowerAttack) {
		return CardPowerAttack
	}

	return CardAttack
}

func fightMummy(s State) CardType {
	if s.Creep.IsFull() {
		if s.Avatar.MP >= 2*FireboltMP && s.Deck[CardFirebolt].Count >= 2 {
			return CardFirebolt
		}

		return CardRetreat
	}

	if s.Creep.HP <= 2 {
		return CardAttack
	} else if s.Creep.HP <= 3 && s.Avatar.MP >= HealMP+MagicArrowMP {
		return CardMagicArrow
	} else if s.Can(CardFirebolt) {
		return CardFirebolt
	}

	return CardRetreat
}

func fightClaws(s State) CardType {
	if s.Creep.IsFull() {
		potentialHPLeft := float64(s.Avatar.HP)
		mpLeft := s.Avatar.MP
		for i := 0; i < s.Deck[CardHeal].Count && mpLeft >= HealMP; i++ {
			mpLeft -= HealMP
			potentialHPLeft += HealAvg
		}

		// Each turn Claws takes out 3 HP but only until your HP >= 20.
		// This variable estimates how many turns we have to kill Claws until
		// our HP drops below 20 and we would have to retreat.
		maxTurnsToFight := int((potentialHPLeft-20.0)/3.0) + 1
		turnsToKillClaws := s.Creep.HP / AttackDmgAvg // should be precisely 4

		if turnsToKillClaws > maxTurnsToFight {
			return CardRetreat
		}
	}

	if s.Avatar.HP <= 20 && s.Creep.HP > 3 {
		return CardRetreat
	}

	if s.Creep.HP == 3 && s.Avatar.MP >= MagicArrowMP+HealMP && s.Can(CardMagicArrow) {
		return CardMagicArrow
	}

	return CardAttack
}

func shouldFightKubus(s State) bool {
	kubusHP := float64(s.Creep.HP)
	myHP := float64(s.Avatar.HP)

	manaLeft := s.Avatar.MP
	fireboltsLeft := s.Deck[CardFirebolt].Count
	healsLeft := s.Deck[CardHeal].Count
	parriesLeft := s.Deck[CardParry].Count
	powersLeft := s.Deck[CardPowerAttack].Count

	turn := 0

	for kubusHP > 0 && myHP > 0 {
		turn++

		if myHP <= 25 && healsLeft > 0 && manaLeft >= HealMP-1 {
			myHP += HealAvg - float64(turn)
			healsLeft--
			manaLeft -= HealMP - 1
			continue
		}

		if manaLeft >= FireboltMP-1+HealMP && fireboltsLeft > 0 {
			kubusHP -= 5
			manaLeft -= FireboltMP - 1
			fireboltsLeft--
			myHP -= float64(turn)
			continue
		}

		// magic arrow
		if kubusHP == 3 {
			kubusHP -= 3
			myHP -= float64(turn)
			continue
		}

		if turn >= 4 && parriesLeft > 0 {
			kubusHP -= float64(turn)
			parriesLeft--
			continue
		}

		myHP -= float64(turn)

		if powersLeft > 0 {
			powersLeft--
			kubusHP -= PowerAttackDmgAvg
			continue
		}

		kubusHP -= AttackDmgAvg
	}

	return myHP > 15 && kubusHP < 0
}

func fightKubus(s State) CardType {
	if s.Creep.IsFull() && !shouldFightKubus(s) {
		return CardRetreat
	}

	if s.Avatar.MP >= FireboltMP-1+HealMP && s.Can(CardFirebolt) {
		return CardFirebolt
	}

	if s.Creep.HP == 3 && s.Can(CardMagicArrow) {
		return CardMagicArrow
	}

	if s.RoundTurn >= 4 && s.Can(CardParry) {
		return CardParry
	}

	if s.Can(CardPowerAttack) {
		return CardPowerAttack
	}

	if s.Avatar.HP <= 10 || s.RoundTurn >= 6 {
		return CardRetreat
	}

	return CardAttack
}

func ChooseCard(s State) CardType {
	if s.Creep.Type == CreepCheepy {
		if s.Creep.IsFull() && s.Avatar.HP <= 25 && s.Can(CardHeal) {
			return CardHeal
		}

		if s.Creep.IsFull() && s.Avatar.HP <= 20 && s.Can(CardRest) {
			return CardRest
		}

		return CardAttack
	}

	if s.Creep.Type == CreepDragon {
		return fightDragon(s)
	}

	// don't need as many stuns to kill the dragon but these cards
	// might be useful for a tougher creeps
	if !s.Creep.IsStunned() && s.Deck[CardStun].Count >= 6 && (s.Creep.Damage.High() >= 5 || s.Creep.Type == CreepKubus && s.RoundTurn >= 4) {
		return CardStun
	}

	minHealHP := 10
	if s.Creep.Damage.High() <= 4 {
		minHealHP = 25
	} else if s.Creep.Traits.Has(TraitBloodlust) {
		minHealHP = 20
	}

	if s.Avatar.HP <= minHealHP && s.Can(CardHeal) {
		return CardHeal
	}

	if s.Creep.Type == CreepClaws {
		return fightClaws(s)
	}

	if s.Creep.Type == CreepKubus {
		return fightKubus(s)
	}

	if s.Creep.Type == CreepMummy {
		return fightMummy(s)
	}

	if s.Round >= 10 && s.Avatar.MP >= FireboltMP+HealMP && s.Can(CardFirebolt) {
		return CardFirebolt
	}

	if s.Creep.Damage.High() >= 5 && s.Can(CardPowerAttack) {
		return CardPowerAttack
	}

	if s.Creep.HP == 3 && s.Avatar.MP >= MagicArrowMP+HealMP && s.Can(CardMagicArrow) {
		return CardMagicArrow
	}

	return CardAttack
}
