package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/rns module sentinel errors
var (
	ErrNoTLD    = errorsmod.Register(ModuleName, 1100, "could not extract the tld from the name provided")
	ErrReserved = errorsmod.Register(ModuleName, 1101, "tld is reserved by the system")
)
