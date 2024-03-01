package component

import "pkg.world.dev/world-engine/cardinal/types"

// Army represents the state and attributes of a player's army.
type Army struct {
	ArmyID        int            `json:"armyId"`    // Unique identifier for the army.
	PlayerID      types.EntityID `json:"playerId"`  // ID of the player who owns the army, using EntityID type.
	Strength      int            `json:"strength"`  // The combat strength of the army.
	LocationQ     int            `json:"locationQ"` // The Q coordinate of the army's location.
	LocationR     int            `json:"locationR"` // The R coordinate of the army's location.
	MovementRange int            `json:"movementRange"`
	HasMoved      bool           `json:"hasMoved"` // New field to track if the army has moved this turn.

}

func (Army) Name() string {
	return "Army"
}
