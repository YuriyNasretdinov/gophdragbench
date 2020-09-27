package yourock

import . "github.com/quasilyte/gophers-and-dragons/game"

const DragonAttack = 5

func Live2(s State) CardType {

	if s.Creep.Type == CreepCheepy {
		return CardAttack
	}

	if s.Avatar.HP <= 7 {
		return CardRetreat
	}

	if s.Avatar.HP <= 20 {
		if s.Can(CardHeal) {
			return CardHeal
		}
	}

	if s.Creep.Type == CreepDragon && s.Creep.IsFull() {
		turnsLeftToLive := (s.Avatar.HP+s.Deck[CardHeal].Count*7)/DragonAttack + s.Deck[CardStun].Count + s.Deck[CardParry].Count
		dragonHPLeft := 30.0
		powerAttLeft := s.Deck[CardPowerAttack].Count + s.Deck[CardParry].Count

		for i := 0; i < turnsLeftToLive; i++ {
			if powerAttLeft > 0 {
				powerAttLeft--
				dragonHPLeft -= 4.5
				continue
			}

			dragonHPLeft -= 3.0
		}

		if dragonHPLeft > 0.0 {
			return CardRetreat
		}
	}

	if s.Creep.Type != CreepFairy && !s.Creep.IsStunned() && s.Can(CardParry) {
		return CardParry
	}

	if !s.Creep.IsStunned() && s.Can(CardStun) {
		return CardStun
	}

	if s.Can(CardPowerAttack) {
		return CardPowerAttack
	}

	if s.Creep.Type == CreepDragon {
		return CardAttack
	}

	if s.Creep.Type == CreepMummy && s.Avatar.MP >= FireboltMP+HealMP && s.Can(CardFirebolt) {
		return CardFirebolt
	}

	if s.Creep.Type == CreepFairy && s.Avatar.MP >= MagicArrowMP+HealMP && s.Can(CardMagicArrow) {
		return CardMagicArrow
	}

	return CardAttack
}
