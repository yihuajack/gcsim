package unforged

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	combat.RegisterWeaponFunc("the unforged", weapon)
}

func weapon(c core.Character, s core.Sim, log core.Logger, r int, param map[string]int) {

	shd := .15 + float64(r)*.05
	s.AddShieldBonus(func() float64 {
		return shd
	})

	stacks := 0
	icd := 0
	duration := 0

	s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
		if ds.ActorIndex != c.CharIndex() {
			return
		}
		if icd > s.Frame() {
			return
		}
		if duration < s.Frame() {
			stacks = 0
		}
		stacks++
		if stacks > 5 {
			stacks = 0
		}
		icd = s.Frame() + 18

	}, fmt.Sprintf("memory-dust-%v", c.Name()))

	val := make([]float64, core.EndStatType)
	atk := 0.03 + 0.01*float64(r)
	c.AddMod(core.CharStatMod{
		Key:    "memory",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if duration > s.Frame() {
				val[core.ATKP] = atk * float64(stacks)
				if s.IsShielded() {
					val[core.ATKP] *= 2
				}
				return val, true
			}
			stacks = 0
			return nil, false
		},
	})

}