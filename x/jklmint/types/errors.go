package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/jklmint module sentinel errors
var (
	ErrCannotParseFloat = errorsmod.Register(ModuleName, 1101, "cannot parse float")
	ErrZeroDivision     = errorsmod.Register(ModuleName, 1102, "cannot use zero value for division")
)
