package property

type SkillProperties string

const (
	Behavior SkillProperties = "Behavior" // Skill broad category: Direct, Reaction, Passive, Counter; Absence means Direct
)

type TargetingProperties string

const (
	Range              TargetingProperties = "Range"              // Range of the Skill // Absence means 1
	Zone               TargetingProperties = "Zone"               // Area of Effect // Absence means 1 tile effect
	TargetNumber       TargetingProperties = "TargetNumber"       // Absence means all targets within the targeted zone.
	Accuracy           TargetingProperties = "Accuracy"           // Absence means 100%
	Dodge              TargetingProperties = "Dodge"              // Absence means 0%
	Parry              TargetingProperties = "Parry"              // Absence means 0%
	TargetType         TargetingProperties = "TargetType"         // Entity, Tile, Both, Self
	TargetingMechanics TargetingProperties = "TargetingMechanics" // Anywhere, Line of Sight, and maybe other mechanics later.
)

type EffectProperties string

const (
	Damage          EffectProperties = "Damage"          // Absence means 0
	Heal            EffectProperties = "Heal"            // Absence means 0
	ShieldPower     EffectProperties = "Shield"          // Absence means 0 , can be negative or positive.
	StunPower       EffectProperties = "Stun"            // Absence means 0 , can be negative or positive.
	StunChance      EffectProperties = "StunChance"      // Absence means 0%
	CriticialChance EffectProperties = "CriticialChance" // Absence means 0%
	CriticialDamage EffectProperties = "CriticialDamage" // Absence means 0%
	Duration        EffectProperties = "Duration"        // Absence means 0
	PoisonPower     EffectProperties = "Poison"          // Absence means 0 , can be negative or positive.
	PoisonChance    EffectProperties = "PoisonChance"    // Absence means 0%
)

type CostProperties string

const (
	Delay    CostProperties = "Delay"    // Absence means 500
	HPLeech  CostProperties = "HPLeech"  // Absence means 0
	MPLeech  CostProperties = "MPLeech"  // Absence means 0
	SPLeech  CostProperties = "SPLeech"  // Absence means 0
	Cooldown CostProperties = "Cooldown" // Absence means 3 turns
)

// Prepare default Properties.

func DefaultTargetNumber() *DefaultIntProperty {
	return MakeIntProperty(TargetNumber, 0, FriendlyController, Skill)
}

func DefaultAccuracy() *DefaultIntProperty {
	return MakeIntProperty(Accuracy, 100, FriendlyController, Skill)
}

func DefaultDodge() *DefaultIntProperty {
	return MakeIntProperty(Dodge, 0, FriendlyController, Skill)
}

func DefaultParry() *DefaultIntProperty {
	return MakeIntProperty(Parry, 0, FriendlyController, Skill)
}

func DefaultDamage() *DefaultIntProperty {
	return MakeIntProperty(Damage, 0, FriendlyController, Skill)
}

func DefaultHeal() *DefaultIntProperty {
	return MakeIntProperty(Heal, 0, FriendlyController, Skill)
}

func DefaultShieldPower() *DefaultIntProperty {
	return MakeIntProperty(ShieldPower, 0, FriendlyController, Skill)
}

func DefaultStunPower() *DefaultIntProperty {
	return MakeIntProperty(StunPower, 0, FriendlyController, Skill)
}

func DefaultStunChance() *DefaultIntProperty {
	return MakeIntProperty(StunChance, 0, FriendlyController, Skill)
}

func DefaultCriticialChance() *DefaultIntProperty {
	return MakeIntProperty(CriticialChance, 0, FriendlyController, Skill)
}

func DefaultCriticialDamage() *DefaultFloatProperty {
	return MakeFloatProperty(CriticialDamage, 0, FriendlyController, Skill)
}

func DefaultDuration() *DefaultIntCounterProperty {
	return MakeIntCounterProperty(Duration, 0, 0, FriendlyController, Skill)
}

func DefaultPoisonPower() *DefaultIntProperty {
	return MakeIntProperty(PoisonPower, 0, FriendlyController, Skill)
}

func DefaultPoisonChance() *DefaultIntProperty {
	return MakeIntProperty(PoisonChance, 0, FriendlyController, Skill)
}

func DefaultDelay() *DefaultIntCounterProperty {
	return MakeIntCounterProperty(Delay, 0, 500, FriendlyController, Skill)
}

func DefaultHPLeech() *DefaultIntProperty {
	return MakeIntProperty(HPLeech, 0, FriendlyController, Skill)
}

func DefaultMPLeech() *DefaultIntProperty {
	return MakeIntProperty(MPLeech, 0, FriendlyController, Skill)
}

func DefaultSPLeech() *DefaultIntProperty {
	return MakeIntProperty(SPLeech, 0, FriendlyController, Skill)
}

func DefaultCooldown() *DefaultIntCounterProperty {
	return MakeIntCounterProperty(Cooldown, 0, 3, FriendlyController, Skill)
}
