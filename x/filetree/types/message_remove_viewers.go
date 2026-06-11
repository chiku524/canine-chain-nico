package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveViewers = "remove_viewers"

var _ sdk.Msg = &MsgRemoveViewers{}

func NewMsgRemoveViewers(creator string, viewerIDs string, address string, fileowner string) *MsgRemoveViewers {
	return &MsgRemoveViewers{
		Creator:   creator,
		ViewerIds: viewerIDs,
		Address:   address,
		FileOwner: fileowner,
	}
}

func (msg *MsgRemoveViewers) Route() string {
	return RouterKey
}

func (msg *MsgRemoveViewers) Type() string {
	return TypeMsgRemoveViewers
}

func (msg *MsgRemoveViewers) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveViewers) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveViewers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.ViewerIds == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
			"invalid viewer ids: %s", msg.ViewerIds)
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
