package component

import "pkg.world.dev/world-engine/cardinal/types"

// CityInfoComponent represents the state and attributes of a city on the map.
type CityInfoComponent struct {
	CityID             int            `json:"cityId"`
	Type               string         `json:"type"`  // Capital or Regular
	Owner              types.EntityID `json:"owner"` // Player ID who owns the city
	ArmyProductionRate int            `json:"armyProductionRate"`
	Defenses           int            `json:"defenses"`
	HexQ               int            `json:"hexQ"`
	HexR               int            `json:"hexR"`
}

// Name returns the name of the component.
func (c CityInfoComponent) Name() string {
	return "CityInfo"
}
