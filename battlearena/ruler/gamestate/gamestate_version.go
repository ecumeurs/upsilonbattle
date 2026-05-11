package gamestate

// UpdateVersion synchronizes the bit-packed version counter with current turn and action indices.
// @spec-link [[mech_version_bit_packing]]
func (gs *GameState) UpdateVersion() {
	gs.Version = (int64(gs.TurnIndex) << 32) | int64(gs.ActionIndex)
}

// IncVersion is a wrapper that increments the action index and updates the packed version.
func (gs *GameState) IncVersion() {
	gs.IncAction()
}

// IncAction increments the action index for the current turn and updates the packed version.
func (gs *GameState) IncAction() {
	gs.ActionIndex++
	gs.UpdateVersion()
}

// IncTurn increments the global turn index, resets the action index, and updates the packed version.
func (gs *GameState) IncTurn() {
	gs.TurnIndex++
	gs.ActionIndex = 0
	gs.UpdateVersion()
}

// GetTurn extracts the 32-bit turn index from the bit-packed Version field.
func (gs *GameState) GetTurn() uint32 {
	return uint32(gs.Version >> 32)
}

// GetAction extracts the 32-bit action index from the bit-packed Version field.
func (gs *GameState) GetAction() uint32 {
	return uint32(gs.Version & 0xFFFFFFFF)
}
