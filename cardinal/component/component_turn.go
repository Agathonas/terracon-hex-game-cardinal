package component

import "pkg.world.dev/world-engine/cardinal/types"

type Turn struct {
	TurnID       int                     // A unique identifier for the turn.
	ActivePlayer types.EntityID          // The ID of the player whose turn it is.
	MovedArmies  map[types.EntityID]bool // A map of army IDs to a boolean indicating if they have moved this turn.
}

func (Turn) Name() string {
	return "Turn"
}
