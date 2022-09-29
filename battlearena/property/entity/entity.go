package entity

import "github.com/ecumeurs/upsilonbattle/battlearena/property"

type EntityProperties string

const (
	HP       EntityProperties = "HP"       // Absence means 10 /!\ this is a counter (has two value: current, max)
	Movement EntityProperties = "Movement" // Absence means 5 /!\ this is a counter (has two value: current, max), reset at end of turn
	SP       EntityProperties = "SP"       // Absence means 10
	MP       EntityProperties = "MP"       // Absence means 10

	Attack     EntityProperties = "Attack"     // Absence means 1, basic attack
	Defence    EntityProperties = "Defence"    // Absence means 0, basic defense
	JumpHeight EntityProperties = "JumpHeight" // Absence means 2

	// Flags

	IsDying     EntityProperties = "IsDying"     // Absence means false
	HasMoved    EntityProperties = "HasMoved"    // Absence means false
	HasAttacked EntityProperties = "HasAttacked" // Absence means false

	// Buffs and other applied, inherited properties ... ?
	AttackRange EntityProperties = "AttackRange" // Absence means 1, Basic attack range. Altered with items, mostly.

	// Computed at end of turn.
	Shield EntityProperties = "Shield" // Absence means 0
	Poison EntityProperties = "Poison" // Absence means 0, Poisoned state. Sum of all poison damage taken per turn. Mostly a temporary status. Negative PoisonPower will cure
	Stun   EntityProperties = "Stun"   // Absence means 0, Stunned state. Sum of all stun taken per turn. Mostly a temporary status. Negative StunPower will cure
)

// Default Properties Declaration

func DefaultHP() *property.DefaultIntCounterProperty {
	return property.MakeIntCounterProperty(HP, 10, 10, property.FriendlyController, property.Character)
}

func DefaultMovement() *property.DefaultIntCounterProperty {
	return property.MakeIntCounterProperty(Movement, 5, 5, property.FriendlyController, property.Character)
}

func DefaultSP() *property.DefaultIntCounterProperty {
	return property.MakeIntCounterProperty(SP, 10, 10, property.FriendlyController, property.Character)
}

func DefaultMP() *property.DefaultIntCounterProperty {
	return property.MakeIntCounterProperty(MP, 10, 10, property.FriendlyController, property.Character)
}

func DefaultAttack() *property.DefaultIntProperty {
	return property.MakeIntProperty(Attack, 1, property.FriendlyController, property.Character)
}

func DefaultDefence() *property.DefaultIntProperty {
	return property.MakeIntProperty(Defence, 0, property.FriendlyController, property.Character)
}

func DefaultJumpHeight() *property.DefaultIntProperty {
	return property.MakeIntProperty(JumpHeight, 2, property.FriendlyController, property.Character)
}

func DefaultIsDying() *property.DefaultBoolProperty {
	return property.MakeBoolProperty(IsDying, false, property.Public, property.Character)
}

func DefaultHasMoved() *property.DefaultBoolProperty {
	return property.MakeBoolProperty(HasMoved, false, property.GameMaster, property.Character)
}

func DefaultHasAttacked() *property.DefaultBoolProperty {
	return property.MakeBoolProperty(HasAttacked, false, property.GameMaster, property.Character)
}

func DefaultAttackRange() *property.DefaultIntProperty {
	return property.MakeIntProperty(AttackRange, 1, property.FriendlyController, property.Character)
}

func DefaultShield() *property.DefaultIntProperty {
	return property.MakeIntProperty(Shield, 0, property.FriendlyController, property.Character)
}

func DefaultPoison() *property.DefaultIntProperty {
	return property.MakeIntProperty(Poison, 0, property.FriendlyController, property.Character)
}

func DefaultStun() *property.DefaultIntProperty {
	return property.MakeIntProperty(Stun, 0, property.FriendlyController, property.Character)
}

// note: futher properties may be added per entity basis.
func DefaultPropertiesForCharacter() []property.Property {
	return []property.Property{
		property.MakeIntCounterProperty(HP, 10, 10, property.Public, property.Character),
		property.MakeIntCounterProperty(Movement, 10, 10, property.Public, property.Character),
		property.MakeIntProperty(Attack, 3, property.Public, property.Character),
		property.MakeIntProperty(AttackRange, 1, property.Public, property.Character),
		property.MakeIntProperty(Defence, 0, property.Public, property.Character),
		property.MakeIntProperty(JumpHeight, 2, property.Public, property.Character),
	}
}
