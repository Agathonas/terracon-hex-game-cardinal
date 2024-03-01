package system

import (
	"fmt"
	"math/rand"
	"time"

	comp "github.com/argus-labs/starter-game-template/cardinal/component"
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
)

const (
	MapWidth  = 11
	MapHeight = 22
)

func HexMapSystem(world cardinal.WorldContext) error {
	search := cardinal.NewSearch(world, filter.Exact(comp.MapInitialized{}))
	count, err := search.Count()
	if err != nil {
		return fmt.Errorf("failed to check map initialization: %w", err)
	}
	if count > 0 {
		return nil
	}

	for q := 0; q < MapWidth; q++ {
		for r := 0; r < MapHeight; r++ {
			hexComponent := comp.NewHex(q, r)
			_, err := cardinal.Create(world, hexComponent)
			if err != nil {
				return fmt.Errorf("failed to create hex tile entity: %w", err)
			}
		}
	}

	playerNicknames := []string{"Player1", "Player2", "Player3", "Player4"}

	capitalPositions := []struct{ q, r int }{
		{1, 1}, {MapWidth - 2, 1}, {1, MapHeight - 2}, {MapWidth - 2, MapHeight - 2},
	}

	rand.Seed(time.Now().UnixNano())

	cityID := 1
	for i, pos := range capitalPositions {
		cityComponent := comp.CityInfoComponent{
			CityID:             cityID,
			Type:               "Capital",
			Owner:              0,
			ArmyProductionRate: 5,
			Defenses:           10,
			HexQ:               pos.q,
			HexR:               pos.r,
		}

		capitalCityEntityID, err := cardinal.Create(world, cityComponent)
		if err != nil {
			return fmt.Errorf("failed to create capital city: %w", err)
		}

		if i < len(playerNicknames) {
			playerComponent := comp.Player{
				Nickname:      playerNicknames[i],
				CapitalCityID: cityID,
				Resources:     100,
			}

			playerEntityID, err := cardinal.Create(world, playerComponent)
			if err != nil {
				return fmt.Errorf("failed to create player entity: %w", err)
			}

			// Set the PlayerID in the Player component to the EntityID of the newly created player entity
			playerComponent.PlayerID = playerEntityID
			if err := cardinal.SetComponent(world, playerEntityID, &playerComponent); err != nil {
				return fmt.Errorf("failed to update player component with PlayerID: %w", err)
			}

			// Update the city owner to be the player
			cityComponent.Owner = playerEntityID
			if err := cardinal.SetComponent(world, capitalCityEntityID, &cityComponent); err != nil {
				return fmt.Errorf("failed to update city owner: %w", err)
			}

			// Create an Army component for the player, positioned at their capital city
			armyComponent := comp.Army{
				ArmyID:    cityID,
				PlayerID:  playerEntityID,
				Strength:  100,
				LocationQ: pos.q,
				LocationR: pos.r,
			}

			_, err = cardinal.Create(world, armyComponent)
			if err != nil {
				return fmt.Errorf("failed to create army entity for player: %w", err)
			}
		}

		cityID++
	}

	numberOfRegularCities := 16
	for i := 0; i < numberOfRegularCities; i++ {
		q := rand.Intn(MapWidth)
		r := rand.Intn(MapHeight)
		isCapital := false
		for _, capPos := range capitalPositions {
			if capPos.q == q && capPos.r == r {
				isCapital = true
				break
			}
		}
		if isCapital {
			continue
		}

		cityComponent := comp.CityInfoComponent{
			CityID:             cityID,
			Type:               "Regular",
			Owner:              0,
			ArmyProductionRate: 3,
			Defenses:           5,
			HexQ:               q,
			HexR:               r,
		}

		_, err := cardinal.Create(world, cityComponent)
		if err != nil {
			return fmt.Errorf("failed to create regular city: %w", err)
		}

		cityID++
	}

	// Mark the hex map as initialized
	_, err = cardinal.Create(world, comp.MapInitialized{})
	if err != nil {
		return fmt.Errorf("failed to mark hex map as initialized: %w", err)
	}

	// Log that the hex map has been initialized
	fmt.Println("Hex Map Initialized")

	return nil
}
