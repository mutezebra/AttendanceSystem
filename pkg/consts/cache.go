package consts

import "time"

// user
const (
	VerifyCodeExpireTime = 3 * 60 * time.Second
)

const (
	CallEventDoKey   = "CallThing:do"
	CallEventUndoKey = "CallThing:undo"

	SvcCallEventKey = "SvcCallEvent"
)
