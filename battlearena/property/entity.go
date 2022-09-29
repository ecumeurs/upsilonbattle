package property

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

func DefaultHP() *DefaultIntCounterProperty {
	return MakeIntCounterProperty(HP, 10, 10, FriendlyController, Character)
}

func DefaultMovement() *DefaultIntCounterProperty {
	return MakeIntCounterProperty(Movement, 5, 5, FriendlyController, Character)
}

func DefaultSP() *DefaultIntCounterProperty {
	return MakeIntCounterProperty(SP, 10, 10, FriendlyController, Character)
}

func DefaultMP() *DefaultIntCounterProperty {
	return MakeIntCounterProperty(MP, 10, 10, FriendlyController, Character)
}

func DefaultAttack() *DefaultIntProperty {
	return MakeIntProperty(Attack, 1, FriendlyController, Character)
}

func DefaultDefence() *DefaultIntProperty {
	return MakeIntProperty(Defence, 0, FriendlyController, Character)
}

func DefaultJumpHeight() *DefaultIntProperty {
	return MakeIntProperty(JumpHeight, 2, FriendlyController, Character)
}

func DefaultIsDying() *DefaultBoolProperty {
	return MakeBoolProperty(IsDying, false, Public, Character)
}

func DefaultHasMoved() *DefaultBoolProperty {
	return MakeBoolProperty(HasMoved, false, GameMaster, Character)
}

func DefaultHasAttacked() *DefaultBoolProperty {
	return MakeBoolProperty(HasAttacked, false, GameMaster, Character)
}

func DefaultAttackRange() *DefaultIntProperty {
	return MakeIntProperty(AttackRange, 1, FriendlyController, Character)
}

func DefaultShield() *DefaultIntProperty {
	return MakeIntProperty(Shield, 0, FriendlyController, Character)
}

func DefaultPoison() *DefaultIntProperty {
	return MakeIntProperty(Poison, 0, FriendlyController, Character)
}

func DefaultStun() *DefaultIntProperty {
	return MakeIntProperty(Stun, 0, FriendlyController, Character)
}
