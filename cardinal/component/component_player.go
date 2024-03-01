package component

import "pkg.world.dev/world-engine/cardinal/types"

// Player stores the state and attributes related to a player.
type Player struct {
	PlayerID      types.EntityID `json:"playerId"`      // Unique identifier for the player.
	Nickname      string         `json:"nickname"`      // Player's chosen nickname.
	CapitalCityID int            `json:"capitalCityId"` // ID of the player's capital city.
	Resources     int            `json:"resources"`     // Resources like $ETH balance, army points, etc.
	IsActiveTurn  bool           `json:"isActiveTurn"`  // Indicates if it's this player's turn.
}

func (Player) Name() string {
	return "Player"
}
