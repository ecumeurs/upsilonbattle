package behavior

// CompositeMode controls how a CompositeBehavior combines its children.
type CompositeMode int

const (
	// CompositeOR returns the first non-Pass decision from its children.
	// Semantics: "do the first thing any behavior thinks we should do."
	CompositeOR CompositeMode = iota
	// CompositeAND requires ALL children to return the same Decision type.
	// If any child disagrees, the composite falls back to DecisionPass.
	// Semantics: "only act if all behaviors agree."
	CompositeAND
)

// CompositeBehavior combines multiple Behaviors into one.
// Children are evaluated in order.
//
// @spec-link [[mech_composite_behavior]]
type CompositeBehavior struct {
	Mode     CompositeMode
	Children []Behavior
}

// NewCompositeOR creates a CompositeBehavior that acts on the first non-Pass child decision.
func NewCompositeOR(children ...Behavior) *CompositeBehavior {
	return &CompositeBehavior{Mode: CompositeOR, Children: children}
}

// NewCompositeAND creates a CompositeBehavior that only acts if all children agree.
func NewCompositeAND(children ...Behavior) *CompositeBehavior {
	return &CompositeBehavior{Mode: CompositeAND, Children: children}
}

// OnTurn implements Behavior.
func (b *CompositeBehavior) OnTurn(ctx GameContext) Decision {
	switch b.Mode {
	case CompositeOR:
		for _, child := range b.Children {
			d := child.OnTurn(ctx)
			if d.Type != DecisionPass {
				return d
			}
		}
		return Decision{Type: DecisionPass}

	case CompositeAND:
		if len(b.Children) == 0 {
			return Decision{Type: DecisionPass}
		}
		first := b.Children[0].OnTurn(ctx)
		if first.Type == DecisionPass {
			return Decision{Type: DecisionPass}
		}
		for _, child := range b.Children[1:] {
			d := child.OnTurn(ctx)
			if d.Type != first.Type {
				return Decision{Type: DecisionPass}
			}
		}
		return first
	}

	return Decision{Type: DecisionPass}
}
