package property

// TriggerTypeValue defines when a positional effect fires.
// Used as the string value of the TriggerType SkillProperty on effect definitions.
//
// @spec-link [[mech_trigger_system]]
type TriggerTypeValue string

const (
	// TriggerOnEnter fires when an entity enters the cell.
	TriggerOnEnter TriggerTypeValue = "OnEnter"
	// TriggerOnExit fires when an entity leaves the cell.
	TriggerOnExit TriggerTypeValue = "OnExit"
	// TriggerOnStep fires for every step through the cell (enter + intermediate steps).
	TriggerOnStep TriggerTypeValue = "OnStep"
	// TriggerOnTurn fires at the beginning of each turn while the entity is in the cell.
	TriggerOnTurn TriggerTypeValue = "OnTurn"
	// TriggerOnDeath fires when the entity dies while in the cell.
	TriggerOnDeath TriggerTypeValue = "OnDeath"
)
