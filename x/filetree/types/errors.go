package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/filetree module sentinel errors
var (
	ErrNoAccess           = errorsmod.Register(ModuleName, 1101, "wrong permissions for file")
	ErrFileNotFound       = errorsmod.Register(ModuleName, 1102, "file not found")
	ErrCantMarshall       = errorsmod.Register(ModuleName, 1103, "cannot marshall data into json")
	ErrCantUnmarshall     = errorsmod.Register(ModuleName, 1104, "cannot unmarshall data from json")
	ErrPubKeyNotFound     = errorsmod.Register(ModuleName, 1105, "user's public key not found. Account not inited or wrong address")
	ErrParentFileNotFound = errorsmod.Register(ModuleName, 1106, "Parent folder does not exist")
	ErrCannotWrite        = errorsmod.Register(ModuleName, 1107, "You are not permitted to write to this folder")
	ErrNoViewingAccess    = errorsmod.Register(ModuleName, 1108, "You do not have viewing access. Failed to decrypt.")
	ErrCannotDelete       = errorsmod.Register(ModuleName, 1110, "You are not authorized to delete this file")
	ErrNotOwner           = errorsmod.Register(ModuleName, 1111, "Not permitted to remove or reset edit/view access. You are not the owner of this file")
	ErrCantGiveAway       = errorsmod.Register(ModuleName, 1112, "You do not own this file and cannot give it away")
	ErrAlreadyExists      = errorsmod.Register(ModuleName, 1113, "Proposed new owner already has a file set with this path name. No duplicates allowed.")
	ErrCannotAllowEdit    = errorsmod.Register(ModuleName, 1114, "Unauthorized. Only the owner can add an editor.")
	ErrCannotAllowView    = errorsmod.Register(ModuleName, 1115, "Unauthorized. Only the owner can add a viewer.")
	ErrMissingAESKey      = errorsmod.Register(ModuleName, 1116, "AES IV and key required")
	ErrIdsKeysLenMismatch = errorsmod.Register(ModuleName, 1117, "ids and keys must have the same number of comma-separated entries")
)
