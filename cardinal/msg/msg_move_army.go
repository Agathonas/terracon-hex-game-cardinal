package msg

import "pkg.world.dev/world-engine/cardinal/types"

type MoveArmyMsg struct {
	ArmyID       types.EntityID
	NewLocationQ int
	NewLocationR int
}

type MoveArmyMsgReply struct {
	Success bool
	Message string
}
