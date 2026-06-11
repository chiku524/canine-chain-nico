package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/notifications module sentinel errors
var (
	ErrCantUnmarshall         = errorsmod.Register(ModuleName, 1101, "cannot unmarshall from JSON")
	ErrBlockedSender          = errorsmod.Register(ModuleName, 1102, "you are a blocked sender")
	ErrNotificationAlreadySet = errorsmod.Register(ModuleName, 1103, "notification already set")
	ErrNotificationNotFound   = errorsmod.Register(ModuleName, 1105, "notification does not exist")
	ErrNotNotificationOwner   = errorsmod.Register(ModuleName, 1106, "you do not own this notification")

	ErrInvalidContents = errorsmod.Register(ModuleName, 1110, "contents must be valid JSON")
)
