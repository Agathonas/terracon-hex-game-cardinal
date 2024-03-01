// system_turn.go
package system

import (
	"fmt"

	"github.com/argus-labs/starter-game-template/cardinal/component"
	"github.com/argus-labs/starter-game-template/cardinal/msg"
	"pkg.world.dev/world-engine/cardinal"
	"pkg.world.dev/world-engine/cardinal/filter"
	"pkg.world.dev/world-engine/cardinal/message"
	"pkg.world.dev/world-engine/cardinal/types"
)

// TurnSystem manages the progression of turns and rounds.
func TurnSystem(world cardinal.WorldContext) error {
	// Initialize the first turn if necessary
	if err := initializeFirstTurn(world); err != nil {
		return err
	}

	// Process the active player's turn
	if err := processActivePlayerTurn(world); err != nil {
		return err
	}

	// Handle end turn messages
	return handleEndTurnMessages(world)
}

func initializeFirstTurn(world cardinal.WorldContext) error {
	search := cardinal.NewSearch(world, filter.Exact(component.Turn{}))
	count, err := search.Count()
	if err != nil {
		return fmt.Errorf("failed to check for existing turn components: %w", err)
	}

	if count == 0 {
		firstPlayerID := types.EntityID(1) // Assuming player IDs start at 1 and are sequential.
		turnComponent := component.Turn{
			TurnID:       1,
			ActivePlayer: firstPlayerID,
			MovedArmies:  make(map[types.EntityID]bool),
		}
		if _, err := cardinal.Create(world, turnComponent); err != nil {
			return fmt.Errorf("failed to create the first turn component: %w", err)
		}
	}

	return nil
}

func processActivePlayerTurn(world cardinal.WorldContext) error {
	turnComponent, err := getTurnComponent(world)
	if err != nil {
		return err
	}

	// Fetch the Player component for the active player
	playerComponent, err := cardinal.GetComponent[component.Player](world, turnComponent.ActivePlayer)
	if err != nil {
		return fmt.Errorf("failed to get player component for entity %d: %w", turnComponent.ActivePlayer, err)
	}

	// Now, you can safely check if it's the active player's turn
	if playerComponent.IsActiveTurn {
		if err := handleArmyMovements(world, turnComponent.ActivePlayer); err != nil {
			return err
		}

		allMoved, err := allArmiesMoved(world, turnComponent) // Corrected to capture both return values
		if err != nil {
			return fmt.Errorf("error checking if all armies have moved: %w", err)
		}

		if allMoved {
			return switchToNextPlayer(world, turnComponent)
		}
	}

	return nil
}

func handleEndTurnMessages(world cardinal.WorldContext) error {
	// Use EachMessage to iterate over messages of type EndTurnMsg.
	return cardinal.EachMessage[msg.EndTurnMsg, msg.EndTurnMsgReply](world,
		func(txData message.TxData[msg.EndTurnMsg]) (msg.EndTurnMsgReply, error) {
			turnComponent, err := getTurnComponent(world)
			if err != nil {
				return msg.EndTurnMsgReply{}, err
			}

			// Directly access Msg properties without calling Msg().
			if txData.Msg.PlayerID != turnComponent.ActivePlayer {
				return msg.EndTurnMsgReply{Success: false, Message: "It's not your turn"}, nil
			}

			if err := switchToNextPlayer(world, turnComponent); err != nil {
				return msg.EndTurnMsgReply{Success: false, Message: "Failed to end turn"}, err
			}

			return msg.EndTurnMsgReply{Success: true, Message: "Turn ended successfully"}, nil
		})
}

func getTurnComponent(world cardinal.WorldContext) (*component.Turn, error) {
	var turnComponent *component.Turn
	found := false

	search := cardinal.NewSearch(world, filter.Exact(component.Turn{}))
	err := search.Each(func(id types.EntityID) bool {
		var err error
		turnComponent, err = cardinal.GetComponent[component.Turn](world, id)
		if err != nil {
			return false // Stop iteration on error
		}
		found = true
		return false // Stop iteration after finding the first component
	})

	if err != nil {
		return nil, fmt.Errorf("error during search: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("no turn component found")
	}

	return turnComponent, nil
}

func handleArmyMovements(world cardinal.WorldContext, playerID types.EntityID) error {
	search := cardinal.NewSearch(world, filter.Exact(component.Army{PlayerID: playerID}))

	err := search.Each(func(armyID types.EntityID) bool {
		armyComponent, err := cardinal.GetComponent[component.Army](world, armyID)
		if err != nil {
			return false // Error occurred, stop iteration
		}

		if !armyComponent.HasMoved {
			armyComponent.HasMoved = true
			if err := cardinal.SetComponent(world, armyID, armyComponent); err != nil {
				return false // Error occurred, stop iteration
			}
		}
		return true // Continue iteration
	})

	if err != nil {
		return fmt.Errorf("failed to process army movements: %w", err)
	}

	return nil
}

func allArmiesMoved(world cardinal.WorldContext, turnComponent *component.Turn) (bool, error) {
	// Retrieve all armies belonging to the active player.
	search := cardinal.NewSearch(world, filter.Exact(component.Army{PlayerID: turnComponent.ActivePlayer}))
	totalArmies := 0
	movedArmies := 0

	err := search.Each(func(armyID types.EntityID) bool {
		armyComponent, err := cardinal.GetComponent[component.Army](world, armyID)
		if err != nil {
			return false // Stop iteration on error
		}

		totalArmies++
		if armyComponent.HasMoved {
			movedArmies++
		}
		return true // Continue iteration
	})

	if err != nil {
		return false, fmt.Errorf("failed to process army movements: %w", err)
	}

	// Check if the number of moved armies matches the total number of armies
	return totalArmies == movedArmies, nil
}

func switchToNextPlayer(world cardinal.WorldContext, turnComponent *component.Turn) error {
	nextPlayerID := getNextPlayerID(turnComponent.ActivePlayer)
	turnComponent.ActivePlayer = nextPlayerID
	if err := cardinal.SetComponent(world, turnComponent.ActivePlayer, turnComponent); err != nil {
		return fmt.Errorf("failed to update the turn component for next player: %w", err)
	}

	return nil
}

// Assuming player IDs are sequential and wrap from the last player back to the first.
func getNextPlayerID(currentPlayerID types.EntityID) types.EntityID {
	nextPlayerID := currentPlayerID + 1
	if nextPlayerID > 4 { // Assuming 4 players in total
		nextPlayerID = 1 // Loop back to the first player
	}
	return nextPlayerID
}
