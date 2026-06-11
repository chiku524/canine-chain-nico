package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/storage module sentinel errors
var (
	ErrDivideByZero        = errorsmod.Register(ModuleName, 1110, "cannot divide by zero")
	ErrProviderNotFound    = errorsmod.Register(ModuleName, 1111, "provider not found please init your provider")
	ErrNotValidTotalSpace  = errorsmod.Register(ModuleName, 1112, "not valid total space please enter total number of bytes to provide")
	ErrDealNotFound        = errorsmod.Register(ModuleName, 1114, "cannot find active deal")
	ErrFormNotFound        = errorsmod.Register(ModuleName, 1115, "cannot find attestation form")
	ErrAttestInvalid       = errorsmod.Register(ModuleName, 1116, "cannot attest to form")
	ErrAttestAlreadyExists = errorsmod.Register(ModuleName, 1117, "attest already exists")
	ErrCannotVerifyProof   = errorsmod.Register(ModuleName, 1118, "cannot verify Proof")
	ErrNoCid               = errorsmod.Register(ModuleName, 1119, "cid does not exist")
	ErrProviderExists      = errorsmod.Register(ModuleName, 1120, "provider already exists")
	ErrBadProofInput       = errorsmod.Register(ModuleName, 1121, "bad proof input")
)
