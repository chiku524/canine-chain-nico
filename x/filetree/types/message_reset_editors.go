package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgResetEditors = "reset_editors"

var _ sdk.Msg = &MsgResetEditors{}

func NewMsgResetEditors(creator string, address string, fileowner string) *MsgResetEditors {
	return &MsgResetEditors{
		Creator:   creator,
		Address:   address,
		FileOwner: fileowner,
	}
}

func (msg *MsgResetEditors) Route() string {
	return RouterKey
}

func (msg *MsgResetEditors) Type() string {
	return TypeMsgResetEditors
}

func (msg *MsgResetEditors) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgResetEditors) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgResetEditors) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Address == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
			"invalid address: %s", msg.Address)
	}
	if msg.FileOwner == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
			"invalid file owner: %s", msg.FileOwner)
	}

	return nil
}
