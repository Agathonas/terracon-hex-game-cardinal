// msg_end_turn.go
package msg

import "pkg.world.dev/world-engine/cardinal/types"

// EndTurnMsg represents a request to end the current player's turn.
type EndTurnMsg struct {
	PlayerID types.EntityID // The ID of the player ending their turn.
}

// EndTurnMsgReply defines the response returned after processing an EndTurnMsg.
type EndTurnMsgReply struct {
	Success bool   // Whether the turn was successfully ended.
	Message string // Additional information or error message.
}
