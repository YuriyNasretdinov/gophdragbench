package alik

import "github.com/golangconf/gophers-and-dragons/game"

func WiningTactic(s game.State) game.CardType {

	if s.Avatar.HP < 10 {

		if s.Can(game.CardHeal) {
			return game.CardHeal
		}
		if s.Creep.IsStunned() {
			if s.Can(game.CardRest) {
				return game.CardRest
			}
		}
		if !s.Creep.IsStunned() && s.Can(game.CardStun) {
			return game.CardStun
		}

		if s.Creep.Traits.Has(game.TraitSlow) {
			return game.CardRetreat // Otherwise run away
		}

		if s.Avatar.HP <= s.Creep.Damage.High() {

			if !s.Creep.Traits.Has(game.TraitRanged) {
				if s.Can(game.CardParry) && s.Creep.Traits.Has(game.TraitCoward) {
					return game.CardParry
				}
			}

			return game.CardRetreat
		}
	}
	switch s.Creep.Type {
	case game.CreepCheepy:
		if s.Can(game.CardPowerAttack) {
			return game.CardPowerAttack
		}
		return game.CardAttack
	case game.CreepImp:
		return game.CardAttack
	case game.CreepLion:
		if s.Can(game.CardPowerAttack) {
			return game.CardPowerAttack
		}
		return game.CardAttack
	case game.CreepFairy:
		if s.Can(game.CardHeal) {
			return game.CardHeal
		}
		if s.Can(game.CardPowerAttack) {
			return game.CardPowerAttack
		}
		if s.Can(game.CardFirebolt) {
			return game.CardFirebolt
		}
		if s.Can(game.CardMagicArrow) && s.Avatar.MP > 5 {
			return game.CardMagicArrow
		}
		return game.CardAttack
	case game.CreepMummy:
		if s.Can(game.CardHeal) && s.Avatar.HP < 30 {
			return game.CardHeal
		}
		if s.Creep.HP == s.Creep.MaxHP {
			if s.Deck[game.CardFirebolt].Count < 2 || s.Deck[game.CardFirebolt].CardStats.MP*2 > s.Avatar.MP {
				if s.Avatar.HP >= s.Creep.Damage.High()*2 {
					return game.CardRetreat
				}
			}
		}
		if s.Creep.IsStunned() && s.Can(game.CardRest) {
			return game.CardRest
		}

		if s.Can(game.CardPowerAttack) {
			return game.CardPowerAttack
		}
		if s.Can(game.CardFirebolt) {
			return game.CardFirebolt
		}

		return game.CardAttack
	default:
		if s.Can(game.CardHeal) {
			return game.CardHeal
		}

		if s.Can(game.CardParry) {
			return game.CardParry
		}
		if !s.Creep.IsStunned() && s.Can(game.CardStun) {
			return game.CardStun
		}

		var roundsToKillMe = s.Avatar.HP / s.Creep.Damage.High()
		var roundsToKillDragon = ((s.Creep.HP - (s.Deck[game.CardPowerAttack].Count * 4)) / 3) + s.Deck[game.CardPowerAttack].Count

		if s.Deck[game.CardStun].Count > 0 {
			if s.Can(game.CardPowerAttack) {
				return game.CardPowerAttack
			}
			return game.CardAttack
		}

		if roundsToKillMe < roundsToKillDragon {
			return game.CardRetreat
		}

		if s.Avatar.HP < 18 && s.Creep.HP > 2 {
			if !s.Creep.IsStunned() && s.Can(game.CardStun) {
				return game.CardStun
			}

			if s.Creep.IsStunned() {

				if s.Can(game.CardRest) && !s.Can(game.CardHeal) && s.Avatar.HP > 3 {
					return game.CardRest
				}
				if s.Can(game.CardPowerAttack) {
					return game.CardPowerAttack
				}
				return game.CardAttack

			}

			return game.CardRetreat
		}

		if s.Can(game.CardPowerAttack) {
			return game.CardPowerAttack
		}

		return game.CardAttack
	}

}
