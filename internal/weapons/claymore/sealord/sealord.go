package sealord

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

func init() {
	core.RegisterWeaponFunc(keys.LuxuriousSeaLord, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	//Increases Elemental Burst DMG by 12%. When Elemental Burst hits opponents,
	//there is a 100% chance of summoning a huge onrush of tuna that deals 100%
	//ATK as AoE DMG. This effect can occur once every 15s.
	w := &Weapon{}
	r := p.Refine

	//perm burst dmg increase
	burstDmgIncrease := .09 + float64(r)*0.03
	val := make([]float64, attributes.EndStatType)
	val[attributes.DmgP] = burstDmgIncrease
	char.AddAttackMod("luxurious-sea-lord", -1, func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		if atk.Info.AttackTag == combat.AttackTagElementalBurst {
			return val, true
		}
		return nil, false
	})

	tunaDmg := .75 + float64(r)*0.25
	icd := -1
	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		if c.F < icd {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagElementalBurst {
			return false
		}
		icd = c.F + 900 //15 * 60
		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Luxurious Sea-Lord Proc",
			AttackTag:  combat.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       tunaDmg,
		}
		c.QueueAttack(ai, combat.NewDefCircHit(1, false, combat.TargettableEnemy), 0, 1)

		return false
	}, fmt.Sprintf("sealord-%v", char.Base.Name))
	return w, nil
}