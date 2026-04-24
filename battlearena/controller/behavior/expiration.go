package behavior

// ExpirationBehavior is the behavior for temporary entities that only exist for N turns.
// On their turn, they perform no meaningful action — they simply pass.
// Their actual work is done via positional effects; the entity itself is just a lifecycle anchor.
// Expiry is handled by GameState.RemoveEntity called from EndOfTurn when EntityDuration reaches 0.
//
// @spec-link [[mech_entity_expiration]]
// @spec-link [[mech_behavior_system]]
type ExpirationBehavior struct{}

func (b *ExpirationBehavior) OnTurn(_ GameContext) Decision {
	return Decision{Type: DecisionPass}
}
