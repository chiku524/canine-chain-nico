package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgChangeOwner = "change_owner"

var _ sdk.Msg = &MsgChangeOwner{}

func NewMsgChangeOwner(creator string, address string, fileOwner string, newOwner string) *MsgChangeOwner {
	return &MsgChangeOwner{
		Creator:   creator,
		Address:   address,
		FileOwner: fileOwner,
		NewOwner:  newOwner,
	}
}

func (msg *MsgChangeOwner) Route() string {
	return RouterKey
}

func (msg *MsgChangeOwner) Type() string {
	return TypeMsgChangeOwner
}

func (msg *MsgChangeOwner) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgChangeOwner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgChangeOwner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Check empty values
	if msg.NewOwner == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
			"invalid new owner: %s", msg.NewOwner)
	}
	if msg.FileOwner == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
			"invalid file owner: %s", msg.FileOwner)
	}
	if msg.Address == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
			"invalid address: %s", msg.Address)
	}

	return nil
}
